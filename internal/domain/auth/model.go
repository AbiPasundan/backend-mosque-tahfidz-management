package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	Password  string     `db:"password"`
	Email     string     `db:"email"`
	Role      string     `db:"role"`
	DeletedAt *time.Time `db:"deleted_at"`
}
