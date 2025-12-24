package models

import (
	"time"

	"github.com/google/uuid"
)

// Courier represents a courier company or individual courier
type Courier struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"-" db:"password"` // Never expose password
	CompanyName    string    `json:"companyName" db:"company_name"`
	OwnerName      string    `json:"ownerName" db:"owner_name"`
	Phone          string    `json:"phone" db:"phone"`
	AlternatePhone string    `json:"alternatePhone,omitempty" db:"alternate_phone"`
	WhatsApp       string    `json:"whatsapp,omitempty" db:"whatsapp"`
	Address        string    `json:"address" db:"address"`
	City           string    `json:"city" db:"city"`
	Country        string    `json:"country" db:"country"`
	LogoURL        string    `json:"logoUrl,omitempty" db:"logo_url"`
	Description    string    `json:"description,omitempty" db:"description"`

	// Service configuration
	ServiceAreas   []string       `json:"serviceAreas" db:"service_areas"`
	VehicleTypes   []string       `json:"vehicleTypes" db:"vehicle_types"`
	MaxWeight      float64        `json:"maxWeight" db:"max_weight"` // in kg
	OperatingHours OperatingHours `json:"operatingHours" db:"operating_hours"`

	// Pricing configuration (overrides system defaults)
	BaseRatePerKm float64 `json:"baseRatePerKm" db:"base_rate_per_km"`
	MinimumFare   float64 `json:"minimumFare" db:"minimum_fare"`
	CustomPricing bool    `json:"customPricing" db:"custom_pricing"`

	// Ratings and statistics
	Rating          float64 `json:"rating" db:"rating"`
	TotalReviews    int     `json:"totalReviews" db:"total_reviews"`
	TotalDeliveries int     `json:"totalDeliveries" db:"total_deliveries"`
	SuccessRate     float64 `json:"successRate" db:"success_rate"`

	// Account status
	IsVerified       bool     `json:"isVerified" db:"is_verified"`
	IsActive         bool     `json:"isActive" db:"is_active"`
	IsFeatured       bool     `json:"isFeatured" db:"is_featured"`
	VerificationDocs []string `json:"verificationDocs,omitempty" db:"verification_docs"`

	// Financial
	WalletBalance float64      `json:"walletBalance" db:"wallet_balance"`
	BankDetails   *BankDetails `json:"bankDetails,omitempty" db:"bank_details"`

	// Timestamps
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty" db:"last_active_at"`
}

// OperatingHours defines business hours
type OperatingHours struct {
	Monday    DayHours `json:"monday"`
	Tuesday   DayHours `json:"tuesday"`
	Wednesday DayHours `json:"wednesday"`
	Thursday  DayHours `json:"thursday"`
	Friday    DayHours `json:"friday"`
	Saturday  DayHours `json:"saturday"`
	Sunday    DayHours `json:"sunday"`
}

// DayHours represents hours for a single day
type DayHours struct {
	Open   string `json:"open"`   // "08:00"
	Close  string `json:"close"`  // "18:00"
	Closed bool   `json:"closed"` // true if closed this day
}

// BankDetails contains banking information for payouts
type BankDetails struct {
	BankName      string `json:"bankName"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	BranchCode    string `json:"branchCode,omitempty"`
	SwiftCode     string `json:"swiftCode,omitempty"`
}

// CourierListItem is a simplified courier for lists
type CourierListItem struct {
	ID              uuid.UUID `json:"id"`
	CompanyName     string    `json:"companyName"`
	LogoURL         string    `json:"logoUrl,omitempty"`
	Rating          float64   `json:"rating"`
	TotalReviews    int       `json:"totalReviews"`
	TotalDeliveries int       `json:"totalDeliveries"`
	BaseRatePerKm   float64   `json:"baseRatePerKm"`
	MinimumFare     float64   `json:"minimumFare"`
	IsVerified      bool      `json:"isVerified"`
	IsFeatured      bool      `json:"isFeatured"`
}

// CourierRegistrationRequest is the request body for registering a new courier
type CourierRegistrationRequest struct {
	Email        string   `json:"email" validate:"required,email"`
	Password     string   `json:"password" validate:"required,min=8"`
	CompanyName  string   `json:"companyName" validate:"required,min=2"`
	OwnerName    string   `json:"ownerName" validate:"required,min=2"`
	Phone        string   `json:"phone" validate:"required"`
	Address      string   `json:"address" validate:"required"`
	City         string   `json:"city" validate:"required"`
	Country      string   `json:"country" validate:"required"`
	ServiceAreas []string `json:"serviceAreas" validate:"required,min=1"`
	VehicleTypes []string `json:"vehicleTypes" validate:"required,min=1"`
}

// CourierLoginRequest is the request body for courier login
type CourierLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CourierLoginResponse is the response after successful login
type CourierLoginResponse struct {
	Token   string   `json:"token"`
	Courier *Courier `json:"courier"`
}

// CourierUpdateRequest is the request body for updating courier profile
type CourierUpdateRequest struct {
	CompanyName    *string         `json:"companyName,omitempty"`
	OwnerName      *string         `json:"ownerName,omitempty"`
	Phone          *string         `json:"phone,omitempty"`
	AlternatePhone *string         `json:"alternatePhone,omitempty"`
	WhatsApp       *string         `json:"whatsapp,omitempty"`
	Address        *string         `json:"address,omitempty"`
	City           *string         `json:"city,omitempty"`
	LogoURL        *string         `json:"logoUrl,omitempty"`
	Description    *string         `json:"description,omitempty"`
	ServiceAreas   []string        `json:"serviceAreas,omitempty"`
	VehicleTypes   []string        `json:"vehicleTypes,omitempty"`
	MaxWeight      *float64        `json:"maxWeight,omitempty"`
	OperatingHours *OperatingHours `json:"operatingHours,omitempty"`
	BaseRatePerKm  *float64        `json:"baseRatePerKm,omitempty"`
	MinimumFare    *float64        `json:"minimumFare,omitempty"`
	BankDetails    *BankDetails    `json:"bankDetails,omitempty"`
}
