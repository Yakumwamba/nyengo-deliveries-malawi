package handlers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/services"
)

type StoreHandler struct {
	courierService         *services.CourierService
	orderService           *services.OrderService
	pricingService         *services.PricingService
	externalCourierService *services.ExternalCourierService
	cfg                    *config.Config
}

func NewStoreHandler(
	courierService *services.CourierService,
	orderService *services.OrderService,
	pricingService *services.PricingService,
	externalCourierService *services.ExternalCourierService,
	cfg *config.Config,
) *StoreHandler {
	return &StoreHandler{
		courierService:         courierService,
		orderService:           orderService,
		pricingService:         pricingService,
		externalCourierService: externalCourierService,
		cfg:                    cfg,
	}
}

// ListCouriers returns available couriers based on delivery distance
// Query params: pickupLat, pickupLon, deliveryLat, deliveryLon
// For local deliveries (< threshold), returns registered local couriers
// For inter-city deliveries (>= threshold), returns external courier services
func (h *StoreHandler) ListCouriers(c *fiber.Ctx) error {
	// Parse coordinates from query params
	pickupLat, err := strconv.ParseFloat(c.Query("pickupLat", "0"), 64)
	if err != nil {
		pickupLat = 0
	}
	pickupLon, err := strconv.ParseFloat(c.Query("pickupLon", "0"), 64)
	if err != nil {
		pickupLon = 0
	}
	deliveryLat, err := strconv.ParseFloat(c.Query("deliveryLat", "0"), 64)
	if err != nil {
		deliveryLat = 0
	}
	deliveryLon, err := strconv.ParseFloat(c.Query("deliveryLon", "0"), 64)
	if err != nil {
		deliveryLon = 0
	}

	// If coordinates provided, use distance-based selection
	if pickupLat != 0 && pickupLon != 0 && deliveryLat != 0 && deliveryLon != 0 {
		return h.listCouriersByDistance(c, pickupLat, pickupLon, deliveryLat, deliveryLon)
	}

	// Fallback to area-based selection (legacy)
	area := c.Query("area")
	couriers, err := h.courierService.ListAvailable(c.Context(), area)
	if err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, couriers)
}

// listCouriersByDistance returns couriers based on calculated distance
func (h *StoreHandler) listCouriersByDistance(c *fiber.Ctx, pickupLat, pickupLon, deliveryLat, deliveryLon float64) error {
	// Calculate distance
	distance := h.pricingService.CalculateDistance(pickupLat, pickupLon, deliveryLat, deliveryLon)
	distance = math.Round(distance*100) / 100

	isLocal := h.pricingService.IsLocalDelivery(distance)
	deliveryType := h.pricingService.GetDeliveryType(distance)
	threshold := h.pricingService.GetDistanceThreshold()

	var courierOptions []models.CourierOption
	var recommended, cheapest, fastest *models.CourierOption

	if isLocal {
		// Local delivery - return registered local couriers
		couriers, err := h.courierService.ListAvailable(c.Context(), "")
		if err != nil {
			return ServerError(c, err.Error())
		}

		for _, courier := range couriers {
			fare := h.pricingService.CalculateLocalCourierFare(distance, courier.BaseRatePerKm, courier.MinimumFare)
			estimatedTime := h.calculateEstimatedTime(distance)

			option := models.CourierOption{
				ID:              courier.ID.String(),
				Type:            models.CourierTypeLocal,
				Name:            courier.CompanyName,
				LogoURL:         courier.LogoURL,
				EstimatedFare:   fare,
				FormattedFare:   h.cfg.FormatCurrency(fare),
				BaseRatePerKm:   courier.BaseRatePerKm,
				MinimumFare:     courier.MinimumFare,
				Rating:          courier.Rating,
				TotalReviews:    courier.TotalReviews,
				TotalDeliveries: courier.TotalDeliveries,
				IsVerified:      courier.IsVerified,
				IsFeatured:      courier.IsFeatured,
				EstimatedTime:   estimatedTime,
			}
			courierOptions = append(courierOptions, option)
		}

		// Set recommendations for local couriers
		if len(courierOptions) > 0 {
			cheapest, fastest, recommended = h.findLocalRecommendations(courierOptions)
		}
	} else {
		// Inter-city delivery - return external couriers
		courierOptions = h.externalCourierService.CalculateExternalCourierOptions(distance)

		// Find recommendations for external couriers
		if len(courierOptions) > 0 {
			cheapest, fastest, recommended = h.findExternalRecommendations(courierOptions)
		}
	}

	response := models.CourierOptionsResponse{
		DeliveryType:    deliveryType,
		IsLocalDelivery: isLocal,
		Distance:        distance,
		Threshold:       threshold,
		Couriers:        courierOptions,
		Recommended:     recommended,
		CheapestOption:  cheapest,
		FastestOption:   fastest,
	}

	return Success(c, response)
}

