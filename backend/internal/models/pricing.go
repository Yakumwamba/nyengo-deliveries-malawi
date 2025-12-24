package models

// PriceEstimateRequest is the request for getting a delivery price estimate
type PriceEstimateRequest struct {
	// Pickup location
	PickupLatitude  float64 `json:"pickupLatitude" validate:"required"`
	PickupLongitude float64 `json:"pickupLongitude" validate:"required"`
	PickupAddress   string  `json:"pickupAddress,omitempty"`

	// Delivery location
	DeliveryLatitude  float64 `json:"deliveryLatitude" validate:"required"`
	DeliveryLongitude float64 `json:"deliveryLongitude" validate:"required"`
	DeliveryAddress   string  `json:"deliveryAddress,omitempty"`

	// Package details (optional for more accurate pricing)
	PackageSize   string  `json:"packageSize,omitempty"`   // small, medium, large
	PackageWeight float64 `json:"packageWeight,omitempty"` // in kg
	IsFragile     bool    `json:"isFragile,omitempty"`
	IsExpress     bool    `json:"isExpress,omitempty"`

	// Optional: specific courier ID for custom pricing
	CourierID string `json:"courierId,omitempty"`
}

// PriceEstimateResponse is the response containing price breakdown
type PriceEstimateResponse struct {
	// Currency information
	Currency       string `json:"currency"`
	CurrencySymbol string `json:"currencySymbol"`

	// Distance and time
	Distance float64 `json:"distance"` // in km
	Duration int     `json:"duration"` // estimated minutes

	// Price breakdown
	BaseFare     float64 `json:"baseFare"`
	DistanceFare float64 `json:"distanceFare"`
	WeightFare   float64 `json:"weightFare,omitempty"`
	FragileFare  float64 `json:"fragileFare,omitempty"`
	ExpressFare  float64 `json:"expressFare,omitempty"`
	SurgeFare    float64 `json:"surgeFare,omitempty"`

	// Totals
	SubTotal    float64 `json:"subTotal"`
	PlatformFee float64 `json:"platformFee"`
	TotalFare   float64 `json:"totalFare"`

	// Formatted prices for display
	FormattedTotal     string                  `json:"formattedTotal"`
	FormattedBreakdown PriceBreakdownFormatted `json:"formattedBreakdown"`

	// Additional info
	EstimatedPickup   string  `json:"estimatedPickup,omitempty"`
	EstimatedDelivery string  `json:"estimatedDelivery,omitempty"`
	SurgeMultiplier   float64 `json:"surgeMultiplier,omitempty"`
	IsSurgeActive     bool    `json:"isSurgeActive"`

	// Pricing tier applied
	PricingTier string `json:"pricingTier"`

	// Delivery type information
	IsLocalDelivery bool   `json:"isLocalDelivery"` // true if distance < threshold
	DeliveryType    string `json:"deliveryType"`    // "local" or "intercity"

	// Disclaimer
	Disclaimer string `json:"disclaimer"`
}

// PriceBreakdownFormatted contains formatted price strings
type PriceBreakdownFormatted struct {
	BaseFare     string `json:"baseFare"`
	DistanceFare string `json:"distanceFare"`
	WeightFare   string `json:"weightFare,omitempty"`
	FragileFare  string `json:"fragileFare,omitempty"`
	ExpressFare  string `json:"expressFare,omitempty"`
	SurgeFare    string `json:"surgeFare,omitempty"`
	SubTotal     string `json:"subTotal"`
	PlatformFee  string `json:"platformFee"`
	TotalFare    string `json:"totalFare"`
}

// SurgeZone represents a geographic area with surge pricing
type SurgeZone struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Radius      float64 `json:"radius"` // in km
	Multiplier  float64 `json:"multiplier"`
	Reason      string  `json:"reason"` // "high_demand", "weather", "event"
	ActiveUntil string  `json:"activeUntil"`
}

// CourierPricingOverview contains pricing info for a specific courier
type CourierPricingOverview struct {
	CourierID      string  `json:"courierId"`
	CompanyName    string  `json:"companyName"`
	BaseRatePerKm  float64 `json:"baseRatePerKm"`
	MinimumFare    float64 `json:"minimumFare"`
	Rating         float64 `json:"rating"`
	EstimatedFare  float64 `json:"estimatedFare"`
	FormattedFare  string  `json:"formattedFare"`
	EstimatedTime  string  `json:"estimatedTime"`
	RecommendedFor string  `json:"recommendedFor,omitempty"` // "fastest", "cheapest", "best_rated"
}

// MultiCourierEstimateResponse compares pricing across couriers
type MultiCourierEstimateResponse struct {
	Distance       float64                  `json:"distance"`
	Couriers       []CourierPricingOverview `json:"couriers"`
	Recommended    *CourierPricingOverview  `json:"recommended,omitempty"`
	CheapestOption *CourierPricingOverview  `json:"cheapestOption,omitempty"`
	FastestOption  *CourierPricingOverview  `json:"fastestOption,omitempty"`
}
