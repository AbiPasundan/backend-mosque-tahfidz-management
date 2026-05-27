package progress

import (
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/internal/domain/student"
	"backend-mosque-tahfidz-management/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProgressHandler struct {
	service            ProgressService
	authService        auth.AuthService
	studentService     student.StudentService
	activityLogService activity_log.ActivityLogService
}

func NewProgressHandler(
	service ProgressService,
	authService auth.AuthService,
	studentService student.StudentService,
	activityLogService activity_log.ActivityLogService,
) *ProgressHandler {
	return &ProgressHandler{
		service:            service,
		authService:        authService,
		studentService:     studentService,
		activityLogService: activityLogService,
	}
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

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				studentName := "Student"
				if st, err := h.studentService.GetStudent(resp.StudentID); err == nil {
					studentName = st.Name
				}
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"CREATE_PROGRESS",
					"progress",
					&resp.ID,
					fmt.Sprintf("Logged progress for %s: Surah %s (Ayat %d-%d, Status: %s)", studentName, resp.Surah, resp.AyatStart, resp.AyatEnd, resp.Status),
				)
			}
		}
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

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"CREATE_PROGRESS_BULK",
					"progress",
					nil,
					fmt.Sprintf("Logged bulk progress entries for %d records", len(resp)),
				)
			}
		}
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
	page, limit := utils.GetPaginationParams(c)

	progressList, total, err := h.service.ListProgress(studentID, date, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	meta := utils.CreatePaginationMeta(total, page, limit)
	return utils.PaginatedResponse(c, fiber.StatusOK, "progress listed", progressList, meta)
}

func (h *ProgressHandler) GetDashboardSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetDashboardSummary()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "dashboard summary", summary)
}
