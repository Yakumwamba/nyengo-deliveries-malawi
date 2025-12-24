package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/repository"
)

// CourierService handles courier business logic
type CourierService struct {
	repo      *repository.CourierRepository
	jwtSecret string
}

// NewCourierService creates a new courier service
func NewCourierService(repo *repository.CourierRepository) *CourierService {
	return &CourierService{
		repo:      repo,
		jwtSecret: "your-super-secret-jwt-key-change-in-production", // Should come from config
	}
}

// SetJWTSecret sets the JWT secret (called from main)
func (s *CourierService) SetJWTSecret(secret string) {
	s.jwtSecret = secret
}

// Register creates a new courier account
func (s *CourierService) Register(ctx context.Context, req *models.CourierRegistrationRequest) (*models.Courier, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create courier
	courier := &models.Courier{
		Email:         req.Email,
		Password:      string(hashedPassword),
		CompanyName:   req.CompanyName,
		OwnerName:     req.OwnerName,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		Country:       req.Country,
		ServiceAreas:  req.ServiceAreas,
		VehicleTypes:  req.VehicleTypes,
		BaseRatePerKm: 5.0,  // Default rate
		MinimumFare:   20.0, // Default minimum
		Rating:        0.0,
		IsActive:      true,
		IsVerified:    false,
	}

	if err := s.repo.Create(ctx, courier); err != nil {
		return nil, err
	}

	// Clear password before returning
	courier.Password = ""

	return courier, nil
}

// Login authenticates a courier and returns a JWT token
func (s *CourierService) Login(ctx context.Context, req *models.CourierLoginRequest) (*models.CourierLoginResponse, error) {
	// Get courier by email
	courier, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if account is active
	if !courier.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(courier.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(courier.ID)
	if err != nil {
		return nil, err
	}

	// Update last active
	_ = s.repo.UpdateLastActive(ctx, courier.ID)

	// Clear password before returning
	courier.Password = ""

	return &models.CourierLoginResponse{
		Token:   token,
		Courier: courier,
	}, nil
}

// GetProfile retrieves the courier's full profile
func (s *CourierService) GetProfile(ctx context.Context, courierID uuid.UUID) (*models.Courier, error) {
	courier, err := s.repo.GetByID(ctx, courierID)
	if err != nil {
		return nil, err
	}

	// Update last active
	_ = s.repo.UpdateLastActive(ctx, courierID)

	return courier, nil
}

// UpdateProfile updates courier profile information
func (s *CourierService) UpdateProfile(ctx context.Context, courierID uuid.UUID, req *models.CourierUpdateRequest) (*models.Courier, error) {
	courier, err := s.repo.GetByID(ctx, courierID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.CompanyName != nil {
		courier.CompanyName = *req.CompanyName
	}
	if req.OwnerName != nil {
		courier.OwnerName = *req.OwnerName
	}
	if req.Phone != nil {
		courier.Phone = *req.Phone
	}
	if req.AlternatePhone != nil {
		courier.AlternatePhone = *req.AlternatePhone
	}
	if req.WhatsApp != nil {
		courier.WhatsApp = *req.WhatsApp
	}
	if req.Address != nil {
		courier.Address = *req.Address
	}
	if req.City != nil {
		courier.City = *req.City
	}
	if req.LogoURL != nil {
		courier.LogoURL = *req.LogoURL
	}
	if req.Description != nil {
		courier.Description = *req.Description
	}
	if len(req.ServiceAreas) > 0 {
		courier.ServiceAreas = req.ServiceAreas
	}
	if len(req.VehicleTypes) > 0 {
		courier.VehicleTypes = req.VehicleTypes
	}
	if req.MaxWeight != nil {
		courier.MaxWeight = *req.MaxWeight
	}
	if req.OperatingHours != nil {
		courier.OperatingHours = *req.OperatingHours
	}
	if req.BaseRatePerKm != nil {
		courier.BaseRatePerKm = *req.BaseRatePerKm
		courier.CustomPricing = true
	}
	if req.MinimumFare != nil {
		courier.MinimumFare = *req.MinimumFare
		courier.CustomPricing = true
	}
	if req.BankDetails != nil {
		courier.BankDetails = req.BankDetails
	}

	if err := s.repo.Update(ctx, courier); err != nil {
		return nil, err
	}

	return courier, nil
}

// ListAvailable lists all available couriers (for stores)
func (s *CourierService) ListAvailable(ctx context.Context, area string) ([]models.CourierListItem, error) {
	if area != "" {
		return s.repo.ListByArea(ctx, area)
	}
	return s.repo.ListActive(ctx)
}

// GetByID retrieves a courier by ID
func (s *CourierService) GetByID(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
	return s.repo.GetByID(ctx, id)
}

// generateToken creates a JWT token for a courier
func (s *CourierService) generateToken(courierID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"courier_id": courierID.String(),
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the courier ID
func (s *CourierService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	courierIDStr, ok := claims["courier_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid courier ID in token")
	}

	return uuid.Parse(courierIDStr)
}
