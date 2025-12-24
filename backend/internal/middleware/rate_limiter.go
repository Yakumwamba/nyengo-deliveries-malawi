package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func RateLimiter(limit int, window time.Duration) fiber.Handler {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Clean old requests
		var valid []time.Time
		for _, t := range rl.requests[ip] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		rl.requests[ip] = valid

		if len(rl.requests[ip]) >= rl.limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   fiber.Map{"code": "RATE_LIMITED", "message": "Too many requests"},
			})
		}

		rl.requests[ip] = append(rl.requests[ip], now)
		return c.Next()
	}
}
