package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/services"
	"nyengo-deliveries/internal/websocket"
)

type OrderHandler struct {
	service      *services.OrderService
	notification *services.NotificationService
	hub          *websocket.Hub
}

func NewOrderHandler(service *services.OrderService, notification *services.NotificationService, hub *websocket.Hub) *OrderHandler {
	return &OrderHandler{service: service, notification: notification, hub: hub}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	courierID := c.Locals("courier_id").(uuid.UUID)
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	order, err := h.service.Create(c.Context(), courierID, &req)
	if err != nil {
		return ServerError(c, err.Error())
	}

	return Created(c, order)
}

func (h *OrderHandler) List(c *fiber.Ctx) error {
	courierID := c.Locals("courier_id").(uuid.UUID)
	filters := &models.OrderListFilters{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("pageSize", 20),
		Search:   c.Query("search"),
	}

	result, err := h.service.List(c.Context(), courierID, filters)
	if err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, result)
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	order, err := h.service.GetByID(c.Context(), orderID)
	if err != nil {
		return NotFound(c, "Order not found")
	}
	return Success(c, order)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	if err := h.service.UpdateStatus(c.Context(), orderID, req.Status); err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, fiber.Map{"message": "Status updated"})
}

func (h *OrderHandler) Accept(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	if err := h.service.Accept(c.Context(), orderID); err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, fiber.Map{"message": "Order accepted"})
}

func (h *OrderHandler) Decline(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid order ID")
	}

	if err := h.service.Decline(c.Context(), orderID); err != nil {
		return ServerError(c, err.Error())
	}
	return Success(c, fiber.Map{"message": "Order declined"})
}
