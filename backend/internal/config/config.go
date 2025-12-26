package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server settings
	Port           string
	AllowedOrigins string

	// Database settings
	DatabaseURL string
	RedisURL    string

	// JWT settings
	JWTSecret     string
	JWTExpiration time.Duration

	// Currency settings (easily adjustable)
	Currency       string // "ZMW", "USD", "ZAR", etc.
	CurrencySymbol string // "K", "$", "R", etc.
	CurrencyLocale string // "en-ZM", "en-US", "en-ZA", etc.

	// Pricing settings (in configured currency)
	BaseRatePerKm    float64 // Base rate per kilometer
	MinimumFare      float64 // Minimum delivery fare
	SurgePricingMult float64 // Surge pricing multiplier (1.0 = no surge)
	PlatformFeePerc  float64 // Platform fee percentage (e.g., 0.15 for 15%)

	// Distance calculation settings
	MaxDeliveryDistance    float64 // Maximum delivery distance in km
	FreeDeliveryRadius     float64 // Free delivery radius in km (if applicable)
	LocalDistanceThreshold float64 // Distance threshold for local vs inter-city (default: 30km)

	// External courier settings
	ExternalCouriers []ExternalCourierConfig

	// Rate limiting
	RateLimitRequests int
	RateLimitDuration time.Duration

	// Business settings
	BusinessName    string
	BusinessCountry string
	SupportEmail    string
	SupportPhone    string

	// Store API Keys (for third-party store integrations)
	StoreAPIKeys []string

	// Webhook settings
	WebhookSecret string // Shared secret for delivery webhook authentication
}

// CurrencyPresets contains preset configurations for different currencies
var CurrencyPresets = map[string]struct {
	Symbol string
	Locale string
}{
	"ZMW": {Symbol: "K", Locale: "en-ZM"},   // Zambian Kwacha
	"USD": {Symbol: "$", Locale: "en-US"},   // US Dollar
	"ZAR": {Symbol: "R", Locale: "en-ZA"},   // South African Rand
	"KES": {Symbol: "KSh", Locale: "en-KE"}, // Kenyan Shilling
	"NGN": {Symbol: "₦", Locale: "en-NG"},   // Nigerian Naira
	"GHS": {Symbol: "GH₵", Locale: "en-GH"}, // Ghanaian Cedi
	"TZS": {Symbol: "TSh", Locale: "sw-TZ"}, // Tanzanian Shilling
	"UGX": {Symbol: "USh", Locale: "en-UG"}, // Ugandan Shilling
	"MWK": {Symbol: "MK", Locale: "en-MW"},  // Malawian Kwacha
	"BWP": {Symbol: "P", Locale: "en-BW"},   // Botswana Pula
	"EUR": {Symbol: "€", Locale: "en-EU"},   // Euro
	"GBP": {Symbol: "£", Locale: "en-GB"},   // British Pound
}

// ExternalCourierConfig represents configuration for an external courier service
type ExternalCourierConfig struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	LogoURL               string  `json:"logoUrl"`
	Description           string  `json:"description"`
	BaseRatePerKm         float64 `json:"baseRatePerKm"`
	MinimumFare           float64 `json:"minimumFare"`
	EstimatedDeliveryDays string  `json:"estimatedDeliveryDays"`
	ServiceType           string  `json:"serviceType"` // "express", "standard", "economy"
	IsActive              bool    `json:"isActive"`
}

