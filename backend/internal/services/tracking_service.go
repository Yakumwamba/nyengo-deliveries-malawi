package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/repository"
)

// TrackingService handles real-time delivery tracking
type TrackingService struct {
	redis        *redis.Client
	deliveryRepo *repository.DeliveryRepository
	orderRepo    *repository.OrderRepository

	// In-memory cache for active deliveries (for fast lookups)
	activeDeliveries sync.Map // map[orderID]*LiveDelivery
}

// LiveDelivery represents an active delivery being tracked
type LiveDelivery struct {
	OrderID      uuid.UUID `json:"orderId"`
	OrderNumber  string    `json:"orderNumber"` // Human-readable order number (NYG-*)
	CourierID    uuid.UUID `json:"courierId"`
	DriverName   string    `json:"driverName"`
	DriverPhone  string    `json:"driverPhone"`
	VehicleType  string    `json:"vehicleType"`
	VehiclePlate string    `json:"vehiclePlate,omitempty"`

	// Current position
	CurrentLocation Location  `json:"currentLocation"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`

	// Destination
	DestinationLat float64 `json:"destinationLat"`
	DestinationLng float64 `json:"destinationLng"`

	// ETA calculations
	DistanceRemaining float64   `json:"distanceRemaining"` // km
	ETAMinutes        int       `json:"etaMinutes"`
	EstimatedArrival  time.Time `json:"estimatedArrival"`

	// Status
	Status   string `json:"status"`
	IsActive bool   `json:"isActive"`
}

// Location represents a GPS position with metadata
type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy,omitempty"` // meters
	Speed     float64   `json:"speed,omitempty"`    // km/h
	Heading   float64   `json:"heading,omitempty"`  // degrees
	Altitude  float64   `json:"altitude,omitempty"` // meters
	Timestamp time.Time `json:"timestamp"`
}

