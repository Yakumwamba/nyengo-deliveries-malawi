package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusAccepted  OrderStatus = "accepted"
	OrderStatusDeclined  OrderStatus = "declined"
	OrderStatusPickedUp  OrderStatus = "picked_up"
	OrderStatusInTransit OrderStatus = "in_transit"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusFailed    OrderStatus = "failed"
)

// PaymentStatus represents the payment state
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// PaymentMethod represents how the customer pays
type PaymentMethod string

const (
	PaymentMethodCash        PaymentMethod = "cash"
	PaymentMethodMobileMoney PaymentMethod = "mobile_money"
	PaymentMethodCard        PaymentMethod = "card"
	PaymentMethodWallet      PaymentMethod = "wallet"
)

// Order represents a delivery order
type Order struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	OrderNumber     string     `json:"orderNumber" db:"order_number"`
	CourierID       uuid.UUID  `json:"courierId" db:"courier_id"`
	StoreID         *uuid.UUID `json:"storeId,omitempty" db:"store_id"`
	ExternalOrderID string     `json:"externalOrderId,omitempty" db:"external_order_id"`

	// Customer information
	CustomerName  string `json:"customerName" db:"customer_name"`
	CustomerPhone string `json:"customerPhone" db:"customer_phone"`
	CustomerEmail string `json:"customerEmail,omitempty" db:"customer_email"`

	// Pickup details
	PickupAddress      string  `json:"pickupAddress" db:"pickup_address"`
	PickupLatitude     float64 `json:"pickupLatitude" db:"pickup_latitude"`
	PickupLongitude    float64 `json:"pickupLongitude" db:"pickup_longitude"`
	PickupNotes        string  `json:"pickupNotes,omitempty" db:"pickup_notes"`
	PickupContactName  string  `json:"pickupContactName,omitempty" db:"pickup_contact_name"`
	PickupContactPhone string  `json:"pickupContactPhone,omitempty" db:"pickup_contact_phone"`

	// Delivery details
	DeliveryAddress   string  `json:"deliveryAddress" db:"delivery_address"`
	DeliveryLatitude  float64 `json:"deliveryLatitude" db:"delivery_latitude"`
	DeliveryLongitude float64 `json:"deliveryLongitude" db:"delivery_longitude"`
	DeliveryNotes     string  `json:"deliveryNotes,omitempty" db:"delivery_notes"`

	// Package details
	PackageDescription string  `json:"packageDescription" db:"package_description"`
	PackageSize        string  `json:"packageSize" db:"package_size"`     // small, medium, large
	PackageWeight      float64 `json:"packageWeight" db:"package_weight"` // in kg
	IsFragile          bool    `json:"isFragile" db:"is_fragile"`
	RequiresSignature  bool    `json:"requiresSignature" db:"requires_signature"`

	// Pricing breakdown
	Distance        float64 `json:"distance" db:"distance"` // in km
	BaseFare        float64 `json:"baseFare" db:"base_fare"`
	DistanceFare    float64 `json:"distanceFare" db:"distance_fare"`
	SurgeFare       float64 `json:"surgeFare" db:"surge_fare"`
	TotalFare       float64 `json:"totalFare" db:"total_fare"`
	PlatformFee     float64 `json:"platformFee" db:"platform_fee"`
	CourierEarnings float64 `json:"courierEarnings" db:"courier_earnings"`

	// Payment
	PaymentMethod    PaymentMethod `json:"paymentMethod" db:"payment_method"`
	PaymentStatus    PaymentStatus `json:"paymentStatus" db:"payment_status"`
	PaymentReference string        `json:"paymentReference,omitempty" db:"payment_reference"`

	// Status
	Status        OrderStatus    `json:"status" db:"status"`
	StatusHistory []StatusChange `json:"statusHistory,omitempty" db:"status_history"`

	// Scheduling
	ScheduledPickup   *time.Time `json:"scheduledPickup,omitempty" db:"scheduled_pickup"`
	ActualPickup      *time.Time `json:"actualPickup,omitempty" db:"actual_pickup"`
	EstimatedDelivery *time.Time `json:"estimatedDelivery,omitempty" db:"estimated_delivery"`
	ActualDelivery    *time.Time `json:"actualDelivery,omitempty" db:"actual_delivery"`

	// Proof of delivery
	DeliveryProofURL string `json:"deliveryProofUrl,omitempty" db:"delivery_proof_url"`
	RecipientName    string `json:"recipientName,omitempty" db:"recipient_name"`
	SignatureURL     string `json:"signatureUrl,omitempty" db:"signature_url"`

	// Rating
	CustomerRating   *int   `json:"customerRating,omitempty" db:"customer_rating"`
	CustomerFeedback string `json:"customerFeedback,omitempty" db:"customer_feedback"`

	// Metadata
	Notes    string         `json:"notes,omitempty" db:"notes"`
	Metadata map[string]any `json:"metadata,omitempty" db:"metadata"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// StatusChange records order status transitions
type StatusChange struct {
	Status    OrderStatus `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Note      string      `json:"note,omitempty"`
	Actor     string      `json:"actor"` // courier, system, customer
}

