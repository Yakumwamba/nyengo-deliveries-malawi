package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/repository"
)

// PaymentVerifier interface for different store payment integrations
// Each store can implement their own payment verification logic
type PaymentVerifier interface {
	VerifyPayment(ctx context.Context, orderID, externalOrderID string) (*models.PaymentVerification, error)
	GetProviderName() string
}

// PaymentService handles all payment and payout operations
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
	courierRepo *repository.CourierRepository
	config      *config.Config
	verifiers   map[string]PaymentVerifier // Map of store ID to verifier
	httpClient  *http.Client
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	orderRepo *repository.OrderRepository,
	courierRepo *repository.CourierRepository,
	cfg *config.Config,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		courierRepo: courierRepo,
		config:      cfg,
		verifiers:   make(map[string]PaymentVerifier),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterVerifier registers a custom payment verifier for a store
func (s *PaymentService) RegisterVerifier(storeID string, verifier PaymentVerifier) {
	s.verifiers[storeID] = verifier
}

// ============================================================
// PAYMENT VERIFICATION
// ============================================================

// VerifyOrderPayment verifies if an order has been paid by querying the store's payment API
func (s *PaymentService) VerifyOrderPayment(ctx context.Context, orderID uuid.UUID) (*models.PaymentVerification, error) {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if order is delivered
	if order.Status != models.OrderStatusDelivered {
		return &models.PaymentVerification{
			OrderID: orderID,
			IsPaid:  false,
			Error:   "order not yet delivered",
		}, nil
	}

	// If payment is already confirmed in our system
	if order.PaymentStatus == models.PaymentStatusPaid {
		return &models.PaymentVerification{
			OrderID:       orderID,
			IsPaid:        true,
			AmountPaid:    order.TotalFare,
			PaymentMethod: string(order.PaymentMethod),
			PaymentRef:    order.PaymentReference,
		}, nil
	}

	// For cash payments, verify with proof of delivery
	if order.PaymentMethod == models.PaymentMethodCash {
		return s.verifyCashPayment(ctx, order)
	}

	// For store orders, verify with store's payment API
	if order.StoreID != nil {
		return s.verifyStorePayment(ctx, order)
	}

	// Default verification for other payment methods
	return s.verifyGenericPayment(ctx, order)
}

// verifyCashPayment verifies cash payment based on delivery proof
func (s *PaymentService) verifyCashPayment(ctx context.Context, order *models.Order) (*models.PaymentVerification, error) {
	verification := &models.PaymentVerification{
		OrderID:       order.ID,
		PaymentMethod: "cash",
	}

	// Cash is considered paid if delivery is confirmed with proof
	if order.DeliveryProofURL != "" || order.RecipientName != "" || order.SignatureURL != "" {
		verification.IsPaid = true
		verification.AmountPaid = order.TotalFare
		now := time.Now()
		verification.PaidAt = &now

		// Update order payment status
		if err := s.updateOrderPaymentStatus(ctx, order.ID, models.PaymentStatusPaid, "cash-verified"); err != nil {
			return nil, err
		}
	} else {
		verification.IsPaid = false
		verification.Error = "cash payment not confirmed - missing delivery proof"
	}

	return verification, nil
}

