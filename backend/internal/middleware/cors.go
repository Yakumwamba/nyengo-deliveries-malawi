package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func CORS(allowedOrigins string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-API-Key")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}
