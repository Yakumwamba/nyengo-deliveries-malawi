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
	webhookHandler := handlers.NewWebhookHandler(orderService, notificationService, orderRepo, cfg)
	trackingHandler := handlers.NewTrackingHandler(trackingService, orderRepo)
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
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-API-Key,X-Webhook-Secret",
		AllowCredentials: true,
	}))

	// Favicon handler
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		// Return a simple 1x1 transparent PNG as favicon
		c.Set("Content-Type", "image/x-icon")
		c.Set("Cache-Control", "public, max-age=31536000")
		// Minimal ICO file (1x1 transparent)
		ico := []byte{
			0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01,
			0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x30, 0x00,
			0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		return c.Send(ico)
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"version":  "1.0.0",
			"currency": cfg.Currency,
		})
	})

	// API root info endpoint
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     "Nyengo Deliveries API",
			"version":     "1.0.0",
			"description": "Delivery management and courier services API",
			"endpoints": fiber.Map{
				"health": "/health",
				"api_v1": "/api/v1",
				"docs":   "See API documentation for full endpoint list",
			},
		})
	})
	app.Get("/api/", func(c *fiber.Ctx) error {
		return c.Redirect("/api", fiber.StatusMovedPermanently)
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// API v1 info endpoint
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     "Nyengo Deliveries API",
			"version":     "v1",
			"description": "Delivery management and courier services API",
			"currency":    cfg.Currency,
			"endpoints": fiber.Map{
				"auth": fiber.Map{
					"register": "POST /api/v1/couriers/register",
					"login":    "POST /api/v1/couriers/login",
				},
				"couriers": fiber.Map{
					"profile":   "GET /api/v1/couriers/profile",
					"available": "GET /api/v1/couriers/available",
					"rates":     "GET /api/v1/couriers/:id/rates",
				},
				"stores": fiber.Map{
					"list_couriers": "GET /api/v1/stores/couriers",
					"create_order":  "POST /api/v1/stores/orders",
					"order_status":  "GET /api/v1/stores/orders/:id/status",
				},
				"orders": fiber.Map{
					"create":        "POST /api/v1/orders",
					"list":          "GET /api/v1/orders",
					"get":           "GET /api/v1/orders/:id",
					"update_status": "PUT /api/v1/orders/:id/status",
					"accept":        "PUT /api/v1/orders/:id/accept",
					"decline":       "PUT /api/v1/orders/:id/decline",
				},
				"tracking": fiber.Map{
					"live":    "GET /api/v1/tracking/:orderId",
					"history": "GET /api/v1/tracking/:orderId/history",
					"start":   "POST /api/v1/tracking/:orderId/start",
					"update":  "POST /api/v1/tracking/:orderId/location",
					"stop":    "POST /api/v1/tracking/:orderId/stop",
				},
				"payments": fiber.Map{
					"verify":       "GET /api/v1/payments/verify/:orderId",
					"payable":      "GET /api/v1/payments/payable-orders",
					"request":      "POST /api/v1/payments/payouts",
					"history":      "GET /api/v1/payments/payouts",
					"earnings":     "GET /api/v1/payments/earnings",
					"transactions": "GET /api/v1/payments/wallet/transactions",
				},
				"pricing": fiber.Map{
					"estimate": "POST /api/v1/pricing/estimate",
				},
			},
		})
	})

	// Rate limiting middleware for API routes
	api.Use(middleware.RateLimiter(cfg.RateLimitRequests, cfg.RateLimitDuration))

	// Delivery webhook from courier platform (not rate limited, auth via X-Webhook-Secret)
	// POST /api/delivery/webhook
	app.Post("/api/delivery/webhook", webhookHandler.HandleDeliveryWebhook)

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
