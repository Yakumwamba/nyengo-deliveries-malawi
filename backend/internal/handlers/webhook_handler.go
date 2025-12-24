package handlers

import (
	"github.com/gofiber/fiber/v2"

	"nyengo-deliveries/internal/services"
)

type WebhookHandler struct {
	orderService        *services.OrderService
	notificationService *services.NotificationService
}

func NewWebhookHandler(orderService *services.OrderService, notificationService *services.NotificationService) *WebhookHandler {
	return &WebhookHandler{orderService: orderService, notificationService: notificationService}
}

func (h *WebhookHandler) HandlePayment(c *fiber.Ctx) error {
	// Handle payment webhook from payment provider
	return Success(c, fiber.Map{"received": true})
}

func (h *WebhookHandler) HandleDeliveryUpdate(c *fiber.Ctx) error {
	// Handle delivery status webhook
	return Success(c, fiber.Map{"received": true})
}