// verifyStorePayment verifies payment through the store's payment API
func (s *PaymentService) verifyStorePayment(ctx context.Context, order *models.Order) (*models.PaymentVerification, error) {
	if order.StoreID == nil {
		return nil, errors.New("order has no associated store")
	}

	// Check if we have a custom verifier for this store
	if verifier, ok := s.verifiers[order.StoreID.String()]; ok {
		return verifier.VerifyPayment(ctx, order.ID.String(), order.ExternalOrderID)
	}

	// Get store payment configuration
	storeConfig, err := s.paymentRepo.GetStorePaymentConfig(ctx, *order.StoreID)
	if err != nil {
		return nil, fmt.Errorf("store payment config not found: %w", err)
	}

	// Build verification request
	verifyReq := StorePaymentVerifyRequest{
		OrderID:         order.ID.String(),
		ExternalOrderID: order.ExternalOrderID,
		Amount:          order.TotalFare,
		Currency:        s.config.Currency,
	}

	reqBody, _ := json.Marshal(verifyReq)

	// Create HTTP request to store's payment API
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		storeConfig.PaymentAPIURL+"/verify-payment",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storeConfig.APIKey)
	req.Header.Set("X-Timestamp", time.Now().UTC().Format(time.RFC3339))

	// Add signature for security
	signature := s.generateSignature(reqBody, storeConfig.WebhookSecret)
	req.Header.Set("X-Signature", signature)

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &models.PaymentVerification{
			OrderID: order.ID,
			IsPaid:  false,
			Error:   fmt.Sprintf("failed to contact store API: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	// Parse response
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &models.PaymentVerification{
			OrderID: order.ID,
			IsPaid:  false,
			Error:   fmt.Sprintf("store API error: %s", string(body)),
		}, nil
	}

	var storeResp StorePaymentVerifyResponse
	if err := json.Unmarshal(body, &storeResp); err != nil {
		return nil, fmt.Errorf("failed to parse store response: %w", err)
	}

	verification := &models.PaymentVerification{
		OrderID:       order.ID,
		IsPaid:        storeResp.IsPaid,
		AmountPaid:    storeResp.AmountPaid,
		PaymentMethod: storeResp.PaymentMethod,
		PaymentRef:    storeResp.TransactionRef,
	}

	if storeResp.PaidAt != "" {
		paidAt, _ := time.Parse(time.RFC3339, storeResp.PaidAt)
		verification.PaidAt = &paidAt
	}

	// Update our order status if payment confirmed
	if storeResp.IsPaid {
		if err := s.updateOrderPaymentStatus(ctx, order.ID, models.PaymentStatusPaid, storeResp.TransactionRef); err != nil {
			return nil, err
		}
	}

	return verification, nil
}

// verifyGenericPayment handles verification for generic payment methods
func (s *PaymentService) verifyGenericPayment(ctx context.Context, order *models.Order) (*models.PaymentVerification, error) {
	// For mobile money, card, etc. - check if we have a payment reference
	if order.PaymentReference != "" {
		return &models.PaymentVerification{
			OrderID:       order.ID,
			IsPaid:        true,
			AmountPaid:    order.TotalFare,
			PaymentMethod: string(order.PaymentMethod),
			PaymentRef:    order.PaymentReference,
		}, nil
	}

	return &models.PaymentVerification{
		OrderID: order.ID,
		IsPaid:  false,
		Error:   "payment not yet confirmed",
	}, nil
}

// ============================================================
// PAYOUT OPERATIONS
// ============================================================

