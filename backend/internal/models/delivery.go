package models

import (
	"time"

	"github.com/google/uuid"
)

// DeliveryTracking represents real-time delivery tracking data
type DeliveryTracking struct {
	ID           uuid.UUID `json:"id" db:"id"`
	OrderID      uuid.UUID `json:"orderId" db:"order_id"`
	CourierID    uuid.UUID `json:"courierId" db:"courier_id"`
	DriverName   string    `json:"driverName" db:"driver_name"`
	DriverPhone  string    `json:"driverPhone" db:"driver_phone"`
	VehicleType  string    `json:"vehicleType" db:"vehicle_type"`
	VehiclePlate string    `json:"vehiclePlate,omitempty" db:"vehicle_plate"`

	// Current location
	CurrentLatitude  float64   `json:"currentLatitude" db:"current_latitude"`
	CurrentLongitude float64   `json:"currentLongitude" db:"current_longitude"`
	LastLocationAt   time.Time `json:"lastLocationAt" db:"last_location_at"`

	// Route information
	RoutePolyline   string          `json:"routePolyline,omitempty" db:"route_polyline"`
	LocationHistory []LocationPoint `json:"locationHistory,omitempty" db:"location_history"`

	// ETA
	EstimatedArrival  *time.Time `json:"estimatedArrival,omitempty" db:"estimated_arrival"`
	DistanceRemaining float64    `json:"distanceRemaining" db:"distance_remaining"` // in km
	DurationRemaining int        `json:"durationRemaining" db:"duration_remaining"` // in seconds

	// Status
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// LocationPoint represents a GPS coordinate with timestamp
type LocationPoint struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Speed     float64   `json:"speed,omitempty"`   // km/h
	Heading   float64   `json:"heading,omitempty"` // degrees
}

// UpdateLocationRequest is the request for updating driver location
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Speed     float64 `json:"speed,omitempty"`
	Heading   float64 `json:"heading,omitempty"`
}

// DeliveryMetrics contains delivery performance metrics
type DeliveryMetrics struct {
	TotalDeliveries     int     `json:"totalDeliveries"`
	CompletedDeliveries int     `json:"completedDeliveries"`
	FailedDeliveries    int     `json:"failedDeliveries"`
	CancelledDeliveries int     `json:"cancelledDeliveries"`
	OnTimeRate          float64 `json:"onTimeRate"`          // percentage
	AverageDeliveryTime float64 `json:"averageDeliveryTime"` // in minutes
	AverageRating       float64 `json:"averageRating"`
}

// DriverAssignment represents a driver assigned to an order
type DriverAssignment struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	OrderID      uuid.UUID  `json:"orderId" db:"order_id"`
	DriverID     uuid.UUID  `json:"driverId" db:"driver_id"`
	DriverName   string     `json:"driverName" db:"driver_name"`
	DriverPhone  string     `json:"driverPhone" db:"driver_phone"`
	VehicleType  string     `json:"vehicleType" db:"vehicle_type"`
	VehiclePlate string     `json:"vehiclePlate,omitempty" db:"vehicle_plate"`
	AssignedAt   time.Time  `json:"assignedAt" db:"assigned_at"`
	AcceptedAt   *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	Status       string     `json:"status" db:"status"` // pending, accepted, rejected
}
