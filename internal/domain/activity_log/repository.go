package activity_log

import (
	"github.com/jmoiron/sqlx"
)

type ActivityLogRepository interface {
	Create(log *ActivityLog) error
	List(page, limit int) ([]ActivityLog, int, error)
}

type activityLogRepository struct {
	db *sqlx.DB
}

func NewActivityLogRepository(db *sqlx.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(log *ActivityLog) error {
	query := `INSERT INTO activity_logs (id, user_id, user_name, action, entity_type, entity_id, description, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, log.ID, log.UserID, log.UserName, log.Action,
		log.EntityType, log.EntityID, log.Description, log.CreatedAt)
	return err
}

func (r *activityLogRepository) List(page, limit int) ([]ActivityLog, int, error) {
	var logs []ActivityLog
	var total int

	baseQuery := `FROM activity_logs`

	// Get total count
	countQuery := `SELECT COUNT(id) ` + baseQuery
	if err := r.db.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query := `SELECT * ` + baseQuery + ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&logs, query, limit, (page-1)*limit)
	return logs, total, err
}
