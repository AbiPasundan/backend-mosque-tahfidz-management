package student

import (
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentHandler struct {
	service            StudentService
	authService        auth.AuthService
	activityLogService activity_log.ActivityLogService
}

func NewStudentHandler(service StudentService, authService auth.AuthService, activityLogService activity_log.ActivityLogService) *StudentHandler {
	return &StudentHandler{
		service:            service,
		authService:        authService,
		activityLogService: activityLogService,
	}
}

func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
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

	resp, err := h.service.CreateStudent(&req, mentorID)
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
					"CREATE_STUDENT",
					"student",
					&resp.ID,
					fmt.Sprintf("Registered new student: %s", resp.Name),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "student created", resp)
}

func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid student id")
	}

	resp, err := h.service.GetStudent(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "student found", resp)
}

func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid student id")
	}

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.UpdateStudent(id, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "student updated", resp)
}

func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid student id")
	}

	// Fetch student first to get their name for the activity log
	studentName := "Unknown Student"
	if st, err := h.service.GetStudent(id); err == nil {
		studentName = st.Name
	}

	if err := h.service.DeleteStudent(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Log Activity
	if actorIDStr, ok := c.Locals("user_id").(string); ok {
		if actorID, err := uuid.Parse(actorIDStr); err == nil {
			if actor, err := h.authService.GetUser(actorID); err == nil {
				_ = h.activityLogService.LogAction(
					&actorID,
					actor.Name,
					"DELETE_STUDENT",
					"student",
					&id,
					fmt.Sprintf("Deleted student: %s", studentName),
				)
			}
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "student deleted", nil)
}

func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	search := c.Query("search")
	status := c.Query("status")
	learningLevel := c.Query("learning_level")
	page, limit := utils.GetPaginationParams(c)

	students, total, err := h.service.ListStudents(search, status, learningLevel, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	meta := utils.CreatePaginationMeta(total, page, limit)
	return utils.PaginatedResponse(c, fiber.StatusOK, "students listed", students, meta)
}