// DefaultExternalCouriers returns default external courier configurations
func DefaultExternalCouriers() []ExternalCourierConfig {
	return []ExternalCourierConfig{
		{
			ID:                    "dhl",
			Name:                  "DHL Express",
			LogoURL:               "https://www.dhl.com/content/dam/dhl/global/core/images/logos/dhl-logo.svg",
			Description:           "Fast international and domestic express delivery",
			BaseRatePerKm:         3.50,
			MinimumFare:           150.0,
			EstimatedDeliveryDays: "1-2 business days",
			ServiceType:           "express",
			IsActive:              true,
		},
		{
			ID:                    "fedex",
			Name:                  "FedEx",
			LogoURL:               "https://www.fedex.com/content/dam/fedex-com/logos/logo.png",
			Description:           "Reliable express shipping nationwide",
			BaseRatePerKm:         3.00,
			MinimumFare:           120.0,
			EstimatedDeliveryDays: "2-3 business days",
			ServiceType:           "express",
			IsActive:              true,
		},
		{
			ID:                    "speedmail",
			Name:                  "Speed Mail Zambia",
			LogoURL:               "",
			Description:           "Local inter-city courier service",
			BaseRatePerKm:         2.00,
			MinimumFare:           80.0,
			EstimatedDeliveryDays: "2-4 business days",
			ServiceType:           "standard",
			IsActive:              true,
		},
		{
			ID:                    "zampost",
			Name:                  "Zambia Postal Services",
			LogoURL:               "",
			Description:           "Economy postal delivery",
			BaseRatePerKm:         1.50,
			MinimumFare:           50.0,
			EstimatedDeliveryDays: "5-7 business days",
			ServiceType:           "economy",
			IsActive:              true,
		},
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		// Server defaults
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),

		// Database defaults
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/nyengo?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		// JWT defaults
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),

		// Currency defaults (Zambian Kwacha as primary)
		Currency:       getEnv("CURRENCY", "ZMW"),
		CurrencySymbol: getEnv("CURRENCY_SYMBOL", ""), // Will be auto-filled
		CurrencyLocale: getEnv("CURRENCY_LOCALE", ""), // Will be auto-filled

		// Pricing defaults (in configured currency)
		BaseRatePerKm:    getFloatEnv("BASE_RATE_PER_KM", 5.0),      // K5 per km default
		MinimumFare:      getFloatEnv("MINIMUM_FARE", 20.0),         // K20 minimum
		SurgePricingMult: getFloatEnv("SURGE_MULTIPLIER", 1.0),      // No surge by default
		PlatformFeePerc:  getFloatEnv("PLATFORM_FEE_PERCENT", 0.10), // 10% platform fee

		// Distance defaults
		MaxDeliveryDistance:    getFloatEnv("MAX_DELIVERY_DISTANCE", 50.0),    // 50km max
		FreeDeliveryRadius:     getFloatEnv("FREE_DELIVERY_RADIUS", 0.0),      // No free delivery
		LocalDistanceThreshold: getFloatEnv("LOCAL_DISTANCE_THRESHOLD", 30.0), // 30km threshold for local vs inter-city

		// External couriers defaults
		ExternalCouriers: DefaultExternalCouriers(),

		// Rate limiting defaults
		RateLimitRequests: getIntEnv("RATE_LIMIT_REQUESTS", 100),
		RateLimitDuration: getDurationEnv("RATE_LIMIT_DURATION", time.Minute),

		// Business defaults
		BusinessName:    getEnv("BUSINESS_NAME", "Nyengo Deliveries"),
		BusinessCountry: getEnv("BUSINESS_COUNTRY", "Zambia"),
		SupportEmail:    getEnv("SUPPORT_EMAIL", "support@nyengo.com"),
		SupportPhone:    getEnv("SUPPORT_PHONE", "+260970000000"),

		// Store API Keys (comma-separated in env)
		StoreAPIKeys: getStoreAPIKeys(),

		// Webhook settings
		WebhookSecret: getEnv("WEBHOOK_SECRET", "nyg_webhook_secret_dev_2024"),
	}

	// Auto-fill currency symbol and locale if not set
	if cfg.CurrencySymbol == "" || cfg.CurrencyLocale == "" {
		if preset, exists := CurrencyPresets[cfg.Currency]; exists {
			if cfg.CurrencySymbol == "" {
				cfg.CurrencySymbol = preset.Symbol
			}
			if cfg.CurrencyLocale == "" {
				cfg.CurrencyLocale = preset.Locale
			}
		} else {
			// Fallback defaults
			if cfg.CurrencySymbol == "" {
				cfg.CurrencySymbol = cfg.Currency + " "
			}
			if cfg.CurrencyLocale == "" {
				cfg.CurrencyLocale = "en-US"
			}
		}
	}

	return cfg
}

// GetPricingTiers returns dynamic pricing tiers based on configuration
func (c *Config) GetPricingTiers() []PricingTier {
	return []PricingTier{
		{
			Name:        "Standard",
			Description: "Regular delivery",
			Multiplier:  1.0,
			MaxWeight:   10.0, // kg
		},
		{
			Name:        "Express",
			Description: "Priority delivery",
			Multiplier:  1.5,
			MaxWeight:   10.0,
		},
		{
			Name:        "Heavy",
			Description: "Heavy items",
			Multiplier:  2.0,
			MaxWeight:   50.0,
		},
		{
			Name:        "Fragile",
			Description: "Fragile handling",
			Multiplier:  1.3,
			MaxWeight:   10.0,
		},
	}
}

// PricingTier represents a delivery pricing category
type PricingTier struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Multiplier  float64 `json:"multiplier"`
	MaxWeight   float64 `json:"maxWeight"`
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// FormatCurrency formats an amount according to the configured currency
func (c *Config) FormatCurrency(amount float64) string {
	// Format with 2 decimal places and thousands separator
	formatted := formatWithCommas(amount)
	return c.CurrencySymbol + formatted
}

func formatWithCommas(amount float64) string {
	// Simple number formatting
	str := strconv.FormatFloat(amount, 'f', 2, 64)
	parts := strings.Split(str, ".")

	intPart := parts[0]
	decPart := parts[1]

	// Add commas for thousands
	var result strings.Builder
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}

	return result.String() + "." + decPart
}

// getStoreAPIKeys returns the list of valid store API keys
// Keys are loaded from STORE_API_KEYS env var (comma-separated)
// A default test key is provided for development
func getStoreAPIKeys() []string {
	envKeys := os.Getenv("STORE_API_KEYS")
	if envKeys != "" {
		keys := strings.Split(envKeys, ",")
		// Trim whitespace from each key
		for i := range keys {
			keys[i] = strings.TrimSpace(keys[i])
		}
		return keys
	}

	// Default test API key for development
	// IMPORTANT: In production, always set STORE_API_KEYS env variable
	return []string{
		"nyg_test_store_api_key_2024_dev",
	}
}

// ValidateStoreAPIKey checks if the provided API key is valid
func (c *Config) ValidateStoreAPIKey(apiKey string) bool {
	for _, validKey := range c.StoreAPIKeys {
		if validKey == apiKey {
			return true
		}
	}
	return false
}
