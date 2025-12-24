package services

import (
	"math"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/models"
)

type PricingService struct {
	cfg *config.Config
}

func NewPricingService(cfg *config.Config) *PricingService {
	return &PricingService{cfg: cfg}
}

func (s *PricingService) CalculateEstimate(req *models.PriceEstimateRequest) (*models.PriceEstimateResponse, error) {
	distance := s.CalculateDistance(req.PickupLatitude, req.PickupLongitude, req.DeliveryLatitude, req.DeliveryLongitude)
	duration := int(distance / 30 * 60)
	if duration < 10 {
		duration = 10
	}

	baseFare := s.cfg.MinimumFare
	distanceFare := distance * s.cfg.BaseRatePerKm
	weightFare, fragileFare, expressFare := 0.0, 0.0, 0.0

	if req.PackageWeight > 5 {
		weightFare = (req.PackageWeight - 5) * 2
	}
	if req.IsFragile {
		fragileFare = baseFare * 0.3
	}
	if req.IsExpress {
		expressFare = baseFare * 0.5
	}

	surgeFare := 0.0
	if s.cfg.SurgePricingMult > 1.0 {
		surgeFare = (baseFare + distanceFare) * (s.cfg.SurgePricingMult - 1)
	}

	subTotal := baseFare + distanceFare + weightFare + fragileFare + expressFare + surgeFare
	if subTotal < s.cfg.MinimumFare {
		subTotal = s.cfg.MinimumFare
	}
	platformFee := subTotal * s.cfg.PlatformFeePerc
	totalFare := math.Round((subTotal+platformFee)*100) / 100

	tier := "Standard"
	if req.IsExpress {
		tier = "Express"
	}

	return &models.PriceEstimateResponse{
		Currency: s.cfg.Currency, CurrencySymbol: s.cfg.CurrencySymbol,
		Distance: math.Round(distance*100) / 100, Duration: duration,
		BaseFare: baseFare, DistanceFare: distanceFare, WeightFare: weightFare,
		FragileFare: fragileFare, ExpressFare: expressFare, SurgeFare: surgeFare,
		SubTotal: subTotal, PlatformFee: platformFee, TotalFare: totalFare,
		FormattedTotal: s.cfg.FormatCurrency(totalFare), PricingTier: tier,
		IsSurgeActive: s.cfg.SurgePricingMult > 1.0, SurgeMultiplier: s.cfg.SurgePricingMult,
		Disclaimer: "Prices are estimates and may vary.",
		// Add delivery type info
		IsLocalDelivery: s.IsLocalDelivery(distance),
		DeliveryType:    s.GetDeliveryType(distance),
	}, nil
}

func (s *PricingService) CalculateCourierEarnings(totalFare float64) (float64, float64) {
	fee := totalFare * s.cfg.PlatformFeePerc
	return math.Round(fee*100) / 100, math.Round((totalFare-fee)*100) / 100
}

// CalculateDistance calculates the distance between two points using Haversine formula
func (s *PricingService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) * 1.3
}

// IsLocalDelivery determines if a delivery is local based on distance threshold
func (s *PricingService) IsLocalDelivery(distance float64) bool {
	return distance < s.cfg.LocalDistanceThreshold
}

// GetDeliveryType returns the delivery type string based on distance
func (s *PricingService) GetDeliveryType(distance float64) string {
	if s.IsLocalDelivery(distance) {
		return "local"
	}
	return "intercity"
}

// GetDistanceThreshold returns the configured distance threshold
func (s *PricingService) GetDistanceThreshold() float64 {
	return s.cfg.LocalDistanceThreshold
}

// CalculateLocalCourierFare calculates fare for a local courier
func (s *PricingService) CalculateLocalCourierFare(distance float64, baseRatePerKm, minimumFare float64) float64 {
	fare := baseRatePerKm * distance
	if fare < minimumFare {
		fare = minimumFare
	}
	return math.Round(fare*100) / 100
}
