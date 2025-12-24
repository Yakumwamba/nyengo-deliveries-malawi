package models

import (
	"time"

	"github.com/google/uuid"
)

// PayoutStatus represents the status of a payout request
type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
	PayoutStatusCancelled  PayoutStatus = "cancelled"
)

// PayoutMethod represents how the courier wants to receive payment
type PayoutMethod string

const (
	PayoutMethodBankTransfer PayoutMethod = "bank_transfer"
	PayoutMethodMobileMoney  PayoutMethod = "mobile_money"
	PayoutMethodWallet       PayoutMethod = "wallet" // Keep in app wallet
)

// Payout represents a courier's payout request
type Payout struct {
	ID             uuid.UUID    `json:"id"`
	CourierID      uuid.UUID    `json:"courierId"`
	OrderIDs       []uuid.UUID  `json:"orderIds"`    // Orders included in this payout
	TotalAmount    float64      `json:"totalAmount"` // Total payout amount
	PlatformFee    float64      `json:"platformFee"` // Deducted platform fee
	NetAmount      float64      `json:"netAmount"`   // Amount after fees
	Currency       string       `json:"currency"`
	Status         PayoutStatus `json:"status"`
	PayoutMethod   PayoutMethod `json:"payoutMethod"`
	PayoutDetails  string       `json:"payoutDetails,omitempty"` // Bank/mobile money details
	TransactionRef string       `json:"transactionRef,omitempty"`
	FailureReason  string       `json:"failureReason,omitempty"`
	ProcessedAt    *time.Time   `json:"processedAt,omitempty"`
	CompletedAt    *time.Time   `json:"completedAt,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

// PayoutRequest is the DTO for requesting a payout
type PayoutRequest struct {
	OrderIDs      []string     `json:"orderIds"` // Orders to include
	PayoutMethod  PayoutMethod `json:"payoutMethod"`
	PayoutDetails string       `json:"payoutDetails,omitempty"` // Bank account or mobile number
}

// PayoutResponse is returned after creating a payout
type PayoutResponse struct {
	PayoutID    uuid.UUID    `json:"payoutId"`
	Status      PayoutStatus `json:"status"`
	TotalAmount string       `json:"totalAmount"`
	NetAmount   string       `json:"netAmount"`
	Message     string       `json:"message"`
}

// CourierWallet represents a courier's in-app wallet
type CourierWallet struct {
	CourierID        uuid.UUID `json:"courierId"`
	AvailableBalance float64   `json:"availableBalance"` // Can be withdrawn
	PendingBalance   float64   `json:"pendingBalance"`   // Pending verification
	TotalEarnings    float64   `json:"totalEarnings"`    // All-time earnings
	Currency         string    `json:"currency"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// WalletTransaction represents a transaction in the wallet
type WalletTransaction struct {
	ID            uuid.UUID  `json:"id"`
	CourierID     uuid.UUID  `json:"courierId"`
	OrderID       *uuid.UUID `json:"orderId,omitempty"`
	PayoutID      *uuid.UUID `json:"payoutId,omitempty"`
	Type          string     `json:"type"` // "earning", "payout", "adjustment", "refund"
	Amount        float64    `json:"amount"`
	BalanceBefore float64    `json:"balanceBefore"`
	BalanceAfter  float64    `json:"balanceAfter"`
	Description   string     `json:"description"`
	Reference     string     `json:"reference,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

// PayableOrder represents an order eligible for payout
type PayableOrder struct {
	OrderID         uuid.UUID  `json:"orderId"`
	OrderNumber     string     `json:"orderNumber"`
	DeliveredAt     time.Time  `json:"deliveredAt"`
	TotalFare       float64    `json:"totalFare"`
	CourierEarnings float64    `json:"courierEarnings"`
	PaymentStatus   string     `json:"paymentStatus"`
	PayoutStatus    string     `json:"payoutStatus"` // "unpaid", "pending", "paid"
	StoreID         *uuid.UUID `json:"storeId,omitempty"`
	CustomerName    string     `json:"customerName"`
}

// PaymentVerification holds verification result from store/payment gateway
type PaymentVerification struct {
	OrderID       uuid.UUID  `json:"orderId"`
	IsPaid        bool       `json:"isPaid"`
	AmountPaid    float64    `json:"amountPaid"`
	PaymentMethod string     `json:"paymentMethod"`
	PaymentRef    string     `json:"paymentRef"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	Error         string     `json:"error,omitempty"`
}

// StorePaymentConfig holds configuration for store payment integration
type StorePaymentConfig struct {
	StoreID       uuid.UUID `json:"storeId"`
	StoreName     string    `json:"storeName"`
	PaymentAPIURL string    `json:"paymentApiUrl"`
	APIKey        string    `json:"apiKey"`
	WebhookSecret string    `json:"webhookSecret"`
	IsActive      bool      `json:"isActive"`
}

// BankDetails is defined in courier.go

// MobileMoneyDetails holds mobile money information
type MobileMoneyDetails struct {
	Provider    string `json:"provider"` // "mtn", "airtel", "zamtel"
	PhoneNumber string `json:"phoneNumber"`
	AccountName string `json:"accountName"`
}

// EarningsSummary provides a summary of courier earnings
type EarningsSummary struct {
	TotalEarnings      float64 `json:"totalEarnings"`
	AvailableBalance   float64 `json:"availableBalance"`
	PendingBalance     float64 `json:"pendingBalance"`
	TotalPaidOut       float64 `json:"totalPaidOut"`
	UnpaidOrders       int     `json:"unpaidOrders"`
	PendingPayouts     int     `json:"pendingPayouts"`
	Currency           string  `json:"currency"`
	FormattedTotal     string  `json:"formattedTotal"`
	FormattedAvailable string  `json:"formattedAvailable"`
}