// RequestPayout creates a payout request for a courier
func (s *PaymentService) RequestPayout(ctx context.Context, courierID uuid.UUID, req *models.PayoutRequest) (*models.PayoutResponse, error) {
	// Validate courier exists
	courier, err := s.courierRepo.GetByID(ctx, courierID)
	if err != nil {
		return nil, fmt.Errorf("courier not found: %w", err)
	}

	if !courier.IsVerified {
		return nil, errors.New("courier account must be verified to request payouts")
	}

	// Parse and validate order IDs
	var orderIDs []uuid.UUID
	for _, id := range req.OrderIDs {
		orderID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid order ID: %s", id)
		}
		orderIDs = append(orderIDs, orderID)
	}

	if len(orderIDs) == 0 {
		return nil, errors.New("at least one order must be specified")
	}

	// Verify all orders are eligible for payout
	var totalAmount float64
	var verifiedOrders []uuid.UUID

	for _, orderID := range orderIDs {
		// Verify payment for this order
		verification, err := s.VerifyOrderPayment(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to verify order %s: %w", orderID, err)
		}

		if !verification.IsPaid {
			return nil, fmt.Errorf("order %s payment not confirmed: %s", orderID, verification.Error)
		}

		// Get order to confirm courier ownership
		order, err := s.orderRepo.GetByID(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("order %s not found", orderID)
		}

		if order.CourierID != courierID {
			return nil, fmt.Errorf("order %s does not belong to this courier", orderID)
		}

		// Check order hasn't already been paid out
		payable, err := s.paymentRepo.GetPayableOrders(ctx, courierID)
		if err != nil {
			return nil, err
		}

		isPayable := false
		for _, p := range payable {
			if p.OrderID == orderID && p.PayoutStatus == "unpaid" {
				isPayable = true
				totalAmount += p.CourierEarnings
				break
			}
		}

		if !isPayable {
			return nil, fmt.Errorf("order %s is not eligible for payout", orderID)
		}

		verifiedOrders = append(verifiedOrders, orderID)
	}

	// Calculate fees
	platformFee := totalAmount * s.config.PlatformFeePerc
	netAmount := totalAmount - platformFee

	// Create payout record
	payout := &models.Payout{
		ID:            uuid.New(),
		CourierID:     courierID,
		OrderIDs:      verifiedOrders,
		TotalAmount:   totalAmount,
		PlatformFee:   platformFee,
		NetAmount:     netAmount,
		Currency:      s.config.Currency,
		Status:        models.PayoutStatusPending,
		PayoutMethod:  req.PayoutMethod,
		PayoutDetails: req.PayoutDetails,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save payout
	if err := s.paymentRepo.CreatePayout(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	// Link orders to payout
	if err := s.paymentRepo.LinkOrdersToPayout(ctx, payout.ID, verifiedOrders); err != nil {
		return nil, fmt.Errorf("failed to link orders: %w", err)
	}

	return &models.PayoutResponse{
		PayoutID:    payout.ID,
		Status:      payout.Status,
		TotalAmount: fmt.Sprintf("%s %.2f", s.config.CurrencySymbol, totalAmount),
		NetAmount:   fmt.Sprintf("%s %.2f", s.config.CurrencySymbol, netAmount),
		Message:     "Payout request submitted successfully. Processing will begin shortly.",
	}, nil
}

// GetPayableOrders retrieves orders eligible for payout
func (s *PaymentService) GetPayableOrders(ctx context.Context, courierID uuid.UUID) ([]models.PayableOrder, error) {
	return s.paymentRepo.GetPayableOrders(ctx, courierID)
}

// GetPayoutHistory retrieves payout history for a courier
func (s *PaymentService) GetPayoutHistory(ctx context.Context, courierID uuid.UUID, status string, page, pageSize int) ([]models.Payout, error) {
	offset := (page - 1) * pageSize
	return s.paymentRepo.GetCourierPayouts(ctx, courierID, status, pageSize, offset)
}

// GetPayoutByID retrieves a specific payout
func (s *PaymentService) GetPayoutByID(ctx context.Context, payoutID uuid.UUID) (*models.Payout, error) {
	return s.paymentRepo.GetPayoutByID(ctx, payoutID)
}

// ProcessPayout processes a pending payout (typically called by admin/system)
func (s *PaymentService) ProcessPayout(ctx context.Context, payoutID uuid.UUID, approve bool, reason string) error {
	payout, err := s.paymentRepo.GetPayoutByID(ctx, payoutID)
	if err != nil {
		return fmt.Errorf("payout not found: %w", err)
	}

	if payout.Status != models.PayoutStatusPending {
		return errors.New("payout is not in pending status")
	}

	if !approve {
		return s.paymentRepo.UpdatePayoutStatus(ctx, payoutID, models.PayoutStatusCancelled, "", reason)
	}

	// Update to processing
	if err := s.paymentRepo.UpdatePayoutStatus(ctx, payoutID, models.PayoutStatusProcessing, "", ""); err != nil {
		return err
	}

	// Process based on payout method
	var transactionRef string
	var processErr error

	switch payout.PayoutMethod {
	case models.PayoutMethodWallet:
		transactionRef, processErr = s.processWalletPayout(ctx, payout)
	case models.PayoutMethodBankTransfer:
		transactionRef, processErr = s.processBankPayout(ctx, payout)
	case models.PayoutMethodMobileMoney:
		transactionRef, processErr = s.processMobileMoneyPayout(ctx, payout)
	default:
		processErr = fmt.Errorf("unsupported payout method: %s", payout.PayoutMethod)
	}

	if processErr != nil {
		s.paymentRepo.UpdatePayoutStatus(ctx, payoutID, models.PayoutStatusFailed, "", processErr.Error())
		return processErr
	}

	return s.paymentRepo.UpdatePayoutStatus(ctx, payoutID, models.PayoutStatusCompleted, transactionRef, "")
}

// processWalletPayout credits the courier's in-app wallet
func (s *PaymentService) processWalletPayout(ctx context.Context, payout *models.Payout) (string, error) {
	// Get current wallet balance
	wallet, err := s.paymentRepo.GetCourierWallet(ctx, payout.CourierID)
	if err != nil {
		return "", err
	}

	// Create wallet transaction
	tx := &models.WalletTransaction{
		ID:            uuid.New(),
		CourierID:     payout.CourierID,
		PayoutID:      &payout.ID,
		Type:          "payout",
		Amount:        payout.NetAmount,
		BalanceBefore: wallet.AvailableBalance,
		BalanceAfter:  wallet.AvailableBalance + payout.NetAmount,
		Description:   fmt.Sprintf("Payout for %d orders", len(payout.OrderIDs)),
		Reference:     payout.ID.String(),
	}

	if err := s.paymentRepo.CreateWalletTransaction(ctx, tx); err != nil {
		return "", err
	}

	// Update wallet balance
	if err := s.paymentRepo.UpdateWalletBalance(ctx, payout.CourierID, payout.NetAmount, 0, payout.NetAmount); err != nil {
		return "", err
	}

	return fmt.Sprintf("WALLET-%s", payout.ID.String()[:8]), nil
}

// processBankPayout initiates a bank transfer (placeholder for actual implementation)
func (s *PaymentService) processBankPayout(ctx context.Context, payout *models.Payout) (string, error) {
	// This would integrate with your bank transfer API
	// For now, return a placeholder reference
	return fmt.Sprintf("BANK-TXN-%s", payout.ID.String()[:8]), nil
}

// processMobileMoneyPayout sends money via mobile money (placeholder for actual implementation)
func (s *PaymentService) processMobileMoneyPayout(ctx context.Context, payout *models.Payout) (string, error) {
	// This would integrate with mobile money APIs (MTN, Airtel, Zamtel)
	// For now, return a placeholder reference
	return fmt.Sprintf("MOMO-TXN-%s", payout.ID.String()[:8]), nil
}

// ============================================================
// EARNINGS & WALLET
// ============================================================

// GetEarningsSummary retrieves earnings summary for a courier
func (s *PaymentService) GetEarningsSummary(ctx context.Context, courierID uuid.UUID) (*models.EarningsSummary, error) {
	summary, err := s.paymentRepo.GetEarningsSummary(ctx, courierID)
	if err != nil {
		return nil, err
	}

	// Format currency strings
	summary.FormattedTotal = fmt.Sprintf("%s %.2f", s.config.CurrencySymbol, summary.TotalEarnings)
	summary.FormattedAvailable = fmt.Sprintf("%s %.2f", s.config.CurrencySymbol, summary.AvailableBalance)
	summary.Currency = s.config.Currency

	return summary, nil
}

// GetWalletTransactions retrieves wallet transaction history
func (s *PaymentService) GetWalletTransactions(ctx context.Context, courierID uuid.UUID, page, pageSize int) ([]models.WalletTransaction, error) {
	offset := (page - 1) * pageSize
	return s.paymentRepo.GetWalletTransactions(ctx, courierID, pageSize, offset)
}

// CreditEarning credits earnings to courier wallet after delivery
func (s *PaymentService) CreditEarning(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusDelivered {
		return errors.New("order must be delivered to credit earnings")
	}

	// Get wallet
	wallet, err := s.paymentRepo.GetCourierWallet(ctx, order.CourierID)
	if err != nil {
		return err
	}

	// Create earning transaction
	tx := &models.WalletTransaction{
		ID:            uuid.New(),
		CourierID:     order.CourierID,
		OrderID:       &orderID,
		Type:          "earning",
		Amount:        order.CourierEarnings,
		BalanceBefore: wallet.PendingBalance,
		BalanceAfter:  wallet.PendingBalance + order.CourierEarnings,
		Description:   fmt.Sprintf("Earnings for order %s", order.OrderNumber),
		Reference:     order.OrderNumber,
	}

	if err := s.paymentRepo.CreateWalletTransaction(ctx, tx); err != nil {
		return err
	}

	// Add to pending balance (will move to available after payment verification)
	return s.paymentRepo.UpdateWalletBalance(ctx, order.CourierID, 0, order.CourierEarnings, 0)
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// updateOrderPaymentStatus updates the payment status of an order
func (s *PaymentService) updateOrderPaymentStatus(ctx context.Context, orderID uuid.UUID, status models.PaymentStatus, reference string) error {
	return s.orderRepo.UpdatePaymentStatus(ctx, orderID, string(status), reference)
}

// generateSignature creates HMAC signature for API requests
func (s *PaymentService) generateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// ============================================================
// REQUEST/RESPONSE TYPES FOR STORE INTEGRATION
// ============================================================

// StorePaymentVerifyRequest is sent to store's payment API
type StorePaymentVerifyRequest struct {
	OrderID         string  `json:"orderId"`
	ExternalOrderID string  `json:"externalOrderId"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

// StorePaymentVerifyResponse is received from store's payment API
type StorePaymentVerifyResponse struct {
	IsPaid         bool    `json:"isPaid"`
	AmountPaid     float64 `json:"amountPaid"`
	PaymentMethod  string  `json:"paymentMethod"`
	TransactionRef string  `json:"transactionRef"`
	PaidAt         string  `json:"paidAt"`
	Error          string  `json:"error,omitempty"`
}

// ============================================================
// CUSTOM STORE VERIFIER EXAMPLE
// ============================================================

// GenericStoreVerifier is a configurable verifier for stores
type GenericStoreVerifier struct {
	StoreName     string
	APIURL        string
	APIKey        string
	WebhookSecret string
	httpClient    *http.Client
}

// NewGenericStoreVerifier creates a new generic store verifier
func NewGenericStoreVerifier(name, apiURL, apiKey, secret string) *GenericStoreVerifier {
	return &GenericStoreVerifier{
		StoreName:     name,
		APIURL:        apiURL,
		APIKey:        apiKey,
		WebhookSecret: secret,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
	}
}

// VerifyPayment implements PaymentVerifier interface
func (v *GenericStoreVerifier) VerifyPayment(ctx context.Context, orderID, externalOrderID string) (*models.PaymentVerification, error) {
	parsedID, _ := uuid.Parse(orderID)

	reqBody, _ := json.Marshal(map[string]any{
		"order_id":          orderID,
		"external_order_id": externalOrderID,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.APIURL+"/verify-payment", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", v.APIKey)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return &models.PaymentVerification{
			OrderID: parsedID,
			IsPaid:  false,
			Error:   err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	var result StorePaymentVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	verification := &models.PaymentVerification{
		OrderID:       parsedID,
		IsPaid:        result.IsPaid,
		AmountPaid:    result.AmountPaid,
		PaymentMethod: result.PaymentMethod,
		PaymentRef:    result.TransactionRef,
	}

	if result.PaidAt != "" {
		paidAt, _ := time.Parse(time.RFC3339, result.PaidAt)
		verification.PaidAt = &paidAt
	}

	return verification, nil
}

// GetProviderName returns the store name
func (v *GenericStoreVerifier) GetProviderName() string {
	return v.StoreName
}
