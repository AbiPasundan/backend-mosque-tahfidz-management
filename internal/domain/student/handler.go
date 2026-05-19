package student

import (
	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentHandler struct {
	service StudentService
}

func NewStudentHandler(service StudentService) *StudentHandler {
	return &StudentHandler{service: service}
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

	if err := h.service.DeleteStudent(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
