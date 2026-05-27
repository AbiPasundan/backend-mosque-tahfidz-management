package activity_log

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLogService interface {
	LogAction(userID *uuid.UUID, userName, action, entityType string, entityID *uuid.UUID, description string) error
	ListActivityLogs(page, limit int) ([]ActivityLogResponse, int, error)
}

type activityLogService struct {
	repo ActivityLogRepository
}

func NewActivityLogService(repo ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) LogAction(userID *uuid.UUID, userName, action, entityType string, entityID *uuid.UUID, description string) error {
	log := &ActivityLog{
		ID:          uuid.New(),
		UserID:      userID,
		UserName:    userName,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.repo.Create(log)
}

func (s *activityLogService) ListActivityLogs(page, limit int) ([]ActivityLogResponse, int, error) {
	logs, total, err := s.repo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []ActivityLogResponse
	for _, log := range logs {
		responses = append(responses, ActivityLogResponse{
			ID:          log.ID,
			UserID:      log.UserID,
			UserName:    log.UserName,
			Action:      log.Action,
			EntityType:  log.EntityType,
			EntityID:    log.EntityID,
			Description: log.Description,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		})
	}
	return responses, total, nil
}
