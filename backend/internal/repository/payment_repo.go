package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"nyengo-deliveries/internal/models"
)

// PaymentRepository handles payment and payout database operations
type PaymentRepository struct {
	db *pgxpool.Pool
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayout creates a new payout request
func (r *PaymentRepository) CreatePayout(ctx context.Context, payout *models.Payout) error {
	orderIDsJSON, _ := json.Marshal(payout.OrderIDs)

	query := `
		INSERT INTO payouts (
			id, courier_id, order_ids, total_amount, platform_fee, net_amount,
			currency, status, payout_method, payout_details, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
	`

	_, err := r.db.Exec(ctx, query,
		payout.ID,
		payout.CourierID,
		orderIDsJSON,
		payout.TotalAmount,
		payout.PlatformFee,
		payout.NetAmount,
		payout.Currency,
		payout.Status,
		payout.PayoutMethod,
		payout.PayoutDetails,
		time.Now(),
	)

	return err
}

// GetPayoutByID retrieves a payout by ID
func (r *PaymentRepository) GetPayoutByID(ctx context.Context, id uuid.UUID) (*models.Payout, error) {
	query := `
		SELECT id, courier_id, order_ids, total_amount, platform_fee, net_amount,
			   currency, status, payout_method, payout_details, transaction_ref,
			   failure_reason, processed_at, completed_at, created_at, updated_at
		FROM payouts WHERE id = $1
	`

	var payout models.Payout
	var orderIDsJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&payout.ID,
		&payout.CourierID,
		&orderIDsJSON,
		&payout.TotalAmount,
		&payout.PlatformFee,
		&payout.NetAmount,
		&payout.Currency,
		&payout.Status,
		&payout.PayoutMethod,
		&payout.PayoutDetails,
		&payout.TransactionRef,
		&payout.FailureReason,
		&payout.ProcessedAt,
		&payout.CompletedAt,
		&payout.CreatedAt,
		&payout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(orderIDsJSON, &payout.OrderIDs)
	return &payout, nil
}

// UpdatePayoutStatus updates the status of a payout
func (r *PaymentRepository) UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status models.PayoutStatus, transactionRef, failureReason string) error {
	query := `
		UPDATE payouts SET 
			status = $2, 
			transaction_ref = COALESCE(NULLIF($3, ''), transaction_ref),
			failure_reason = COALESCE(NULLIF($4, ''), failure_reason),
			processed_at = CASE WHEN $2 = 'processing' THEN NOW() ELSE processed_at END,
			completed_at = CASE WHEN $2 IN ('completed', 'failed') THEN NOW() ELSE completed_at END,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, payoutID, status, transactionRef, failureReason)
	return err
}

// GetCourierPayouts retrieves payouts for a courier
func (r *PaymentRepository) GetCourierPayouts(ctx context.Context, courierID uuid.UUID, status string, limit, offset int) ([]models.Payout, error) {
	query := `
		SELECT id, courier_id, order_ids, total_amount, platform_fee, net_amount,
			   currency, status, payout_method, payout_details, transaction_ref,
			   failure_reason, processed_at, completed_at, created_at, updated_at
		FROM payouts 
		WHERE courier_id = $1
		AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, courierID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payouts []models.Payout
	for rows.Next() {
		var payout models.Payout
		var orderIDsJSON []byte

		if err := rows.Scan(
			&payout.ID,
			&payout.CourierID,
			&orderIDsJSON,
			&payout.TotalAmount,
			&payout.PlatformFee,
			&payout.NetAmount,
			&payout.Currency,
			&payout.Status,
			&payout.PayoutMethod,
			&payout.PayoutDetails,
			&payout.TransactionRef,
			&payout.FailureReason,
			&payout.ProcessedAt,
			&payout.CompletedAt,
			&payout.CreatedAt,
			&payout.UpdatedAt,
		); err != nil {
			return nil, err
		}

		json.Unmarshal(orderIDsJSON, &payout.OrderIDs)
		payouts = append(payouts, payout)
	}

	return payouts, nil
}

