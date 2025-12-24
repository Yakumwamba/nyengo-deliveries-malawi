package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/services"
)

type CourierHandler struct {
	service *services.CourierService
}

func NewCourierHandler(service *services.CourierService) *CourierHandler {
	return &CourierHandler{service: service}
}

func (h *CourierHandler) Register(c *fiber.Ctx) error {
	var req models.CourierRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	courier, err := h.service.Register(c.Context(), &req)
	if err != nil {
		return BadRequest(c, err.Error())
	}

	return Created(c, courier)
}

func (h *CourierHandler) Login(c *fiber.Ctx) error {
	var req models.CourierLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	result, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return Unauthorized(c, err.Error())
	}

	return Success(c, result)
}

func (h *CourierHandler) GetProfile(c *fiber.Ctx) error {
	courierID := c.Locals("courier_id").(uuid.UUID)
	courier, err := h.service.GetProfile(c.Context(), courierID)
	if err != nil {
		return NotFound(c, "Courier not found")
	}
	return Success(c, courier)
}

func (h *CourierHandler) UpdateProfile(c *fiber.Ctx) error {
	courierID := c.Locals("courier_id").(uuid.UUID)
	var req models.CourierUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	courier, err := h.service.UpdateProfile(c.Context(), courierID, &req)
	if err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, courier)
}

func (h *CourierHandler) GetDashboard(c *fiber.Ctx) error {
	return Success(c, fiber.Map{"message": "Dashboard data"})
}

func (h *CourierHandler) ListAvailable(c *fiber.Ctx) error {
	area := c.Query("area")
	couriers, err := h.service.ListAvailable(c.Context(), area)
	if err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, couriers)
}

func (h *CourierHandler) GetRates(c *fiber.Ctx) error {
	courierID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid courier ID")
	}

	courier, err := h.service.GetByID(c.Context(), courierID)
	if err != nil {
		return NotFound(c, "Courier not found")
	}

	return Success(c, fiber.Map{
		"baseRatePerKm": courier.BaseRatePerKm,
		"minimumFare":   courier.MinimumFare,
		"rating":        courier.Rating,
		"totalReviews":  courier.TotalReviews,
	})
}
