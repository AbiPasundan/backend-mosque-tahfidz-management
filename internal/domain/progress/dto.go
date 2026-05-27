package progress

import (
	"github.com/google/uuid"
)

type CreateProgressRequest struct {
	StudentID uuid.UUID `json:"student_id" validate:"required"`
	Surah     string    `json:"surah" validate:"required"`
	AyatStart int       `json:"ayat_start" validate:"required,min=1"`
	AyatEnd   int       `json:"ayat_end" validate:"required,min=1,gtefield=AyatStart"`
	Status    string    `json:"status" validate:"required"`
	Notes     string    `json:"notes"`
}

type UpdateProgressRequest struct {
	Surah     string `json:"surah" validate:"required"`
	AyatStart int    `json:"ayat_start" validate:"required,min=1"`
	AyatEnd   int    `json:"ayat_end" validate:"required,min=1,gtefield=AyatStart"`
	Status    string `json:"status" validate:"required"`
	Notes     string `json:"notes"`
}

type BulkCreateProgressRequest struct {
	Items []CreateProgressRequest `json:"items" validate:"required,dive"`
}

type ProgressResponse struct {
	ID           uuid.UUID `json:"id"`
	StudentID    uuid.UUID `json:"student_id"`
	MentorID     uuid.UUID `json:"mentor_id"`
	MentorName   string    `json:"mentor_name,omitempty"`
	Surah        string    `json:"surah"`
	Status       string    `json:"status"`
	AyatStart    int       `json:"ayat_start"`
	AyatEnd      int       `json:"ayat_end"`
	Notes        string    `json:"notes,omitempty"`
	ProgressDate string    `json:"progress_date"`
}

type DashboardSummaryResponse struct {
	TotalStudents            int                    `json:"total_students"`
	ActiveToday              int                    `json:"active_today"`
	WeeklyProgressPercentage float64                `json:"weekly_progress_percentage"`
	WeeklyActivity           []DailyActivityCount   `json:"weekly_activity"`
	RecentProgress           []RecentProgressItem   `json:"recent_progress"`
}
