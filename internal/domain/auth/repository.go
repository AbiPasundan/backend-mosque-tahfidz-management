package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(search, role string, page, limit int) ([]User, int, error)
	UpdatePassword(id uuid.UUID, hashedPassword string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	query := `INSERT INTO users (id, name, password, email, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Password, user.Email, user.Role)
	return err
}

func (r *userRepository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	query := `SELECT id, name, password, email, role, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&user, query, id)
	return &user, err
}

func (r *userRepository) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, name, password, email, role, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.Get(&user, query, email)
	return &user, err
}

func (r *userRepository) Update(user *User) error {
	query := `UPDATE users SET name = $2, email = $3, role = $4 WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Role)
	return err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) List(search, role string, page, limit int) ([]User, int, error) {
	var users []User
	var total int

	baseQuery := `FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	idx := 1

	if search != "" {
		baseQuery += fmt.Sprintf(` AND (name ILIKE $%d OR email ILIKE $%d)`, idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}

	if role != "" {
		baseQuery += fmt.Sprintf(` AND role = $%d`, idx)
		args = append(args, role)
		idx++
	}

	countQuery := `SELECT COUNT(id) ` + baseQuery
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, email, role ` + baseQuery + fmt.Sprintf(` ORDER BY name LIMIT $%d OFFSET $%d`, idx, idx+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&users, query, args...)
	return users, total, err
}

func (r *userRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	query := `UPDATE users SET password = $2 WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id, hashedPassword)
	return err
}
