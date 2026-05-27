package activity_log

import (
	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogHandler struct {
	service ActivityLogService
}

func NewActivityLogHandler(service ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

func (h *ActivityLogHandler) ListActivityLogs(c *fiber.Ctx) error {
	page, limit := utils.GetPaginationParams(c)

	logs, total, err := h.service.ListActivityLogs(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	meta := utils.CreatePaginationMeta(total, page, limit)
	return utils.PaginatedResponse(c, fiber.StatusOK, "activity logs listed", logs, meta)
}
