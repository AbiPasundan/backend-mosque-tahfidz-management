package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalData   int `json:"total_data"`
	Limit       int `json:"limit"`
}

func GetPaginationParams(c *fiber.Ctx) (page, limit int) {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// Cap limit to 100 to prevent abuse
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

func CreatePaginationMeta(totalData, page, limit int) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginationMeta{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalData:   totalData,
		Limit:       limit,
	}
}
