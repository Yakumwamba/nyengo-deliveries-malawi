package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/services"
)

// PaymentHandler handles payment and payout HTTP requests
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// VerifyPayment verifies payment status for an order
// @Summary Verify order payment status
// @Description Checks if payment has been received for a delivered order by querying the store's payment API
// @Tags Payments
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID (UUID)"
// @Success 200 {object} models.PaymentVerification
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/verify/{orderId} [get]
func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	orderIDStr := c.Params("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid order ID format",
		})
	}

	verification, err := h.paymentService.VerifyOrderPayment(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data:    verification,
	})
}

// RequestPayout handles courier payout requests
// @Summary Request payout for completed deliveries
// @Description Allows verified couriers to request payment for their delivered orders
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body models.PayoutRequest true "Payout request details"
// @Success 200 {object} models.PayoutResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/payouts [post]
func (h *PaymentHandler) RequestPayout(c *fiber.Ctx) error {
	// Get courier ID from JWT claims
	courierIDStr := c.Locals("courierID").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid authentication",
		})
	}

	var req models.PayoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	response, err := h.paymentService.RequestPayout(c.Context(), courierID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Message: response.Message,
		Data:    response,
	})
}

// GetPayableOrders retrieves orders eligible for payout
// @Summary Get payable orders
// @Description Returns list of delivered orders that are eligible for payout
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse{data=[]models.PayableOrder}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/payable-orders [get]
func (h *PaymentHandler) GetPayableOrders(c *fiber.Ctx) error {
	courierIDStr := c.Locals("courierID").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid authentication",
		})
	}

	orders, err := h.paymentService.GetPayableOrders(c.Context(), courierID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data:    orders,
	})
}

// GetPayoutHistory retrieves payout history
// @Summary Get payout history
// @Description Returns paginated list of courier's payout requests
// @Tags Payments
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (pending, processing, completed, failed, cancelled)"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Success 200 {object} SuccessResponse{data=[]models.Payout}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/payouts [get]
func (h *PaymentHandler) GetPayoutHistory(c *fiber.Ctx) error {
	courierIDStr := c.Locals("courierID").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid authentication",
		})
	}

	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	payouts, err := h.paymentService.GetPayoutHistory(c.Context(), courierID, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"payouts":  payouts,
			"page":     page,
			"pageSize": limit,
		},
	})
}

// GetPayoutByID retrieves a specific payout
// @Summary Get payout details
// @Description Returns details of a specific payout request
// @Tags Payments
// @Accept json
// @Produce json
// @Param payoutId path string true "Payout ID (UUID)"
// @Success 200 {object} models.Payout
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/payouts/{payoutId} [get]
func (h *PaymentHandler) GetPayoutByID(c *fiber.Ctx) error {
	payoutIDStr := c.Params("payoutId")
	payoutID, err := uuid.Parse(payoutIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid payout ID format",
		})
	}

	payout, err := h.paymentService.GetPayoutByID(c.Context(), payoutID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Error:   "Payout not found",
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data:    payout,
	})
}

// GetEarningsSummary retrieves earnings summary
// @Summary Get earnings summary
// @Description Returns summary of courier's earnings including available balance, pending, and total
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {object} models.EarningsSummary
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/earnings [get]
func (h *PaymentHandler) GetEarningsSummary(c *fiber.Ctx) error {
	courierIDStr := c.Locals("courierID").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid authentication",
		})
	}

	summary, err := h.paymentService.GetEarningsSummary(c.Context(), courierID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data:    summary,
	})
}

// GetWalletTransactions retrieves wallet transaction history
// @Summary Get wallet transactions
// @Description Returns paginated list of wallet transactions (earnings, payouts, etc.)
// @Tags Payments
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Success 200 {object} SuccessResponse{data=[]models.WalletTransaction}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/payments/wallet/transactions [get]
func (h *PaymentHandler) GetWalletTransactions(c *fiber.Ctx) error {
	courierIDStr := c.Locals("courierID").(string)
	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid authentication",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	transactions, err := h.paymentService.GetWalletTransactions(c.Context(), courierID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"transactions": transactions,
			"page":         page,
			"pageSize":     limit,
		},
	})
}

// ============================================================
// ADMIN ENDPOINTS
// ============================================================

// ProcessPayout processes a pending payout (admin only)
// @Summary Process payout request (Admin)
// @Description Approves or rejects a pending payout request
// @Tags Admin - Payments
// @Accept json
// @Produce json
// @Param payoutId path string true "Payout ID (UUID)"
// @Param request body ProcessPayoutRequest true "Processing decision"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/payments/payouts/{payoutId}/process [post]
func (h *PaymentHandler) ProcessPayout(c *fiber.Ctx) error {
	payoutIDStr := c.Params("payoutId")
	payoutID, err := uuid.Parse(payoutIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid payout ID format",
		})
	}

	var req ProcessPayoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := h.paymentService.ProcessPayout(c.Context(), payoutID, req.Approve, req.Reason); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	status := "rejected"
	if req.Approve {
		status = "approved"
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Payout " + status + " successfully",
	})
}

// ProcessPayoutRequest is the request body for processing payouts
type ProcessPayoutRequest struct {
	Approve bool   `json:"approve"`
	Reason  string `json:"reason,omitempty"`
}

// ============================================================
// STORE INTEGRATION ENDPOINTS
// ============================================================

// StoreVerifyPayment allows stores to query payment status
// @Summary Verify payment status (Store API)
// @Description Allows integrated stores to check payment verification status for their orders
// @Tags Store Integration
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID (UUID)"
// @Success 200 {object} models.PaymentVerification
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/stores/payments/verify/{orderId} [get]
func (h *PaymentHandler) StoreVerifyPayment(c *fiber.Ctx) error {
	orderIDStr := c.Params("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid order ID format",
		})
	}

	verification, err := h.paymentService.VerifyOrderPayment(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Data:    verification,
	})
}

// PaymentWebhook handles payment status updates from stores
// @Summary Payment webhook (Store Integration)
// @Description Receives payment confirmation webhooks from integrated stores
// @Tags Store Integration
// @Accept json
// @Produce json
// @Param request body PaymentWebhookPayload true "Payment update payload"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/webhooks/payment-confirm [post]
func (h *PaymentHandler) PaymentWebhook(c *fiber.Ctx) error {
	var payload PaymentWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid payload",
		})
	}

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Error:   "Invalid order ID",
		})
	}

	// Verify payment
	verification, err := h.paymentService.VerifyOrderPayment(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Payment status updated",
		Data:    verification,
	})
}

// PaymentWebhookPayload is the webhook payload structure
type PaymentWebhookPayload struct {
	OrderID         string  `json:"orderId"`
	ExternalOrderID string  `json:"externalOrderId,omitempty"`
	Status          string  `json:"status"` // "paid", "failed", "refunded"
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	TransactionRef  string  `json:"transactionRef"`
	PaymentMethod   string  `json:"paymentMethod"`
	PaidAt          string  `json:"paidAt,omitempty"`
}
