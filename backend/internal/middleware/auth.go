package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"nyengo-deliveries/internal/config"
)

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "UNAUTHORIZED", "message": "Missing authorization header"},
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "UNAUTHORIZED", "message": "Invalid token"},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "UNAUTHORIZED", "message": "Invalid token claims"},
			})
		}

		courierIDStr, _ := claims["courier_id"].(string)
		courierID, err := uuid.Parse(courierIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "UNAUTHORIZED", "message": "Invalid courier ID"},
			})
		}

		c.Locals("courier_id", courierID)
		return c.Next()
	}
}

// APIKeyAuth validates API keys for third-party store integrations
// Valid API keys are configured via STORE_API_KEYS environment variable
func APIKeyAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "UNAUTHORIZED", "message": "Missing API key"},
			})
		}

		// Validate the API key against configured keys
		if !cfg.ValidateStoreAPIKey(apiKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false, "error": fiber.Map{"code": "INVALID_API_KEY", "message": "Invalid API key"},
			})
		}

		return c.Next()
	}
}
