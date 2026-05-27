package progress

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProgressRepository interface {
	Create(progress *Progress) error
	GetByID(id uuid.UUID) (*Progress, error)
	Update(progress *Progress) error
	List(studentID, date string, page, limit int) ([]Progress, int, error)
	GetDashboardSummary() (*DashboardSummary, error)
	GetWeeklyActivity() ([]DailyActivityCount, error)
	GetRecentProgress(limit int) ([]RecentProgressItem, error)
}

type progressRepository struct {
	db *sqlx.DB
}

func NewProgressRepository(db *sqlx.DB) ProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) Create(progress *Progress) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO progress (id, student_id, mentor_id, surah, status, ayat_start, ayat_end, notes, progress_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(query, progress.ID, progress.StudentID, progress.MentorID, progress.Surah,
		progress.Status, progress.AyatStart, progress.AyatEnd, progress.Notes, progress.ProgressDate)
	if err != nil {
		return err
	}

	// Update student's last_progress
	updateQuery := `UPDATE students SET last_progress = $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, progress.ProgressDate, progress.StudentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *progressRepository) GetByID(id uuid.UUID) (*Progress, error) {
	var progress Progress
	query := `
		SELECT p.*, u.name AS mentor_name 
		FROM progress p
		LEFT JOIN users u ON p.mentor_id = u.id
		WHERE p.id = $1
	`
	err := r.db.Get(&progress, query, id)
	return &progress, err
}

func (r *progressRepository) Update(progress *Progress) error {
	query := `UPDATE progress SET surah = $2, status = $3, ayat_start = $4, ayat_end = $5, notes = $6 WHERE id = $1`
	_, err := r.db.Exec(query, progress.ID, progress.Surah, progress.Status,
		progress.AyatStart, progress.AyatEnd, progress.Notes)
	return err
}

func (r *progressRepository) List(studentID, date string, page, limit int) ([]Progress, int, error) {
	var progress []Progress
	var total int

	baseQuery := `
		FROM progress p
		LEFT JOIN users u ON p.mentor_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	idx := 1

	if studentID != "" {
		baseQuery += fmt.Sprintf(` AND p.student_id = $%d`, idx)
		args = append(args, studentID)
		idx++
	}
	if date != "" {
		baseQuery += fmt.Sprintf(` AND p.progress_date = $%d`, idx)
		args = append(args, date)
		idx++
	}

	// Get total count
	countQuery := `SELECT COUNT(p.id) ` + baseQuery
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query := `SELECT p.*, u.name AS mentor_name ` + baseQuery + fmt.Sprintf(` ORDER BY p.progress_date DESC LIMIT $%d OFFSET $%d`, idx, idx+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&progress, query, args...)
	return progress, total, err
}

func (r *progressRepository) GetDashboardSummary() (*DashboardSummary, error) {
	var summary DashboardSummary
	query := `
		SELECT 
			(SELECT COUNT(*) FROM students WHERE deleted_at IS NULL) AS total_students,
			(SELECT COUNT(DISTINCT student_id) FROM progress WHERE progress_date = CURRENT_DATE) AS active_today,
			COALESCE(
				(SELECT COUNT(DISTINCT student_id) FROM progress WHERE progress_date >= CURRENT_DATE - INTERVAL '7 days') * 100.0 / 
				NULLIF((SELECT COUNT(*) FROM students WHERE deleted_at IS NULL), 0), 0
			) AS weekly_progress_percentage
	`
	err := r.db.Get(&summary, query)
	return &summary, err
}

func (r *progressRepository) GetWeeklyActivity() ([]DailyActivityCount, error) {
	var results []DailyActivityCount
	query := `
		SELECT 
			TO_CHAR(d.date, 'Dy') AS day,
			TO_CHAR(d.date, 'YYYY-MM-DD') AS date,
			COALESCE(COUNT(p.id), 0) AS count
		FROM generate_series(
			CURRENT_DATE - INTERVAL '6 days',
			CURRENT_DATE,
			'1 day'
		) AS d(date)
		LEFT JOIN progress p ON p.progress_date::date = d.date::date
		GROUP BY d.date
		ORDER BY d.date ASC
	`
	err := r.db.Select(&results, query)
	return results, err
}

func (r *progressRepository) GetRecentProgress(limit int) ([]RecentProgressItem, error) {
	var results []RecentProgressItem
	query := `
		SELECT 
			s.id::text AS student_id,
			s.name AS student_name,
			COALESCE(s.profile_img, '') AS profile_img,
			p.surah,
			p.ayat_start,
			p.ayat_end,
			p.status,
			COALESCE(u.name, 'Unknown') AS mentor_name,
			TO_CHAR(p.progress_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS progress_date
		FROM progress p
		JOIN students s ON s.id = p.student_id
		LEFT JOIN users u ON u.id = p.mentor_id
		ORDER BY p.progress_date DESC
		LIMIT $1
	`
	err := r.db.Select(&results, query, limit)
	return results, err
}