// GetPayableOrders retrieves orders eligible for payout
func (r *PaymentRepository) GetPayableOrders(ctx context.Context, courierID uuid.UUID) ([]models.PayableOrder, error) {
	query := `
		SELECT o.id, o.order_number, o.actual_delivery, o.total_fare, o.courier_earnings,
			   o.payment_status, o.store_id, o.customer_name,
			   CASE 
				   WHEN EXISTS (SELECT 1 FROM payout_orders po JOIN payouts p ON po.payout_id = p.id 
							   WHERE po.order_id = o.id AND p.status IN ('pending', 'processing', 'completed'))
				   THEN 'paid'
				   WHEN EXISTS (SELECT 1 FROM payout_orders po JOIN payouts p ON po.payout_id = p.id 
							   WHERE po.order_id = o.id AND p.status = 'pending')
				   THEN 'pending'
				   ELSE 'unpaid'
			   END as payout_status
		FROM orders o
		WHERE o.courier_id = $1
		AND o.status = 'delivered'
		AND o.payment_status = 'paid'
		AND NOT EXISTS (
			SELECT 1 FROM payout_orders po 
			JOIN payouts p ON po.payout_id = p.id 
			WHERE po.order_id = o.id 
			AND p.status IN ('completed', 'processing')
		)
		ORDER BY o.actual_delivery DESC
	`

	rows, err := r.db.Query(ctx, query, courierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.PayableOrder
	for rows.Next() {
		var order models.PayableOrder
		if err := rows.Scan(
			&order.OrderID,
			&order.OrderNumber,
			&order.DeliveredAt,
			&order.TotalFare,
			&order.CourierEarnings,
			&order.PaymentStatus,
			&order.StoreID,
			&order.CustomerName,
			&order.PayoutStatus,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// LinkOrdersToPayout associates orders with a payout
func (r *PaymentRepository) LinkOrdersToPayout(ctx context.Context, payoutID uuid.UUID, orderIDs []uuid.UUID) error {
	for _, orderID := range orderIDs {
		query := `INSERT INTO payout_orders (payout_id, order_id) VALUES ($1, $2)`
		if _, err := r.db.Exec(ctx, query, payoutID, orderID); err != nil {
			return err
		}
	}
	return nil
}

// GetCourierWallet retrieves or creates a courier's wallet
func (r *PaymentRepository) GetCourierWallet(ctx context.Context, courierID uuid.UUID) (*models.CourierWallet, error) {
	query := `
		SELECT courier_id, available_balance, pending_balance, total_earnings, currency, updated_at
		FROM courier_wallets WHERE courier_id = $1
	`

	var wallet models.CourierWallet
	err := r.db.QueryRow(ctx, query, courierID).Scan(
		&wallet.CourierID,
		&wallet.AvailableBalance,
		&wallet.PendingBalance,
		&wallet.TotalEarnings,
		&wallet.Currency,
		&wallet.UpdatedAt,
	)

	if err != nil {
		// Create wallet if doesn't exist
		createQuery := `
			INSERT INTO courier_wallets (courier_id, available_balance, pending_balance, total_earnings, currency, updated_at)
			VALUES ($1, 0, 0, 0, 'ZMW', NOW())
			ON CONFLICT (courier_id) DO NOTHING
			RETURNING courier_id, available_balance, pending_balance, total_earnings, currency, updated_at
		`
		err = r.db.QueryRow(ctx, createQuery, courierID).Scan(
			&wallet.CourierID,
			&wallet.AvailableBalance,
			&wallet.PendingBalance,
			&wallet.TotalEarnings,
			&wallet.Currency,
			&wallet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &wallet, nil
}

// UpdateWalletBalance updates the courier's wallet balance
func (r *PaymentRepository) UpdateWalletBalance(ctx context.Context, courierID uuid.UUID, availableDelta, pendingDelta, totalDelta float64) error {
	query := `
		UPDATE courier_wallets SET
			available_balance = available_balance + $2,
			pending_balance = pending_balance + $3,
			total_earnings = total_earnings + $4,
			updated_at = NOW()
		WHERE courier_id = $1
	`

	_, err := r.db.Exec(ctx, query, courierID, availableDelta, pendingDelta, totalDelta)
	return err
}

// CreateWalletTransaction records a wallet transaction
func (r *PaymentRepository) CreateWalletTransaction(ctx context.Context, tx *models.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (
			id, courier_id, order_id, payout_id, type, amount,
			balance_before, balance_after, description, reference, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		tx.ID,
		tx.CourierID,
		tx.OrderID,
		tx.PayoutID,
		tx.Type,
		tx.Amount,
		tx.BalanceBefore,
		tx.BalanceAfter,
		tx.Description,
		tx.Reference,
		time.Now(),
	)

	return err
}

// GetWalletTransactions retrieves wallet transaction history
func (r *PaymentRepository) GetWalletTransactions(ctx context.Context, courierID uuid.UUID, limit, offset int) ([]models.WalletTransaction, error) {
	query := `
		SELECT id, courier_id, order_id, payout_id, type, amount,
			   balance_before, balance_after, description, reference, created_at
		FROM wallet_transactions
		WHERE courier_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, courierID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.WalletTransaction
	for rows.Next() {
		var tx models.WalletTransaction
		if err := rows.Scan(
			&tx.ID,
			&tx.CourierID,
			&tx.OrderID,
			&tx.PayoutID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceBefore,
			&tx.BalanceAfter,
			&tx.Description,
			&tx.Reference,
			&tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// GetStorePaymentConfig retrieves payment config for a store
func (r *PaymentRepository) GetStorePaymentConfig(ctx context.Context, storeID uuid.UUID) (*models.StorePaymentConfig, error) {
	query := `
		SELECT store_id, store_name, payment_api_url, api_key, webhook_secret, is_active
		FROM store_payment_configs WHERE store_id = $1 AND is_active = true
	`

	var config models.StorePaymentConfig
	err := r.db.QueryRow(ctx, query, storeID).Scan(
		&config.StoreID,
		&config.StoreName,
		&config.PaymentAPIURL,
		&config.APIKey,
		&config.WebhookSecret,
		&config.IsActive,
	)

	if err != nil {
		return nil, fmt.Errorf("store payment config not found: %w", err)
	}

	return &config, nil
}

// GetEarningsSummary calculates earnings summary for a courier
func (r *PaymentRepository) GetEarningsSummary(ctx context.Context, courierID uuid.UUID) (*models.EarningsSummary, error) {
	wallet, err := r.GetCourierWallet(ctx, courierID)
	if err != nil {
		return nil, err
	}

	// Count unpaid orders
	var unpaidCount int
	unpaidQuery := `
		SELECT COUNT(*) FROM orders 
		WHERE courier_id = $1 AND status = 'delivered' AND payment_status = 'paid'
		AND NOT EXISTS (
			SELECT 1 FROM payout_orders po JOIN payouts p ON po.payout_id = p.id
			WHERE po.order_id = orders.id AND p.status IN ('completed', 'processing')
		)
	`
	r.db.QueryRow(ctx, unpaidQuery, courierID).Scan(&unpaidCount)

	// Count pending payouts
	var pendingCount int
	pendingQuery := `SELECT COUNT(*) FROM payouts WHERE courier_id = $1 AND status IN ('pending', 'processing')`
	r.db.QueryRow(ctx, pendingQuery, courierID).Scan(&pendingCount)

	// Total paid out
	var totalPaidOut float64
	paidQuery := `SELECT COALESCE(SUM(net_amount), 0) FROM payouts WHERE courier_id = $1 AND status = 'completed'`
	r.db.QueryRow(ctx, paidQuery, courierID).Scan(&totalPaidOut)

	return &models.EarningsSummary{
		TotalEarnings:    wallet.TotalEarnings,
		AvailableBalance: wallet.AvailableBalance,
		PendingBalance:   wallet.PendingBalance,
		TotalPaidOut:     totalPaidOut,
		UnpaidOrders:     unpaidCount,
		PendingPayouts:   pendingCount,
		Currency:         wallet.Currency,
	}, nil
}
