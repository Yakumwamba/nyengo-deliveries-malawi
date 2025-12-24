package services

import (
	"math"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/models"
)

// ExternalCourierService handles external courier operations
type ExternalCourierService struct {
	cfg *config.Config
}

// NewExternalCourierService creates a new external courier service
func NewExternalCourierService(cfg *config.Config) *ExternalCourierService {
	return &ExternalCourierService{cfg: cfg}
}

// GetExternalCouriers returns all active external couriers
func (s *ExternalCourierService) GetExternalCouriers() []config.ExternalCourierConfig {
	active := []config.ExternalCourierConfig{}
	for _, c := range s.cfg.ExternalCouriers {
		if c.IsActive {
			active = append(active, c)
		}
	}
	return active
}

// CalculateExternalCourierOptions calculates pricing for all external couriers for a given distance
func (s *ExternalCourierService) CalculateExternalCourierOptions(distance float64) []models.CourierOption {
	options := []models.CourierOption{}

	for _, courier := range s.GetExternalCouriers() {
		fare := s.CalculateFare(courier, distance)

		option := models.CourierOption{
			ID:                    courier.ID,
			Type:                  models.CourierTypeExternal,
			Name:                  courier.Name,
			LogoURL:               courier.LogoURL,
			Description:           courier.Description,
			EstimatedFare:         fare,
			FormattedFare:         s.cfg.FormatCurrency(fare),
			BaseRatePerKm:         courier.BaseRatePerKm,
			MinimumFare:           courier.MinimumFare,
			EstimatedDeliveryDays: courier.EstimatedDeliveryDays,
			ServiceType:           courier.ServiceType,
		}

		// Set recommendations based on service type
		switch courier.ServiceType {
		case "express":
			option.RecommendedFor = "fastest"
		case "economy":
			option.RecommendedFor = "cheapest"
		}

		options = append(options, option)
	}

	// Sort by price and set cheapest/fastest recommendations
	s.setRecommendations(options)

	return options
}

// CalculateFare calculates the fare for an external courier
func (s *ExternalCourierService) CalculateFare(courier config.ExternalCourierConfig, distance float64) float64 {
	fare := courier.BaseRatePerKm * distance
	if fare < courier.MinimumFare {
		fare = courier.MinimumFare
	}
	// Round to 2 decimal places
	return math.Round(fare*100) / 100
}

// setRecommendations sets the recommended options (cheapest, fastest)
func (s *ExternalCourierService) setRecommendations(options []models.CourierOption) {
	if len(options) == 0 {
		return
	}

	cheapestIdx := 0
	fastestIdx := 0

	for i, opt := range options {
		if opt.EstimatedFare < options[cheapestIdx].EstimatedFare {
			cheapestIdx = i
		}
		if opt.ServiceType == "express" && options[fastestIdx].ServiceType != "express" {
			fastestIdx = i
		} else if opt.ServiceType == "express" && opt.EstimatedFare < options[fastestIdx].EstimatedFare {
			fastestIdx = i
		}
	}

	options[cheapestIdx].RecommendedFor = "cheapest"
	if fastestIdx != cheapestIdx {
		options[fastestIdx].RecommendedFor = "fastest"
	}
}

// IsLocalDelivery determines if a delivery is local based on distance threshold
func (s *ExternalCourierService) IsLocalDelivery(distance float64) bool {
	return distance < s.cfg.LocalDistanceThreshold
}

// GetDeliveryType returns the delivery type string based on distance
func (s *ExternalCourierService) GetDeliveryType(distance float64) string {
	if s.IsLocalDelivery(distance) {
		return "local"
	}
	return "intercity"
}

// GetDistanceThreshold returns the configured distance threshold
func (s *ExternalCourierService) GetDistanceThreshold() float64 {
	return s.cfg.LocalDistanceThreshold
}
