package student

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentRepository interface {
	Create(student *Student) error
	GetByID(id uuid.UUID) (*Student, error)
	Update(student *Student) error
	Delete(id uuid.UUID) error
	List(search, status, learningLevel string, page, limit int) ([]Student, int, error)
}

type studentRepository struct {
	db *sqlx.DB
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *Student) error {
	query := `INSERT INTO students (id, mentor_id, name, password, profile_img, cover_img, age, learning_level, fluency, status, contact, join_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(query, student.ID, student.MentorID, student.Name, student.Password,
		student.ProfileImg, student.CoverImg, student.Age, student.LearningLevel,
		student.Fluency, student.Status, student.Contact, student.JoinDate)
	return err
}

func (r *studentRepository) GetByID(id uuid.UUID) (*Student, error) {
	var student Student
	query := `
		SELECT s.*, u.name AS mentor_name 
		FROM students s
		LEFT JOIN users u ON s.mentor_id = u.id
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	err := r.db.Get(&student, query, id)
	return &student, err
}

func (r *studentRepository) Update(student *Student) error {
	query := `UPDATE students SET name = $2, age = $3, learning_level = $4, fluency = $5, status = $6, contact = $7, updated_at = NOW()
	          WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, student.ID, student.Name, student.Age,
		student.LearningLevel, student.Fluency, student.Status, student.Contact)
	return err
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	query := `UPDATE students SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *studentRepository) List(search, status, learningLevel string, page, limit int) ([]Student, int, error) {
	var students []Student
	var total int

	baseQuery := `
		FROM students s
		LEFT JOIN users u ON s.mentor_id = u.id
		WHERE s.deleted_at IS NULL
	`
	args := []interface{}{}
	idx := 1

	if search != "" {
		baseQuery += fmt.Sprintf(` AND s.name ILIKE $%d`, idx)
		args = append(args, "%"+search+"%")
		idx++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(` AND s.status = $%d`, idx)
		args = append(args, status)
		idx++
	}
	if learningLevel != "" {
		baseQuery += fmt.Sprintf(` AND s.learning_level = $%d`, idx)
		args = append(args, learningLevel)
		idx++
	}

	// Get total count
	countQuery := `SELECT COUNT(s.id) ` + baseQuery
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query := `SELECT s.*, u.name AS mentor_name ` + baseQuery + fmt.Sprintf(` ORDER BY s.created_at DESC LIMIT $%d OFFSET $%d`, idx, idx+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&students, query, args...)
	return students, total, err
}
