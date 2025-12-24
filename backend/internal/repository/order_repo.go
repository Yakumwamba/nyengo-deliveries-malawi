package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"nyengo-deliveries/internal/models"
)

// OrderRepository handles order data access
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts a new order into the database
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (
			id, order_number, courier_id, store_id, external_order_id,
			customer_name, customer_phone, customer_email,
			pickup_address, pickup_latitude, pickup_longitude, pickup_notes,
			pickup_contact_name, pickup_contact_phone,
			delivery_address, delivery_latitude, delivery_longitude, delivery_notes,
			package_description, package_size, package_weight, is_fragile, requires_signature,
			distance, base_fare, distance_fare, surge_fare, total_fare, platform_fee, courier_earnings,
			payment_method, payment_status, status, scheduled_pickup,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36
		)
	`

	order.ID = uuid.New()
	order.OrderNumber = generateOrderNumber()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = models.OrderStatusPending
	order.PaymentStatus = models.PaymentStatusPending

	_, err := r.db.Exec(ctx, query,
		order.ID,
		order.OrderNumber,
		order.CourierID,
		order.StoreID,
		order.ExternalOrderID,
		order.CustomerName,
		order.CustomerPhone,
		order.CustomerEmail,
		order.PickupAddress,
		order.PickupLatitude,
		order.PickupLongitude,
		order.PickupNotes,
		order.PickupContactName,
		order.PickupContactPhone,
		order.DeliveryAddress,
		order.DeliveryLatitude,
		order.DeliveryLongitude,
		order.DeliveryNotes,
		order.PackageDescription,
		order.PackageSize,
		order.PackageWeight,
		order.IsFragile,
		order.RequiresSignature,
		order.Distance,
		order.BaseFare,
		order.DistanceFare,
		order.SurgeFare,
		order.TotalFare,
		order.PlatformFee,
		order.CourierEarnings,
		order.PaymentMethod,
		order.PaymentStatus,
		order.Status,
		order.ScheduledPickup,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query := `
		SELECT id, order_number, courier_id, store_id, external_order_id,
			customer_name, customer_phone, customer_email,
			pickup_address, pickup_latitude, pickup_longitude, pickup_notes,
			pickup_contact_name, pickup_contact_phone,
			delivery_address, delivery_latitude, delivery_longitude, delivery_notes,
			package_description, package_size, package_weight, is_fragile, requires_signature,
			distance, base_fare, distance_fare, surge_fare, total_fare, platform_fee, courier_earnings,
			payment_method, payment_status, payment_reference, status,
			scheduled_pickup, actual_pickup, estimated_delivery, actual_delivery,
			delivery_proof_url, recipient_name, signature_url,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM orders WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CourierID,
		&order.StoreID,
		&order.ExternalOrderID,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CustomerEmail,
		&order.PickupAddress,
		&order.PickupLatitude,
		&order.PickupLongitude,
		&order.PickupNotes,
		&order.PickupContactName,
		&order.PickupContactPhone,
		&order.DeliveryAddress,
		&order.DeliveryLatitude,
		&order.DeliveryLongitude,
		&order.DeliveryNotes,
		&order.PackageDescription,
		&order.PackageSize,
		&order.PackageWeight,
		&order.IsFragile,
		&order.RequiresSignature,
		&order.Distance,
		&order.BaseFare,
		&order.DistanceFare,
		&order.SurgeFare,
		&order.TotalFare,
		&order.PlatformFee,
		&order.CourierEarnings,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&order.PaymentReference,
		&order.Status,
		&order.ScheduledPickup,
		&order.ActualPickup,
		&order.EstimatedDelivery,
		&order.ActualDelivery,
		&order.DeliveryProofURL,
		&order.RecipientName,
		&order.SignatureURL,
		&order.CustomerRating,
		&order.CustomerFeedback,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	query := `SELECT id FROM orders WHERE order_number = $1`
	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, orderNumber).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