// CreateOrderRequest is the request body for creating a new order
type CreateOrderRequest struct {
	// Customer
	CustomerName  string `json:"customerName" validate:"required"`
	CustomerPhone string `json:"customerPhone" validate:"required"`
	CustomerEmail string `json:"customerEmail,omitempty" validate:"omitempty,email"`

	// Pickup
	PickupAddress      string  `json:"pickupAddress" validate:"required"`
	PickupLatitude     float64 `json:"pickupLatitude" validate:"required"`
	PickupLongitude    float64 `json:"pickupLongitude" validate:"required"`
	PickupNotes        string  `json:"pickupNotes,omitempty"`
	PickupContactName  string  `json:"pickupContactName,omitempty"`
	PickupContactPhone string  `json:"pickupContactPhone,omitempty"`

	// Delivery
	DeliveryAddress   string  `json:"deliveryAddress" validate:"required"`
	DeliveryLatitude  float64 `json:"deliveryLatitude" validate:"required"`
	DeliveryLongitude float64 `json:"deliveryLongitude" validate:"required"`
	DeliveryNotes     string  `json:"deliveryNotes,omitempty"`

	// Package
	PackageDescription string  `json:"packageDescription" validate:"required"`
	PackageSize        string  `json:"packageSize" validate:"required,oneof=small medium large"`
	PackageWeight      float64 `json:"packageWeight,omitempty"`
	IsFragile          bool    `json:"isFragile,omitempty"`
	RequiresSignature  bool    `json:"requiresSignature,omitempty"`

	// Payment
	PaymentMethod PaymentMethod `json:"paymentMethod" validate:"required"`

	// Scheduling
	ScheduledPickup *time.Time `json:"scheduledPickup,omitempty"`

	// External reference
	ExternalOrderID string     `json:"externalOrderId,omitempty"`
	StoreID         *uuid.UUID `json:"storeId,omitempty"`
}

// UpdateOrderStatusRequest is the request for updating order status
type UpdateOrderStatusRequest struct {
	Status        OrderStatus `json:"status" validate:"required"`
	Note          string      `json:"note,omitempty"`
	ProofURL      string      `json:"proofUrl,omitempty"`
	Signature     string      `json:"signature,omitempty"`
	RecipientName string      `json:"recipientName,omitempty"`
}

// OrderListFilters contains filters for listing orders
type OrderListFilters struct {
	Status    []OrderStatus `json:"status,omitempty"`
	DateFrom  *time.Time    `json:"dateFrom,omitempty"`
	DateTo    *time.Time    `json:"dateTo,omitempty"`
	Search    string        `json:"search,omitempty"`
	SortBy    string        `json:"sortBy,omitempty"`
	SortOrder string        `json:"sortOrder,omitempty"`
	Page      int           `json:"page,omitempty"`
	PageSize  int           `json:"pageSize,omitempty"`
}

// OrderListResponse contains paginated order results
type OrderListResponse struct {
	Orders     []Order `json:"orders"`
	TotalCount int     `json:"totalCount"`
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	TotalPages int     `json:"totalPages"`
}
