package progress

import (
	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProgressHandler struct {
	service ProgressService
}

func NewProgressHandler(service ProgressService) *ProgressHandler {
	return &ProgressHandler{service: service}
}

func (h *ProgressHandler) CreateProgress(c *fiber.Ctx) error {
	var req CreateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	mentorID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid user context")
	}

	resp, err := h.service.CreateProgress(&req, mentorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "progress created", resp)
}

func (h *ProgressHandler) BulkCreateProgress(c *fiber.Ctx) error {
	var req BulkCreateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	mentorID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid user context")
	}

	resp, err := h.service.BulkCreateProgress(&req, mentorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "progress records created", resp)
}

func (h *ProgressHandler) UpdateProgress(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid progress id")
	}

	var req UpdateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.UpdateProgress(id, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "progress updated", resp)
}

func (h *ProgressHandler) ListProgress(c *fiber.Ctx) error {
	studentID := c.Query("student_id")
	date := c.Query("date")

	progressList, err := h.service.ListProgress(studentID, date)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "progress listed", progressList)
}

func (h *ProgressHandler) GetDashboardSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetDashboardSummary()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "dashboard summary", summary)
}
