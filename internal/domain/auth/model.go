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

type MentorStudent struct {
	ID            uuid.UUID  `db:"id"`
	MentorID      uuid.UUID  `db:"mentor_id"`
	Name          string     `db:"name"`
	ProfileImg    string     `db:"profile_img"`
	CoverImg      string     `db:"cover_img"`
	Age           int        `db:"age"`
	LearningLevel string     `db:"learning_level"`
	Fluency       string     `db:"fluency"`
	Status        string     `db:"status"`
	Contact       string     `db:"contact"`
	JoinDate      time.Time  `db:"join_date"`
	LastProgress  *time.Time `db:"last_progress"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}
