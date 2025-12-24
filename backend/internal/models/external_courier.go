package models

// CourierType defines whether a courier is local or external
type CourierType string

const (
	CourierTypeLocal    CourierType = "local"
	CourierTypeExternal CourierType = "external"
)

// ExternalCourier represents an external courier service like DHL, FedEx, etc.
type ExternalCourier struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	LogoURL               string  `json:"logoUrl"`
	Description           string  `json:"description"`
	BaseRatePerKm         float64 `json:"baseRatePerKm"`
	MinimumFare           float64 `json:"minimumFare"`
	EstimatedDeliveryDays string  `json:"estimatedDeliveryDays"` // e.g., "2-3 business days"
	ServiceType           string  `json:"serviceType"`           // "express", "standard", "economy"
	TrackingURL           string  `json:"trackingUrl,omitempty"`
	IsActive              bool    `json:"isActive"`
}

// CourierOption is a unified struct for both local and external couriers
type CourierOption struct {
	ID          string      `json:"id"`
	Type        CourierType `json:"type"`
	Name        string      `json:"name"`
	LogoURL     string      `json:"logoUrl,omitempty"`
	Description string      `json:"description,omitempty"`

	// Pricing
	EstimatedFare float64 `json:"estimatedFare"`
	FormattedFare string  `json:"formattedFare"`
	BaseRatePerKm float64 `json:"baseRatePerKm"`
	MinimumFare   float64 `json:"minimumFare"`

	// Ratings (for local couriers)
	Rating          float64 `json:"rating,omitempty"`
	TotalReviews    int     `json:"totalReviews,omitempty"`
	TotalDeliveries int     `json:"totalDeliveries,omitempty"`
	IsVerified      bool    `json:"isVerified,omitempty"`
	IsFeatured      bool    `json:"isFeatured,omitempty"`

	// Delivery time
	EstimatedTime         string `json:"estimatedTime,omitempty"`         // e.g., "30-45 mins" for local
	EstimatedDeliveryDays string `json:"estimatedDeliveryDays,omitempty"` // e.g., "2-3 days" for external

	// External courier specific
	ServiceType string `json:"serviceType,omitempty"` // "express", "standard", "economy"
	TrackingURL string `json:"trackingUrl,omitempty"`

	// Recommendation
	RecommendedFor string `json:"recommendedFor,omitempty"` // "fastest", "cheapest", "best_rated"
}

// CourierOptionsResponse is the response for listing couriers for an order
type CourierOptionsResponse struct {
	DeliveryType    string          `json:"deliveryType"` // "local" or "intercity"
	IsLocalDelivery bool            `json:"isLocalDelivery"`
	Distance        float64         `json:"distance"`  // in km
	Threshold       float64         `json:"threshold"` // distance threshold used
	Couriers        []CourierOption `json:"couriers"`

	// Recommendations
	Recommended    *CourierOption `json:"recommended,omitempty"`
	CheapestOption *CourierOption `json:"cheapestOption,omitempty"`
	FastestOption  *CourierOption `json:"fastestOption,omitempty"`
}
