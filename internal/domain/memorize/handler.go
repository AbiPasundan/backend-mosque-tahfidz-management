package memorize

import (
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/pkg/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MemorizeHandler struct {
	service            MemorizeService
	authService        auth.AuthService
	activityLogService activity_log.ActivityLogService
}

func NewMemorizeHandler(
	service MemorizeService,
	authService auth.AuthService,
	activityLogService activity_log.ActivityLogService,
) *MemorizeHandler {
	return &MemorizeHandler{
		service:            service,
		authService:        authService,
		activityLogService: activityLogService,
	}
}

func (h *MemorizeHandler) CreateMemorize(c *fiber.Ctx) error {
	var req CreateMemorizeRequest
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

	resp, err := h.service.CreateMemorize(&req, mentorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"CREATE_MEMORIZE",
					"memorize",
					&resp.ID,
					fmt.Sprintf("Logged memorize entry: Surah %s (Ayat %d-%d, Status: %s)", resp.Surah, resp.AyatStart, resp.AyatEnd, resp.Status),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "memorize record created", resp)
}

func (h *MemorizeHandler) ListByStudent(c *fiber.Ctx) error {
	studentIDStr := c.Query("student_id")
	surah := c.Query("surah")
	status := c.Query("status")
	page, limit := utils.GetPaginationParams(c)

	var studentID uuid.UUID
	if studentIDStr != "" {
		var err error
		studentID, err = uuid.Parse(studentIDStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid student_id")
		}
	}

	list, total, err := h.service.ListByStudent(studentID, surah, status, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	meta := utils.CreatePaginationMeta(total, page, limit)
	return utils.PaginatedResponse(c, fiber.StatusOK, "memorize records listed", list, meta)
}

func (h *MemorizeHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid memorize id")
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "memorize record found", resp)
}

func (h *MemorizeHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid memorize id")
	}

	var req UpdateMemorizeStatusRequest
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

	resp, err := h.service.UpdateStatus(id, &req, mentorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"UPDATE_MEMORIZE_STATUS",
					"memorize",
					&resp.ID,
					fmt.Sprintf("Updated memorize status for Surah %s (Ayat %d-%d) to %s", resp.Surah, resp.AyatStart, resp.AyatEnd, resp.Status),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "memorize status updated", resp)
}

func (h *MemorizeHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	var req BulkUpdateStatusRequest
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

	err = h.service.BulkUpdateStatus(&req, mentorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"BULK_UPDATE_MEMORIZE_STATUS",
					"memorize",
					nil,
					fmt.Sprintf("Bulk updated status to %s for %d memorize records", req.Status, len(req.IDs)),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "memorize records status updated", nil)
}

func (h *MemorizeHandler) DeleteMemorize(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid memorize id")
	}

	err = h.service.DeleteMemorize(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"DELETE_MEMORIZE",
					"memorize",
					&id,
					fmt.Sprintf("Deleted memorize record: %s", id),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "memorize record deleted", nil)
}

func (h *MemorizeHandler) GetStudentSurahDetail(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid student id")
	}

	surahNumber, err := strconv.Atoi(c.Params("number"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid surah number")
	}

	list, err := h.service.GetStudentSurahDetail(studentID, surahNumber)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "student surah details listed", list)
}
