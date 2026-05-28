package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewApp creates a pre-configured Fiber instance with body limit and error handler.
func NewApp() *fiber.App {
	return fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB for file uploads
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})
}

// CORS applies the CORS middleware with configurable allowed origins.
func CORS(allowOrigins string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: true,
	})
}
