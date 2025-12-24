package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"nyengo-deliveries/internal/config"
	"nyengo-deliveries/internal/database"
	"nyengo-deliveries/internal/handlers"
	"nyengo-deliveries/internal/middleware"
	"nyengo-deliveries/internal/repository"
	"nyengo-deliveries/internal/services"
	"nyengo-deliveries/internal/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load application configuration
	cfg := config.LoadConfig()

	// Initialize database connections
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis for real-time features
	redisClient, err := database.NewRedisConnection(cfg.RedisURL)
	if err != nil {
		log.Printf("Redis connection failed (real-time features disabled): %v", err)
	}

	// Initialize repositories
	courierRepo := repository.NewCourierRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	deliveryRepo := repository.NewDeliveryRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	courierService := services.NewCourierService(courierRepo)
	pricingService := services.NewPricingService(cfg)
	orderService := services.NewOrderService(orderRepo, courierRepo, pricingService)
	notificationService := services.NewNotificationService(redisClient)
	trackingService := services.NewTrackingService(redisClient, deliveryRepo, orderRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo, courierRepo, cfg)
	externalCourierService := services.NewExternalCourierService(cfg)

	// Initialize WebSocket hub with Redis for cross-instance communication
	wsHub := websocket.NewHub()
	wsHub.SetRedis(redisClient)
	go wsHub.Run()

	// Initialize handlers
	courierHandler := handlers.NewCourierHandler(courierService)
	orderHandler := handlers.NewOrderHandler(orderService, notificationService, wsHub)
	pricingHandler := handlers.NewPricingHandler(pricingService)
	storeHandler := handlers.NewStoreHandler(courierService, orderService, pricingService, externalCourierService, cfg)
	webhookHandler := handlers.NewWebhookHandler(orderService, notificationService)
	trackingHandler := handlers.NewTrackingHandler(trackingService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Nyengo Deliveries API",
		ServerHeader: "Nyengo",
		ErrorHandler: handlers.CustomErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-API-Key",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"version":  "1.0.0",
			"currency": cfg.Currency,
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Rate limiting middleware for API routes
	api.Use(middleware.RateLimiter(cfg.RateLimitRequests, cfg.RateLimitDuration))

	// Public routes
	api.Post("/couriers/register", courierHandler.Register)
	api.Post("/couriers/login", courierHandler.Login)

	// Pricing routes (public for stores)
	pricing := api.Group("/pricing")
	pricing.Post("/estimate", pricingHandler.GetEstimate)

	// Store integration routes (API key authenticated)
	stores := api.Group("/stores")
	stores.Use(middleware.APIKeyAuth(cfg))
	stores.Get("/couriers", storeHandler.ListCouriers)
	stores.Post("/orders", storeHandler.CreateOrder)
	stores.Get("/orders/:id/status", storeHandler.GetOrderStatus)

	// Protected courier routes
	couriers := api.Group("/couriers")
	couriers.Use(middleware.JWTAuth(cfg.JWTSecret))
	couriers.Get("/profile", courierHandler.GetProfile)
	couriers.Put("/profile", courierHandler.UpdateProfile)
	couriers.Get("/dashboard", courierHandler.GetDashboard)
	couriers.Get("/available", courierHandler.ListAvailable)
	couriers.Get("/:id/rates", courierHandler.GetRates)

	// Protected order routes
	orders := api.Group("/orders")
	orders.Use(middleware.JWTAuth(cfg.JWTSecret))
	orders.Post("/", orderHandler.Create)
	orders.Get("/", orderHandler.List)
	orders.Get("/:id", orderHandler.GetByID)
	orders.Put("/:id/status", orderHandler.UpdateStatus)
	orders.Put("/:id/accept", orderHandler.Accept)
	orders.Put("/:id/decline", orderHandler.Decline)

	// WebSocket endpoint for real-time updates
	app.Get("/ws", middleware.JWTAuth(cfg.JWTSecret), websocket.HandleWebSocket(wsHub))

	// Webhook routes
	webhooks := api.Group("/webhooks")
	webhooks.Post("/payment", webhookHandler.HandlePayment)
	webhooks.Post("/delivery-update", webhookHandler.HandleDeliveryUpdate)

	// Tracking routes - live location tracking
	tracking := api.Group("/tracking")
	tracking.Get("/:orderId", trackingHandler.GetLiveTracking)                                             // Get current tracking
	tracking.Get("/:orderId/history", trackingHandler.GetLocationHistory)                                  // Get location history
	tracking.Post("/:orderId/start", middleware.JWTAuth(cfg.JWTSecret), trackingHandler.StartTracking)     // Start tracking
	tracking.Post("/:orderId/location", middleware.JWTAuth(cfg.JWTSecret), trackingHandler.UpdateLocation) // Update location
	tracking.Post("/:orderId/stop", middleware.JWTAuth(cfg.JWTSecret), trackingHandler.StopTracking)       // Stop tracking

	// Payment and Payout routes (courier authenticated)
	payments := api.Group("/payments")
	payments.Use(middleware.JWTAuth(cfg.JWTSecret))
	payments.Get("/verify/:orderId", paymentHandler.VerifyPayment)             // Verify order payment
	payments.Get("/payable-orders", paymentHandler.GetPayableOrders)           // Get orders eligible for payout
	payments.Post("/payouts", paymentHandler.RequestPayout)                    // Request payout
	payments.Get("/payouts", paymentHandler.GetPayoutHistory)                  // Get payout history
	payments.Get("/payouts/:payoutId", paymentHandler.GetPayoutByID)           // Get specific payout
	payments.Get("/earnings", paymentHandler.GetEarningsSummary)               // Get earnings summary
	payments.Get("/wallet/transactions", paymentHandler.GetWalletTransactions) // Get wallet transactions

	// Store payment verification routes (API key authenticated)
	stores.Get("/payments/verify/:orderId", paymentHandler.StoreVerifyPayment)

	// Admin payment management routes (should have admin auth in production)
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(cfg.JWTSecret)) // TODO: Add admin role check
	admin.Post("/payments/payouts/:payoutId/process", paymentHandler.ProcessPayout)

	log.Printf("📍 Live tracking enabled")
	log.Printf("💳 Payment & Payout system enabled")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 Nyengo Deliveries API starting on port %s", port)
	log.Printf("💰 Currency: %s", cfg.Currency)
	log.Printf("📍 Base rate: %s %.2f/km", cfg.CurrencySymbol, cfg.BaseRatePerKm)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
