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

// DailyActivityCount represents progress count per day — used for the dashboard chart.
type DailyActivityCount struct {
	Day   string `db:"day" json:"day"`
	Date  string `db:"date" json:"date"`
	Count int    `db:"count" json:"count"`
}

// RecentProgressItem represents a single recent progress entry for the dashboard table.
type RecentProgressItem struct {
	StudentID    string `db:"student_id" json:"student_id"`
	StudentName  string `db:"student_name" json:"student_name"`
	ProfileImg   string `db:"profile_img" json:"profile_img"`
	Surah        string `db:"surah" json:"surah"`
	AyatStart    int    `db:"ayat_start" json:"ayat_start"`
	AyatEnd      int    `db:"ayat_end" json:"ayat_end"`
	Status       string `db:"status" json:"status"`
	MentorName   string `db:"mentor_name" json:"mentor_name"`
	ProgressDate string `db:"progress_date" json:"progress_date"`
}
