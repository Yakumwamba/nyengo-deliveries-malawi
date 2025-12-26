package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/repository"
	"nyengo-deliveries/internal/services"
)

type WebhookHandler struct {
	orderService        *services.OrderService
	notificationService *services.NotificationService
	orderRepo           *repository.OrderRepository
	cfg                 *config.Config
}

func NewWebhookHandler(
	orderService *services.OrderService,
	notificationService *services.NotificationService,
	orderRepo *repository.OrderRepository,
	cfg *config.Config,
) *WebhookHandler {
	return &WebhookHandler{
		orderService:        orderService,
		notificationService: notificationService,
		orderRepo:           orderRepo,
		cfg:                 cfg,
	}
}

func (h *WebhookHandler) HandlePayment(c *fiber.Ctx) error {
	// Handle payment webhook from payment provider
	return Success(c, fiber.Map{"received": true})
}

// DeliveryWebhookPayload represents the incoming webhook payload from courier
type DeliveryWebhookPayload struct {
	Event     string              `json:"event"`
	Timestamp string              `json:"timestamp"`
	Data      DeliveryWebhookData `json:"data"`
}

// DeliveryWebhookData represents the data field in webhook payload
type DeliveryWebhookData struct {
	OrderID         string `json:"orderId"`
	ID              string `json:"id"`
	ExternalOrderID string `json:"externalOrderId"`
	Status          string `json:"status"`
	DeliveryStatus  string `json:"deliveryStatus"`
	NewStatus       string `json:"newStatus"`
	TrackingNumber  string `json:"trackingNumber"`
}

// courierStatusMapping maps courier status to internal order status
var courierStatusMapping = map[string]models.OrderStatus{
	"pending":    models.OrderStatusPending,
	"accepted":   models.OrderStatusAccepted,
	"picked_up":  models.OrderStatusPickedUp,
	"in_transit": models.OrderStatusInTransit,
	"delivered":  models.OrderStatusDelivered,
	"cancelled":  models.OrderStatusCancelled,
	"failed":     models.OrderStatusFailed,
}

// HandleDeliveryWebhook processes delivery status updates from courier platform
// POST /api/delivery/webhook
func (h *WebhookHandler) HandleDeliveryWebhook(c *fiber.Ctx) error {
	// Validate webhook secret if configured
	if h.cfg.WebhookSecret != "" {
		secret := c.Get("X-Webhook-Secret")
		if secret != h.cfg.WebhookSecret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or missing webhook secret",
			})
		}
	}

	// Parse the payload
	var payload DeliveryWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON payload",
		})
	}

	// Get order identifier - try multiple fields
	orderIDStr := payload.Data.OrderID
	if orderIDStr == "" {
		orderIDStr = payload.Data.ID
	}
	externalOrderID := payload.Data.ExternalOrderID

	// Validate that at least one identifier is present
	if orderIDStr == "" && externalOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing order identifier. Provide orderId, id, or externalOrderId",
		})
	}

	// Get status from multiple possible fields
	status := payload.Data.Status
	if status == "" {
		status = payload.Data.DeliveryStatus
	}
	if status == "" {
		status = payload.Data.NewStatus
	}
	status = strings.ToLower(strings.TrimSpace(status))

	// Find the order
	var order *models.Order
	var err error

	// Try order ID first (as UUID)
	if orderIDStr != "" {
		orderID, parseErr := uuid.Parse(orderIDStr)
		if parseErr == nil {
			order, err = h.orderRepo.GetByID(c.Context(), orderID)
		}
	}

	// If not found by UUID, try by order number (NYG-...)
	if order == nil && orderIDStr != "" && strings.HasPrefix(orderIDStr, "NYG-") {
		order, err = h.orderRepo.GetByOrderNumber(c.Context(), orderIDStr)
	}

	// If still not found, try by external order ID
	if order == nil && externalOrderID != "" {
		// Try parsing as UUID first
		extID, parseErr := uuid.Parse(externalOrderID)
		if parseErr == nil {
			order, err = h.orderRepo.GetByID(c.Context(), extID)
		}
		// Try as order number
		if order == nil && strings.HasPrefix(externalOrderID, "NYG-") {
			order, err = h.orderRepo.GetByOrderNumber(c.Context(), externalOrderID)
		}
	}

	if order == nil || err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Order not found",
		})
	}

	// Map courier status to internal status
	internalStatus, exists := courierStatusMapping[status]
	if !exists {
		// Default to accepted/processing for unknown statuses
		internalStatus = models.OrderStatusAccepted
	}

	// Update order status
	if err := h.orderRepo.UpdateStatus(c.Context(), order.ID, internalStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update order status",
		})
	}

	// Update tracking number if provided
	if payload.Data.TrackingNumber != "" {
		// Could update order with tracking number if field exists
		// For now, we'll just log it
	}

	// Notify via Redis pub/sub if notification service is available
	if h.notificationService != nil {
		_ = h.notificationService.SendOrderUpdate(c.Context(), order.CourierID.String(), order.ID.String(), string(internalStatus))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Delivery status updated",
		"data": fiber.Map{
			"orderId":   order.ID,
			"newStatus": internalStatus,
		},
	})
}

// HandleDeliveryUpdate is the old handler - kept for backward compatibility
func (h *WebhookHandler) HandleDeliveryUpdate(c *fiber.Ctx) error {
	return h.HandleDeliveryWebhook(c)
}
