package middleware

import (
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/gofiber/fiber/v2"
)

func JWT(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("auth_token")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authentication cookie",
			})
		}

		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid or expired token",
			})
		}

		c.Locals("user_id", payload.UserID)
		c.Locals("role", payload.Role)

		return c.Next()
	}
}
