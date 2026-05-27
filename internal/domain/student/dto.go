package student

import (
	"github.com/google/uuid"
)

type CreateStudentRequest struct {
	Name          string `json:"name" validate:"required"`
	Username      string `json:"username" validate:"required,min=3,max=50"`
	Password      string `json:"password" validate:"required,min=6"`
	Age           int    `json:"age" validate:"required,min=5,max=25"`
	LearningLevel string `json:"learning_level" validate:"required"`
	Contact       string `json:"contact" validate:"required"`
	Status        string `json:"status" validate:"required"`
}

type UpdateStudentRequest struct {
	Name          string `json:"name"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ProfileImg    string `json:"profile_img"`
	CoverImg      string `json:"cover_img"`
	Age           int    `json:"age" validate:"omitempty,min=5,max=25"`
	LearningLevel string `json:"learning_level"`
	Fluency       string `json:"fluency"`
	Contact       string `json:"contact"`
	Status        string `json:"status"`
}

type StudentResponse struct {
	ID            uuid.UUID `json:"id"`
	MentorID      uuid.UUID `json:"mentor_id"`
	MentorName    string    `json:"mentor_name,omitempty"`
	Name          string    `json:"name"`
	Username      string    `json:"username"`
	ProfileImg    string    `json:"profile_img,omitempty"`
	CoverImg      string    `json:"cover_img,omitempty"`
	Age           int       `json:"age"`
	LearningLevel string    `json:"learning_level"`
	Fluency       string    `json:"fluency,omitempty"`
	Status        string    `json:"status"`
	Contact       string    `json:"contact"`
	JoinDate      string    `json:"join_date"`
	LastProgress  string    `json:"last_progress,omitempty"`
}