// calculateEstimatedTime calculates estimated delivery time for local deliveries
func (h *StoreHandler) calculateEstimatedTime(distance float64) string {
	// Assume average speed of 30 km/h in city traffic
	minutes := int(distance / 30 * 60)
	if minutes < 15 {
		return "10-20 mins"
	} else if minutes < 30 {
		return "20-35 mins"
	} else if minutes < 45 {
		return "35-50 mins"
	} else if minutes < 60 {
		return "45-60 mins"
	}
	return "1-2 hours"
}

// findLocalRecommendations finds recommended couriers for local delivery
func (h *StoreHandler) findLocalRecommendations(options []models.CourierOption) (*models.CourierOption, *models.CourierOption, *models.CourierOption) {
	if len(options) == 0 {
		return nil, nil, nil
	}

	cheapestIdx := 0
	bestRatedIdx := 0

	for i, opt := range options {
		if opt.EstimatedFare < options[cheapestIdx].EstimatedFare {
			cheapestIdx = i
		}
		if opt.Rating > options[bestRatedIdx].Rating {
			bestRatedIdx = i
		}
	}

	cheapest := options[cheapestIdx]
	cheapest.RecommendedFor = "cheapest"

	bestRated := options[bestRatedIdx]
	bestRated.RecommendedFor = "best_rated"

	// Recommended is best rated if rating > 4.0, otherwise cheapest
	var recommended models.CourierOption
	if options[bestRatedIdx].Rating >= 4.0 {
		recommended = bestRated
		recommended.RecommendedFor = "recommended"
	} else {
		recommended = cheapest
		recommended.RecommendedFor = "recommended"
	}

	return &cheapest, &bestRated, &recommended
}

// findExternalRecommendations finds recommended couriers for external delivery
func (h *StoreHandler) findExternalRecommendations(options []models.CourierOption) (*models.CourierOption, *models.CourierOption, *models.CourierOption) {
	if len(options) == 0 {
		return nil, nil, nil
	}

	cheapestIdx := 0
	fastestIdx := 0

	for i, opt := range options {
		if opt.EstimatedFare < options[cheapestIdx].EstimatedFare {
			cheapestIdx = i
		}
		if opt.ServiceType == "express" {
			if options[fastestIdx].ServiceType != "express" || opt.EstimatedFare < options[fastestIdx].EstimatedFare {
				fastestIdx = i
			}
		}
	}

	cheapest := options[cheapestIdx]
	cheapest.RecommendedFor = "cheapest"

	fastest := options[fastestIdx]
	fastest.RecommendedFor = "fastest"

	// Recommended is express with best value (balance of speed and price)
	recommended := fastest
	recommended.RecommendedFor = "recommended"

	return &cheapest, &fastest, &recommended
}

func (h *StoreHandler) CreateOrder(c *fiber.Ctx) error {
	var req struct {
		CourierID string `json:"courierId"`
		models.CreateOrderRequest
	}
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return BadRequest(c, "Invalid courier ID")
	}

	order, err := h.orderService.Create(c.Context(), courierID, &req.CreateOrderRequest)
	if err != nil {
		return ServerError(c, err.Error())
	}

	return Created(c, order)
}

func (h *StoreHandler) GetOrderStatus(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	order, err := h.orderService.GetByID(c.Context(), orderID)
	if err != nil {
		return NotFound(c, "Order not found")
	}

	return Success(c, fiber.Map{
		"orderId":     order.ID,
		"orderNumber": order.OrderNumber,
		"status":      order.Status,
		"updatedAt":   order.UpdatedAt,
	})
}
