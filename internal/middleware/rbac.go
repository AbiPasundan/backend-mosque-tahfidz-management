package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func RBAC(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if slices.Contains(allowedRoles, role) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "access forbidden",
		})
	}
}
