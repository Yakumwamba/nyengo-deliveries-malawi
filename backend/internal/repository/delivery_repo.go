package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"nyengo-deliveries/internal/models"
)

// DeliveryRepository handles delivery tracking data access
type DeliveryRepository struct {
	db *pgxpool.Pool
}

// NewDeliveryRepository creates a new delivery repository
func NewDeliveryRepository(db *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// CreateTracking creates a new delivery tracking record
func (r *DeliveryRepository) CreateTracking(ctx context.Context, tracking *models.DeliveryTracking) error {
	query := `
		INSERT INTO delivery_tracking (
			id, order_id, courier_id, driver_name, driver_phone,
			vehicle_type, vehicle_plate, current_latitude, current_longitude,
			last_location_at, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	tracking.ID = uuid.New()
	tracking.CreatedAt = time.Now()
	tracking.UpdatedAt = time.Now()
	tracking.IsActive = true

	_, err := r.db.Exec(ctx, query,
		tracking.ID,
		tracking.OrderID,
		tracking.CourierID,
		tracking.DriverName,
		tracking.DriverPhone,
		tracking.VehicleType,
		tracking.VehiclePlate,
		tracking.CurrentLatitude,
		tracking.CurrentLongitude,
		tracking.LastLocationAt,
		tracking.IsActive,
		tracking.CreatedAt,
		tracking.UpdatedAt,
	)

	return err
}

// GetByOrderID retrieves tracking data for an order
func (r *DeliveryRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.DeliveryTracking, error) {
	query := `
		SELECT id, order_id, courier_id, driver_name, driver_phone,
			vehicle_type, vehicle_plate, current_latitude, current_longitude,
			last_location_at, route_polyline, estimated_arrival, distance_remaining,
			duration_remaining, is_active, created_at, updated_at
		FROM delivery_tracking
		WHERE order_id = $1 AND is_active = true
	`

	var tracking models.DeliveryTracking
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&tracking.ID,
		&tracking.OrderID,
		&tracking.CourierID,
		&tracking.DriverName,
		&tracking.DriverPhone,
		&tracking.VehicleType,
		&tracking.VehiclePlate,
		&tracking.CurrentLatitude,
		&tracking.CurrentLongitude,
		&tracking.LastLocationAt,
		&tracking.RoutePolyline,
		&tracking.EstimatedArrival,
		&tracking.DistanceRemaining,
		&tracking.DurationRemaining,
		&tracking.IsActive,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &tracking, nil
}

// UpdateLocation updates the current location of a delivery
func (r *DeliveryRepository) UpdateLocation(ctx context.Context, orderID uuid.UUID, lat, lng float64) error {
	query := `
		UPDATE delivery_tracking SET
			current_latitude = $2,
			current_longitude = $3,
			last_location_at = $4,
			updated_at = $4
		WHERE order_id = $1 AND is_active = true
	`

	now := time.Now()
	_, err := r.db.Exec(ctx, query, orderID, lat, lng, now)
	return err
}

// UpdateETA updates the estimated arrival time
func (r *DeliveryRepository) UpdateETA(ctx context.Context, orderID uuid.UUID, eta time.Time, distanceRemaining float64, durationRemaining int) error {
	query := `
		UPDATE delivery_tracking SET
			estimated_arrival = $2,
			distance_remaining = $3,
			duration_remaining = $4,
			updated_at = $5
		WHERE order_id = $1 AND is_active = true
	`

	_, err := r.db.Exec(ctx, query, orderID, eta, distanceRemaining, durationRemaining, time.Now())
	return err
}

// Complete marks a delivery tracking as completed
func (r *DeliveryRepository) Complete(ctx context.Context, orderID uuid.UUID) error {
	query := `
		UPDATE delivery_tracking SET
			is_active = false,
			updated_at = $2
		WHERE order_id = $1
	`

	_, err := r.db.Exec(ctx, query, orderID, time.Now())
	return err
}

// GetActiveDeliveriesForCourier gets all active deliveries for a courier
func (r *DeliveryRepository) GetActiveDeliveriesForCourier(ctx context.Context, courierID uuid.UUID) ([]models.DeliveryTracking, error) {
	query := `
		SELECT id, order_id, courier_id, driver_name, driver_phone,
			vehicle_type, current_latitude, current_longitude,
			last_location_at, estimated_arrival, distance_remaining
		FROM delivery_tracking
		WHERE courier_id = $1 AND is_active = true
	`

	rows, err := r.db.Query(ctx, query, courierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.DeliveryTracking
	for rows.Next() {
		var t models.DeliveryTracking
		err := rows.Scan(
			&t.ID,
			&t.OrderID,
			&t.CourierID,
			&t.DriverName,
			&t.DriverPhone,
			&t.VehicleType,
			&t.CurrentLatitude,
			&t.CurrentLongitude,
			&t.LastLocationAt,
			&t.EstimatedArrival,
			&t.DistanceRemaining,
		)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, t)
	}

	return deliveries, nil
}

// SaveLocationHistory saves a location point to history
func (r *DeliveryRepository) SaveLocationHistory(ctx context.Context, trackingID uuid.UUID, point models.LocationPoint) error {
	query := `
		INSERT INTO location_history (id, tracking_id, latitude, longitude, speed, heading, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		uuid.New(),
		trackingID,
		point.Latitude,
		point.Longitude,
		point.Speed,
		point.Heading,
		point.Timestamp,
	)

	return err
}

// GetLocationHistory retrieves location history for a tracking record
func (r *DeliveryRepository) GetLocationHistory(ctx context.Context, trackingID uuid.UUID) ([]models.LocationPoint, error) {
	query := `
		SELECT latitude, longitude, speed, heading, recorded_at
		FROM location_history
		WHERE tracking_id = $1
		ORDER BY recorded_at ASC
	`

	rows, err := r.db.Query(ctx, query, trackingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.LocationPoint
	for rows.Next() {
		var p models.LocationPoint
		err := rows.Scan(&p.Latitude, &p.Longitude, &p.Speed, &p.Heading, &p.Timestamp)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return points, nil
}
