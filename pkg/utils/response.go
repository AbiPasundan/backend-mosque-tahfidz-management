package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func ValidationErrorResponse(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "validation error",
		"errors":  errors,
	})
}
