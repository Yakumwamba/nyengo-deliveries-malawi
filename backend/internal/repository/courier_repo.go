package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"nyengo-deliveries/internal/models"
)

// CourierRepository handles courier data access
type CourierRepository struct {
	db *pgxpool.Pool
}

// NewCourierRepository creates a new courier repository
func NewCourierRepository(db *pgxpool.Pool) *CourierRepository {
	return &CourierRepository{db: db}
}

// Create inserts a new courier into the database
func (r *CourierRepository) Create(ctx context.Context, courier *models.Courier) error {
	query := `
		INSERT INTO couriers (
			id, email, password, company_name, owner_name, phone, address, city, country,
			service_areas, vehicle_types, base_rate_per_km, minimum_fare, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	courier.ID = uuid.New()
	courier.CreatedAt = time.Now()
	courier.UpdatedAt = time.Now()
	courier.IsActive = true

	_, err := r.db.Exec(ctx, query,
		courier.ID,
		courier.Email,
		courier.Password,
		courier.CompanyName,
		courier.OwnerName,
		courier.Phone,
		courier.Address,
		courier.City,
		courier.Country,
		courier.ServiceAreas,
		courier.VehicleTypes,
		courier.BaseRatePerKm,
		courier.MinimumFare,
		courier.IsActive,
		courier.CreatedAt,
		courier.UpdatedAt,
	)

	return err
}

// GetByID retrieves a courier by ID
func (r *CourierRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
	query := `
		SELECT id, email, company_name, owner_name, phone, alternate_phone, whatsapp,
			address, city, country, logo_url, description, service_areas, vehicle_types,
			max_weight, base_rate_per_km, minimum_fare, custom_pricing, rating, total_reviews,
			total_deliveries, success_rate, is_verified, is_active, is_featured, wallet_balance,
			created_at, updated_at, last_active_at
		FROM couriers WHERE id = $1
	`

	var courier models.Courier
	err := r.db.QueryRow(ctx, query, id).Scan(
		&courier.ID,
		&courier.Email,
		&courier.CompanyName,
		&courier.OwnerName,
		&courier.Phone,
		&courier.AlternatePhone,
		&courier.WhatsApp,
		&courier.Address,
		&courier.City,
		&courier.Country,
		&courier.LogoURL,
		&courier.Description,
		&courier.ServiceAreas,
		&courier.VehicleTypes,
		&courier.MaxWeight,
		&courier.BaseRatePerKm,
		&courier.MinimumFare,
		&courier.CustomPricing,
		&courier.Rating,
		&courier.TotalReviews,
		&courier.TotalDeliveries,
		&courier.SuccessRate,
		&courier.IsVerified,
		&courier.IsActive,
		&courier.IsFeatured,
		&courier.WalletBalance,
		&courier.CreatedAt,
		&courier.UpdatedAt,
		&courier.LastActiveAt,
	)

	if err != nil {
		return nil, err
	}

	return &courier, nil
}

// GetByEmail retrieves a courier by email
func (r *CourierRepository) GetByEmail(ctx context.Context, email string) (*models.Courier, error) {
	query := `
		SELECT id, email, password, company_name, owner_name, phone, is_active, created_at
		FROM couriers WHERE email = $1
	`

	var courier models.Courier
	err := r.db.QueryRow(ctx, query, email).Scan(
		&courier.ID,
		&courier.Email,
		&courier.Password,
		&courier.CompanyName,
		&courier.OwnerName,
		&courier.Phone,
		&courier.IsActive,
		&courier.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &courier, nil
}

// Update updates a courier's information
func (r *CourierRepository) Update(ctx context.Context, courier *models.Courier) error {
	query := `
		UPDATE couriers SET
			company_name = $2, owner_name = $3, phone = $4, alternate_phone = $5,
			whatsapp = $6, address = $7, city = $8, logo_url = $9, description = $10,
			service_areas = $11, vehicle_types = $12, max_weight = $13, base_rate_per_km = $14,
			minimum_fare = $15, updated_at = $16
		WHERE id = $1
	`

	courier.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		courier.ID,
		courier.CompanyName,
		courier.OwnerName,
		courier.Phone,
		courier.AlternatePhone,
		courier.WhatsApp,
		courier.Address,
		courier.City,
		courier.LogoURL,
		courier.Description,
		courier.ServiceAreas,
		courier.VehicleTypes,
		courier.MaxWeight,
		courier.BaseRatePerKm,
		courier.MinimumFare,
		courier.UpdatedAt,
	)

	return err
}

// ListActive retrieves all active couriers
func (r *CourierRepository) ListActive(ctx context.Context) ([]models.CourierListItem, error) {
	query := `
		SELECT id, company_name, logo_url, rating, total_reviews, total_deliveries,
			base_rate_per_km, minimum_fare, is_verified, is_featured
		FROM couriers
		WHERE is_active = true
		ORDER BY is_featured DESC, rating DESC, total_deliveries DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []models.CourierListItem
	for rows.Next() {
		var c models.CourierListItem
		err := rows.Scan(
			&c.ID,
			&c.CompanyName,
			&c.LogoURL,
			&c.Rating,
			&c.TotalReviews,
			&c.TotalDeliveries,
			&c.BaseRatePerKm,
			&c.MinimumFare,
			&c.IsVerified,
			&c.IsFeatured,
		)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}

	return couriers, nil
}

// ListByArea retrieves couriers that serve a specific area
func (r *CourierRepository) ListByArea(ctx context.Context, area string) ([]models.CourierListItem, error) {
	query := `
		SELECT id, company_name, logo_url, rating, total_reviews, total_deliveries,
			base_rate_per_km, minimum_fare, is_verified, is_featured
		FROM couriers
		WHERE is_active = true AND $1 = ANY(service_areas)
		ORDER BY is_featured DESC, rating DESC
	`

	rows, err := r.db.Query(ctx, query, area)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []models.CourierListItem
	for rows.Next() {
		var c models.CourierListItem
		err := rows.Scan(
			&c.ID,
			&c.CompanyName,
			&c.LogoURL,
			&c.Rating,
			&c.TotalReviews,
			&c.TotalDeliveries,
			&c.BaseRatePerKm,
			&c.MinimumFare,
			&c.IsVerified,
			&c.IsFeatured,
		)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}

	return couriers, nil
}

// UpdateLastActive updates the courier's last active timestamp
func (r *CourierRepository) UpdateLastActive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE couriers SET last_active_at = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, time.Now())
	return err
}

// UpdateStats updates courier statistics
func (r *CourierRepository) UpdateStats(ctx context.Context, id uuid.UUID, totalDeliveries int, successRate, rating float64) error {
	query := `
		UPDATE couriers SET
			total_deliveries = $2,
			success_rate = $3,
			rating = $4,
			updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, totalDeliveries, successRate, rating, time.Now())
	return err
}

// UpdateWalletBalance updates the courier's wallet balance
func (r *CourierRepository) UpdateWalletBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	query := `UPDATE couriers SET wallet_balance = wallet_balance + $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, amount, time.Now())
	return err
}

// EmailExists checks if an email already exists
func (r *CourierRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM couriers WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
