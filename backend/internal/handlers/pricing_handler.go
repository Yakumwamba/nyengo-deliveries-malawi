package handlers

import (
	"github.com/gofiber/fiber/v2"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/services"
)

type PricingHandler struct {
	service *services.PricingService
}

func NewPricingHandler(service *services.PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

func (h *PricingHandler) GetEstimate(c *fiber.Ctx) error {
	var req models.PriceEstimateRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	if req.PickupLatitude == 0 || req.PickupLongitude == 0 ||
		req.DeliveryLatitude == 0 || req.DeliveryLongitude == 0 {
		return BadRequest(c, "Pickup and delivery coordinates are required")
	}

	estimate, err := h.service.CalculateEstimate(&req)
	if err != nil {
		return ServerError(c, err.Error())
	}

	return Success(c, estimate)
}