// LocationUpdate is the payload from driver's device
type LocationUpdate struct {
	OrderID   uuid.UUID `json:"orderId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy,omitempty"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
	Altitude  float64   `json:"altitude,omitempty"`
}

// TrackingEvent is broadcast to subscribers
type TrackingEvent struct {
	Type      string      `json:"type"` // "location_update", "eta_update", "status_change"
	OrderID   string      `json:"orderId"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewTrackingService creates a new tracking service
func NewTrackingService(redis *redis.Client, deliveryRepo *repository.DeliveryRepository, orderRepo *repository.OrderRepository) *TrackingService {
	service := &TrackingService{
		redis:        redis,
		deliveryRepo: deliveryRepo,
		orderRepo:    orderRepo,
	}

	// Start background cleanup for stale deliveries
	go service.cleanupStaleDeliveries()

	return service
}

// StartTracking initiates tracking for an order
func (s *TrackingService) StartTracking(ctx context.Context, orderID uuid.UUID, driverInfo *DriverInfo) error {
	// Get order details for destination
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	delivery := &LiveDelivery{
		OrderID:        orderID,
		OrderNumber:    order.OrderNumber, // Use orderNumber for tracking
		CourierID:      order.CourierID,
		DriverName:     driverInfo.Name,
		DriverPhone:    driverInfo.Phone,
		VehicleType:    driverInfo.VehicleType,
		VehiclePlate:   driverInfo.VehiclePlate,
		DestinationLat: order.DeliveryLatitude,
		DestinationLng: order.DeliveryLongitude,
		Status:         "tracking",
		IsActive:       true,
		LastUpdatedAt:  time.Now(),
	}

	// Store in memory using orderNumber as key
	s.activeDeliveries.Store(order.OrderNumber, delivery)

	// Store in Redis for persistence across restarts
	if s.redis != nil {
		data, _ := json.Marshal(delivery)
		s.redis.Set(ctx, s.getTrackingKey(order.OrderNumber), data, 24*time.Hour)
	}

	// Create database record
	tracking := &models.DeliveryTracking{
		OrderID:      orderID,
		CourierID:    order.CourierID,
		DriverName:   driverInfo.Name,
		DriverPhone:  driverInfo.Phone,
		VehicleType:  driverInfo.VehicleType,
		VehiclePlate: driverInfo.VehiclePlate,
		IsActive:     true,
	}

	if err := s.deliveryRepo.CreateTracking(ctx, tracking); err != nil {
		return fmt.Errorf("failed to create tracking record: %w", err)
	}

	// Broadcast tracking started event
	s.broadcastEvent(ctx, order.OrderNumber, "tracking_started", delivery)

	return nil
}

// UpdateLocation processes a location update from driver
func (s *TrackingService) UpdateLocation(ctx context.Context, update *LocationUpdate) (*LiveDelivery, error) {
	// First, we need to find the delivery by order ID
	// Try to load from memory or Redis using order ID to get the orderNumber
	order, err := s.orderRepo.GetByID(ctx, update.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	orderNumber := order.OrderNumber

	// Get active delivery using orderNumber
	val, exists := s.activeDeliveries.Load(orderNumber)
	if !exists {
		// Try to load from Redis
		if s.redis != nil {
			data, err := s.redis.Get(ctx, s.getTrackingKey(orderNumber)).Bytes()
			if err == nil {
				var delivery LiveDelivery
				if json.Unmarshal(data, &delivery) == nil {
					val = &delivery
					s.activeDeliveries.Store(orderNumber, &delivery)
					exists = true
				}
			}
		}

		if !exists {
			return nil, fmt.Errorf("no active tracking for order %s", orderNumber)
		}
	}

	delivery := val.(*LiveDelivery)

	// Update location
	now := time.Now()
	delivery.CurrentLocation = Location{
		Latitude:  update.Latitude,
		Longitude: update.Longitude,
		Accuracy:  update.Accuracy,
		Speed:     update.Speed,
		Heading:   update.Heading,
		Altitude:  update.Altitude,
		Timestamp: now,
	}
	delivery.LastUpdatedAt = now

	// Calculate distance remaining and ETA
	delivery.DistanceRemaining = s.calculateDistance(
		update.Latitude, update.Longitude,
		delivery.DestinationLat, delivery.DestinationLng,
	)

	// Calculate ETA based on speed or default
	avgSpeed := update.Speed
	if avgSpeed < 5 {
		avgSpeed = 25 // Default average speed in urban areas (km/h)
	}
	delivery.ETAMinutes = int(delivery.DistanceRemaining / avgSpeed * 60)
	if delivery.ETAMinutes < 1 {
		delivery.ETAMinutes = 1
	}
	delivery.EstimatedArrival = now.Add(time.Duration(delivery.ETAMinutes) * time.Minute)

	// Update in memory using orderNumber
	s.activeDeliveries.Store(orderNumber, delivery)

	// Update in Redis
	if s.redis != nil {
		data, _ := json.Marshal(delivery)
		s.redis.Set(ctx, s.getTrackingKey(orderNumber), data, 24*time.Hour)

		// Also store location in a sorted set for history
		locationData, _ := json.Marshal(delivery.CurrentLocation)
		s.redis.ZAdd(ctx, s.getHistoryKey(orderNumber), redis.Z{
			Score:  float64(now.UnixMilli()),
			Member: locationData,
		})

		// Keep only last 1000 points
		s.redis.ZRemRangeByRank(ctx, s.getHistoryKey(orderNumber), 0, -1001)
	}

	// Update database (async to not block)
	go func() {
		bgCtx := context.Background()
		_ = s.deliveryRepo.UpdateLocation(bgCtx, update.OrderID, update.Latitude, update.Longitude)
		_ = s.deliveryRepo.UpdateETA(bgCtx, update.OrderID, delivery.EstimatedArrival, delivery.DistanceRemaining, delivery.ETAMinutes*60)

		// Save to location history
		point := models.LocationPoint{
			Latitude:  update.Latitude,
			Longitude: update.Longitude,
			Speed:     update.Speed,
			Heading:   update.Heading,
			Timestamp: now,
		}

		// Get tracking ID
		if tracking, err := s.deliveryRepo.GetByOrderID(bgCtx, update.OrderID); err == nil {
			_ = s.deliveryRepo.SaveLocationHistory(bgCtx, tracking.ID, point)
		}
	}()

	// Broadcast location update to subscribers
	s.broadcastEvent(ctx, orderNumber, "location_update", map[string]interface{}{
		"location":          delivery.CurrentLocation,
		"distanceRemaining": delivery.DistanceRemaining,
		"etaMinutes":        delivery.ETAMinutes,
		"estimatedArrival":  delivery.EstimatedArrival,
	})

	return delivery, nil
}

// GetLiveTracking retrieves current tracking data for an order
func (s *TrackingService) GetLiveTracking(ctx context.Context, orderID uuid.UUID) (*LiveDelivery, error) {
	// Get order to retrieve orderNumber
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	orderNumber := order.OrderNumber

	// Check memory first using orderNumber
	if val, exists := s.activeDeliveries.Load(orderNumber); exists {
		return val.(*LiveDelivery), nil
	}

	// Check Redis
	if s.redis != nil {
		data, err := s.redis.Get(ctx, s.getTrackingKey(orderNumber)).Bytes()
		if err == nil {
			var delivery LiveDelivery
			if json.Unmarshal(data, &delivery) == nil {
				s.activeDeliveries.Store(orderNumber, &delivery)
				return &delivery, nil
			}
		}
	}

	return nil, fmt.Errorf("no tracking data for order %s", orderNumber)
}

// GetLocationHistory retrieves location history for an order
func (s *TrackingService) GetLocationHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]Location, error) {
	// Get order to retrieve orderNumber
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	orderNumber := order.OrderNumber

	if s.redis == nil {
		// Fall back to database
		tracking, err := s.deliveryRepo.GetByOrderID(ctx, orderID)
		if err != nil {
			return nil, err
		}

		points, err := s.deliveryRepo.GetLocationHistory(ctx, tracking.ID)
		if err != nil {
			return nil, err
		}

		locations := make([]Location, len(points))
		for i, p := range points {
			locations[i] = Location{
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
				Speed:     p.Speed,
				Heading:   p.Heading,
				Timestamp: p.Timestamp,
			}
		}
		return locations, nil
	}

	// Get from Redis using orderNumber
	results, err := s.redis.ZRevRange(ctx, s.getHistoryKey(orderNumber), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	locations := make([]Location, len(results))
	for i, result := range results {
		var loc Location
		if json.Unmarshal([]byte(result), &loc) == nil {
			locations[i] = loc
		}
	}

	return locations, nil
}

// StopTracking ends tracking for an order
func (s *TrackingService) StopTracking(ctx context.Context, orderID uuid.UUID, reason string) error {
	// Get order to retrieve orderNumber
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}
	orderNumber := order.OrderNumber

	// Remove from memory using orderNumber
	s.activeDeliveries.Delete(orderNumber)

	// Remove from Redis
	if s.redis != nil {
		s.redis.Del(ctx, s.getTrackingKey(orderNumber))
		// Keep history for 7 days
		s.redis.Expire(ctx, s.getHistoryKey(orderNumber), 7*24*time.Hour)
	}

	// Update database
	if err := s.deliveryRepo.Complete(ctx, orderID); err != nil {
		return err
	}

	// Broadcast tracking stopped
	s.broadcastEvent(ctx, orderNumber, "tracking_stopped", map[string]string{
		"reason": reason,
	})

	return nil
}

