package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/repository"
	"nyengo-deliveries/internal/services"
)

// TrackingHandler handles live tracking API endpoints
type TrackingHandler struct {
	trackingService *services.TrackingService
	orderRepo       *repository.OrderRepository
}

// NewTrackingHandler creates a new tracking handler
func NewTrackingHandler(trackingService *services.TrackingService, orderRepo *repository.OrderRepository) *TrackingHandler {
	return &TrackingHandler{
		trackingService: trackingService,
		orderRepo:       orderRepo,
	}
}

// resolveOrderID resolves an order identifier (UUID or order number) to a UUID
func (h *TrackingHandler) resolveOrderID(c *fiber.Ctx, orderIDParam string) (uuid.UUID, error) {
	// Try to parse as UUID first
	orderID, err := uuid.Parse(orderIDParam)
	if err == nil {
		return orderID, nil
	}

	// Not a UUID, might be an order number (e.g., NYG-20251226-AF857C71)
	if strings.HasPrefix(orderIDParam, "NYG-") {
		order, err := h.orderRepo.GetByOrderNumber(c.Context(), orderIDParam)
		if err != nil {
			return uuid.Nil, err
		}
		return order.ID, nil
	}

	return uuid.Nil, err
}

// StartTracking initiates tracking for an order
// POST /api/v1/tracking/:orderId/start
func (h *TrackingHandler) StartTracking(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	var driverInfo services.DriverInfo
	if err := c.BodyParser(&driverInfo); err != nil {
		return BadRequest(c, "Invalid driver info")
	}

	if driverInfo.Name == "" || driverInfo.Phone == "" {
		return BadRequest(c, "Driver name and phone are required")
	}

	if err := h.trackingService.StartTracking(c.Context(), orderID, &driverInfo); err != nil {
		return ServerError(c, err.Error())
	}

	return Success(c, fiber.Map{
		"message": "Tracking started",
		"orderId": orderID,
	})
}

// UpdateLocation receives location update from driver
// POST /api/v1/tracking/:orderId/location
func (h *TrackingHandler) UpdateLocation(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  float64 `json:"accuracy,omitempty"`
		Speed     float64 `json:"speed,omitempty"`
		Heading   float64 `json:"heading,omitempty"`
		Altitude  float64 `json:"altitude,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid location data")
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return BadRequest(c, "Latitude and longitude are required")
	}

	update := &services.LocationUpdate{
		OrderID:   orderID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Accuracy:  req.Accuracy,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Altitude:  req.Altitude,
	}

	delivery, err := h.trackingService.UpdateLocation(c.Context(), update)
	if err != nil {
		return ServerError(c, err.Error())
	}

	return Success(c, fiber.Map{
		"location": fiber.Map{
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
		},
		"distanceRemaining": delivery.DistanceRemaining,
		"etaMinutes":        delivery.ETAMinutes,
		"estimatedArrival":  delivery.EstimatedArrival,
	})
}

// GetLiveTracking retrieves current tracking data
// GET /api/v1/tracking/:orderId
// Accepts both UUID and order number (e.g., NYG-20251226-AF857C71)
func (h *TrackingHandler) GetLiveTracking(c *fiber.Ctx) error {
	orderIDParam := c.Params("orderId")

	orderID, err := h.resolveOrderID(c, orderIDParam)
	if err != nil {
		return BadRequest(c, "Invalid order ID or order number not found")
	}

	delivery, err := h.trackingService.GetLiveTracking(c.Context(), orderID)
	if err != nil {
		return NotFound(c, "No active tracking for this order")
	}

	return Success(c, delivery)
}

// GetLocationHistory retrieves location history
// GET /api/v1/tracking/:orderId/history
func (h *TrackingHandler) GetLocationHistory(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	limit := c.QueryInt("limit", 100)
	if limit > 1000 {
		limit = 1000
	}

	history, err := h.trackingService.GetLocationHistory(c.Context(), orderID, limit)
	if err != nil {
		return ServerError(c, err.Error())
	}

	return Success(c, fiber.Map{
		"orderId": orderID,
		"points":  history,
		"count":   len(history),
	})
}

// StopTracking ends tracking for an order
// POST /api/v1/tracking/:orderId/stop
func (h *TrackingHandler) StopTracking(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.BodyParser(&req)

	if req.Reason == "" {
		req.Reason = "completed"
	}

	if err := h.trackingService.StopTracking(c.Context(), orderID, req.Reason); err != nil {
		return ServerError(c, err.Error())
	}

	return Success(c, fiber.Map{
		"message": "Tracking stopped",
		"orderId": orderID,
	})
}
