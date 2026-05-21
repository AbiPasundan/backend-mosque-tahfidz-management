package auth

import (
	"backend-mosque-tahfidz-management/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary Login
// @Description Login user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} utils.BaseResponse{data=LoginResponse} "login successful"
// @Failure 400 {object} utils.BaseResponse "invalid request body"
// @Failure 401 {object} utils.BaseResponse "invalid credentials"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    resp.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "login successful", fiber.Map{
		"role": resp.Role,
	})
}

// logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "logout successful", nil)
}

// me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid user session")
	}

	resp, err := h.service.GetUser(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user profile fetched", resp)
}

// @Summary Register new user (Admin only)
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body CreateUserRequest true "User data"
// @Success 201 {object} utils.BaseResponse{data=UserResponse} "user created"
// @Failure 400 {object} utils.BaseResponse "invalid request body"
// @Failure 401 {object} utils.BaseResponse "unauthorized"
// @Failure 500 {object} utils.BaseResponse "internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.CreateUser(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "user created", resp)
}

func (h *AuthHandler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid user id")
	}

	resp, err := h.service.GetUser(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user found", resp)
}

func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid user id")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user updated", resp)
}

func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteUser(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user deleted", nil)
}

func (h *AuthHandler) ListUsers(c *fiber.Ctx) error {
	search := c.Query("search")
	role := c.Query("role")
	page, limit := utils.GetPaginationParams(c)

	users, total, err := h.service.ListUsers(search, role, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	meta := utils.CreatePaginationMeta(total, page, limit)
	return utils.PaginatedResponse(c, fiber.StatusOK, "users listed", users, meta)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid user context")
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	resp, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "profile updated", resp)
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid user context")
	}

	var req UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.service.UpdatePassword(userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "password updated", nil)
}