// SubscribeToOrder subscribes to real-time updates for an order
func (s *TrackingService) SubscribeToOrder(ctx context.Context, orderID uuid.UUID) (<-chan TrackingEvent, func()) {
	ch := make(chan TrackingEvent, 10)

	if s.redis == nil {
		return ch, func() { close(ch) }
	}

	// Get order to retrieve orderNumber
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ch, func() { close(ch) }
	}
	orderNumber := order.OrderNumber

	pubsub := s.redis.Subscribe(ctx, s.getChannelKey(orderNumber))

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				var event TrackingEvent
				if json.Unmarshal([]byte(msg.Payload), &event) == nil {
					select {
					case ch <- event:
					default:
						// Channel full, skip
					}
				}
			}
		}
	}()

	cancel := func() {
		pubsub.Close()
	}

	return ch, cancel
}

// broadcastEvent publishes a tracking event
func (s *TrackingService) broadcastEvent(ctx context.Context, orderNumber string, eventType string, data interface{}) {
	event := TrackingEvent{
		Type:      eventType,
		OrderID:   orderNumber, // Use orderNumber as the identifier
		Timestamp: time.Now(),
		Data:      data,
	}

	if s.redis != nil {
		payload, _ := json.Marshal(event)
		s.redis.Publish(ctx, s.getChannelKey(orderNumber), payload)
	}
}

// calculateDistance uses Haversine formula
func (s *TrackingService) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // Distance in km
}

// cleanupStaleDeliveries removes old tracking data
func (s *TrackingService) cleanupStaleDeliveries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		staleThreshold := time.Now().Add(-30 * time.Minute)

		s.activeDeliveries.Range(func(key, value interface{}) bool {
			delivery := value.(*LiveDelivery)
			if delivery.LastUpdatedAt.Before(staleThreshold) {
				s.activeDeliveries.Delete(key)
			}
			return true
		})
	}
}

// Helper functions for Redis keys - use orderNumber for human-readable tracking
func (s *TrackingService) getTrackingKey(orderNumber string) string {
	return fmt.Sprintf("tracking:%s", orderNumber)
}

func (s *TrackingService) getHistoryKey(orderNumber string) string {
	return fmt.Sprintf("tracking:%s:history", orderNumber)
}

func (s *TrackingService) getChannelKey(orderNumber string) string {
	return fmt.Sprintf("tracking:%s:events", orderNumber)
}

// DriverInfo contains driver details for tracking
type DriverInfo struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	VehicleType  string `json:"vehicleType"`
	VehiclePlate string `json:"vehiclePlate,omitempty"`
}
