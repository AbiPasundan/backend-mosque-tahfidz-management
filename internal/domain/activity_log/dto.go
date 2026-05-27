package activity_log

import (
	"github.com/google/uuid"
)

type ActivityLogResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	UserName    string     `json:"user_name"`
	Action      string     `json:"action"`
	EntityType  string     `json:"entity_type"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"created_at"`
}
