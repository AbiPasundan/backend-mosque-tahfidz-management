package middleware

import (
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/gofiber/fiber/v2"
)

func JWT(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Cek token dari cookie (Prioritas utama untuk keamanan HttpOnly)
		tokenString := c.Cookies("auth_token")

		// 2. Jika cookie kosong, cek dari Authorization header (Fallback untuk mobile/API)
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

		// 3. Proteksi CSRF sederhana: Jika menggunakan cookie, wajib ada custom header dari client (Axios)
		if c.Cookies("auth_token") != "" && c.Get("X-Requested-With") == "" {
			if c.Method() != "GET" && c.Method() != "OPTIONS" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "potential CSRF detected: missing X-Requested-With header",
				})
			}
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
