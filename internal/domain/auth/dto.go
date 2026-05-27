package auth

import (
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
	Role   string    `json:"role"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin mentor"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin mentor"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}

type MentorDetailResponse struct {
	ID       uuid.UUID               `json:"id"`
	Name     string                  `json:"name"`
	Email    string                  `json:"email"`
	Role     string                  `json:"role"`
	Students []MentorStudentResponse `json:"students"`
}

type MentorStudentResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ProfileImg    string    `json:"profile_img"`
	CoverImg      string    `json:"cover_img"`
	Age           int       `json:"age"`
	LearningLevel string    `json:"learning_level"`
	Fluency       string    `json:"fluency"`
	Status        string    `json:"status"`
	Contact       string    `json:"contact"`
	JoinDate      string    `json:"join_date"`
	LastProgress  string    `json:"last_progress,omitempty"`
}
