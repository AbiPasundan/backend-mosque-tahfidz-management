package progress

import (
	"time"

	"github.com/google/uuid"
)

type Progress struct {
	ID           uuid.UUID `db:"id"`
	StudentID    uuid.UUID `db:"student_id"`
	MentorID     uuid.UUID `db:"mentor_id"`
	MentorName   string    `db:"mentor_name"`
	Surah        string    `db:"surah"`
	Status       string    `db:"status"`
	AyatStart    int       `db:"ayat_start"`
	AyatEnd      int       `db:"ayat_end"`
	Notes        string    `db:"notes"`
	ProgressDate time.Time `db:"progress_date"`
}

type DashboardSummary struct {
	TotalStudents            int     `db:"total_students"`
	ActiveToday              int     `db:"active_today"`
	WeeklyProgressPercentage float64 `db:"weekly_progress_percentage"`
}