// List retrieves orders for a courier with filters
func (r *OrderRepository) List(ctx context.Context, courierID uuid.UUID, filters *models.OrderListFilters) (*models.OrderListResponse, error) {
	// Build query dynamically
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("courier_id = $%d", argIndex))
	args = append(args, courierID)
	argIndex++

	if len(filters.Status) > 0 {
		placeholders := make([]string, len(filters.Status))
		for i, s := range filters.Status {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, s)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(customer_name ILIKE $%d OR order_number ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE %s", whereClause)
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Set defaults for pagination
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	// Sort
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filters.Page - 1) * filters.PageSize

	// Main query
	query := fmt.Sprintf(`
		SELECT id, order_number, customer_name, customer_phone,
			pickup_address, delivery_address, package_size,
			distance, total_fare, status, payment_status, created_at
		FROM orders
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortOrder, argIndex, argIndex+1)

	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.OrderNumber,
			&o.CustomerName,
			&o.CustomerPhone,
			&o.PickupAddress,
			&o.DeliveryAddress,
			&o.PackageSize,
			&o.Distance,
			&o.TotalFare,
			&o.Status,
			&o.PaymentStatus,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize

	return &models.OrderListResponse{
		Orders:     orders,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateStatus updates the order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	query := `UPDATE orders SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status, time.Now())
	return err
}

// UpdatePaymentStatus updates the payment status and reference
func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string, reference string) error {
	query := `UPDATE orders SET payment_status = $2, payment_reference = $3, updated_at = $4 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status, reference, time.Now())
	return err
}

// UpdateDeliveryProof updates delivery proof information
func (r *OrderRepository) UpdateDeliveryProof(ctx context.Context, id uuid.UUID, proofURL, recipientName, signatureURL string) error {
	query := `
		UPDATE orders SET
			delivery_proof_url = $2,
			recipient_name = $3,
			signature_url = $4,
			actual_delivery = $5,
			updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, proofURL, recipientName, signatureURL, time.Now())
	return err
}

// GetDailyStats retrieves daily statistics for a courier
func (r *OrderRepository) GetDailyStats(ctx context.Context, courierID uuid.UUID, date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT
			COUNT(*) as total_orders,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN status = 'pending' OR status = 'accepted' OR status = 'in_transit' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(total_fare), 0) as total_revenue,
			COALESCE(SUM(courier_earnings), 0) as total_earnings,
			COALESCE(SUM(platform_fee), 0) as total_fees
		FROM orders
		WHERE courier_id = $1 AND created_at >= $2 AND created_at < $3
	`

	var stats struct {
		TotalOrders   int
		Completed     int
		Pending       int
		TotalRevenue  float64
		TotalEarnings float64
		TotalFees     float64
	}

	err := r.db.QueryRow(ctx, query, courierID, startOfDay, endOfDay).Scan(
		&stats.TotalOrders,
		&stats.Completed,
		&stats.Pending,
		&stats.TotalRevenue,
		&stats.TotalEarnings,
		&stats.TotalFees,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalOrders":   stats.TotalOrders,
		"completed":     stats.Completed,
		"pending":       stats.Pending,
		"totalRevenue":  stats.TotalRevenue,
		"totalEarnings": stats.TotalEarnings,
		"totalFees":     stats.TotalFees,
	}, nil
}

// GetMonthlyStats retrieves monthly statistics
func (r *OrderRepository) GetMonthlyStats(ctx context.Context, courierID uuid.UUID, year, month int) (map[string]interface{}, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	query := `
		SELECT
			COUNT(*) as total_orders,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(total_fare), 0) as total_revenue,
			COALESCE(SUM(courier_earnings), 0) as total_earnings,
			COALESCE(AVG(customer_rating), 0) as avg_rating
		FROM orders
		WHERE courier_id = $1 AND created_at >= $2 AND created_at < $3
	`

	var stats struct {
		TotalOrders   int
		Completed     int
		TotalRevenue  float64
		TotalEarnings float64
		AvgRating     float64
	}

	err := r.db.QueryRow(ctx, query, courierID, startOfMonth, endOfMonth).Scan(
		&stats.TotalOrders,
		&stats.Completed,
		&stats.TotalRevenue,
		&stats.TotalEarnings,
		&stats.AvgRating,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalOrders":   stats.TotalOrders,
		"completed":     stats.Completed,
		"totalRevenue":  stats.TotalRevenue,
		"totalEarnings": stats.TotalEarnings,
		"avgRating":     stats.AvgRating,
	}, nil
}

// generateOrderNumber creates a unique order number
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102")
	unique := uuid.New().String()[:8]
	return fmt.Sprintf("NYG-%s-%s", timestamp, strings.ToUpper(unique))
}
