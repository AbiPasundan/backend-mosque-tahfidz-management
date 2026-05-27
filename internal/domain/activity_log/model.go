package activity_log

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	UserID      *uuid.UUID `db:"user_id" json:"user_id,omitempty"`
	UserName    string     `db:"user_name" json:"user_name"`
	Action      string     `db:"action" json:"action"`
	EntityType  string     `db:"entity_type" json:"entity_type"`
	EntityID    *uuid.UUID `db:"entity_id" json:"entity_id,omitempty"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}
