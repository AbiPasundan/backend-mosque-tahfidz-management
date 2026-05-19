package auth

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(page, limit int) ([]User, int, error)
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

func (r *userRepository) List(page, limit int) ([]User, int, error) {
	var users []User
	var total int

	countQuery := `SELECT COUNT(id) FROM users WHERE deleted_at IS NULL`
	if err := r.db.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, email, role FROM users WHERE deleted_at IS NULL ORDER BY name LIMIT $1 OFFSET $2`
	err := r.db.Select(&users, query, limit, (page-1)*limit)
	return users, total, err
}
