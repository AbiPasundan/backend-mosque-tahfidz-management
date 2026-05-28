package middleware

import (
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/gofiber/fiber/v2"
)

func JWT(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Cek token dari cookie
		tokenString := c.Cookies("auth_token")

		// 2. Jika cookie kosong, cek dari Authorization header (standard Bearer token)
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authentication token",
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
