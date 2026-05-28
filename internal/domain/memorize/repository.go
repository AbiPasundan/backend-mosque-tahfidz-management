package memorize

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MemorizeRepository interface {
	Create(m *Memorize) error
	GetByID(id uuid.UUID) (*Memorize, error)
	Update(m *Memorize) error
	Delete(id uuid.UUID) error
	ListByStudent(studentID uuid.UUID, surah, status string, page, limit int) ([]Memorize, int, error)
	GetStudentSurahDetail(studentID uuid.UUID, surahNumber int) ([]Memorize, error)
	BulkUpdateStatus(ids []uuid.UUID, status string, verifiedBy uuid.UUID) error
}

type memorizeRepository struct {
	db *sqlx.DB
}

func NewMemorizeRepository(db *sqlx.DB) MemorizeRepository {
	return &memorizeRepository{db: db}
}

func (r *memorizeRepository) Create(m *Memorize) error {
	query := `
		INSERT INTO memorize (
			id, student_id, verified_by, surah, surah_number, 
			ayat_start, ayat_end, status, notes, 
			memorized_at, last_reviewed_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`
	_, err := r.db.Exec(
		query, m.ID, m.StudentID, m.VerifiedBy, m.Surah, m.SurahNumber,
		m.AyatStart, m.AyatEnd, m.Status, m.Notes,
		m.MemorizedAt, m.LastReviewedAt,
	)
	return err
}

func (r *memorizeRepository) GetByID(id uuid.UUID) (*Memorize, error) {
	var m Memorize
	query := `
		SELECT m.*, u.name AS verifier_name 
		FROM memorize m 
		LEFT JOIN users u ON m.verified_by = u.id 
		WHERE m.id = $1
	`
	err := r.db.Get(&m, query, id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *memorizeRepository) Update(m *Memorize) error {
	query := `
		UPDATE memorize 
		SET surah = $2, surah_number = $3, ayat_start = $4, ayat_end = $5, 
		    status = $6, notes = $7, verified_by = $8, 
		    memorized_at = $9, last_reviewed_at = $10, updated_at = NOW() 
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query, m.ID, m.Surah, m.SurahNumber, m.AyatStart, m.AyatEnd,
		m.Status, m.Notes, m.VerifiedBy, m.MemorizedAt, m.LastReviewedAt,
	)
	return err
}

func (r *memorizeRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM memorize WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *memorizeRepository) ListByStudent(studentID uuid.UUID, surah, status string, page, limit int) ([]Memorize, int, error) {
	var list []Memorize
	var total int

	baseQuery := `
		FROM memorize m
		LEFT JOIN users u ON m.verified_by = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	idx := 1

	if studentID != uuid.Nil {
		baseQuery += fmt.Sprintf(` AND m.student_id = $%d`, idx)
		args = append(args, studentID)
		idx++
	}
	if surah != "" {
		baseQuery += fmt.Sprintf(` AND m.surah ILIKE $%d`, idx)
		args = append(args, "%"+surah+"%")
		idx++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(` AND m.status = $%d`, idx)
		args = append(args, status)
		idx++
	}

	// Get total count
	countQuery := `SELECT COUNT(m.id) ` + baseQuery
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Apply pagination and order
	query := `SELECT m.*, u.name AS verifier_name ` + baseQuery + fmt.Sprintf(` ORDER BY m.surah_number ASC, m.ayat_start ASC LIMIT $%d OFFSET $%d`, idx, idx+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&list, query, args...)
	return list, total, err
}

func (r *memorizeRepository) GetStudentSurahDetail(studentID uuid.UUID, surahNumber int) ([]Memorize, error) {
	var list []Memorize
	query := `
		SELECT m.*, u.name AS verifier_name 
		FROM memorize m 
		LEFT JOIN users u ON m.verified_by = u.id 
		WHERE m.student_id = $1 AND m.surah_number = $2 
		ORDER BY m.ayat_start ASC
	`
	err := r.db.Select(&list, query, studentID, surahNumber)
	return list, err
}

func (r *memorizeRepository) BulkUpdateStatus(ids []uuid.UUID, status string, verifiedBy uuid.UUID) error {
	query := `
		UPDATE memorize 
		SET status = $1, verified_by = $2, updated_at = NOW() 
		WHERE id = ANY($3)
	`
	_, err := r.db.Exec(query, status, verifiedBy, pq.Array(ids))
	return err
}
