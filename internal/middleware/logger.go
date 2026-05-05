package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", latency).
			Str("request_id", c.Locals("request_id").(string)).
			Err(c.Context().Err()). //there is modif here
			Msg("request completed")

		return err
	}
}
