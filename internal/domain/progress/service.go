package progress

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ProgressService interface {
	CreateProgress(req *CreateProgressRequest, mentorID uuid.UUID) (*ProgressResponse, error)
	BulkCreateProgress(req *BulkCreateProgressRequest, mentorID uuid.UUID) ([]ProgressResponse, error)
	UpdateProgress(id uuid.UUID, req *UpdateProgressRequest) (*ProgressResponse, error)
	ListProgress(studentID, date string, page, limit int) ([]ProgressResponse, int, error)
	GetDashboardSummary() (*DashboardSummaryResponse, error)
}

type progressService struct {
	repo ProgressRepository
}

func NewProgressService(repo ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) CreateProgress(req *CreateProgressRequest, mentorID uuid.UUID) (*ProgressResponse, error) {
	id := uuid.New()
	progress := &Progress{
		ID:           id,
		StudentID:    req.StudentID,
		MentorID:     mentorID,
		Surah:        req.Surah,
		Status:       req.Status,
		AyatStart:    req.AyatStart,
		AyatEnd:      req.AyatEnd,
		Notes:        req.Notes,
		ProgressDate: time.Now(),
	}

	if err := s.repo.Create(progress); err != nil {
		return nil, err
	}

	// Fetch full record to get joined mentor_name
	fullProgress, err := s.repo.GetByID(id)
	if err != nil {
		return s.toResponse(progress), nil // Fallback
	}

	return s.toResponse(fullProgress), nil
}

func (s *progressService) BulkCreateProgress(req *BulkCreateProgressRequest, mentorID uuid.UUID) ([]ProgressResponse, error) {
	var responses []ProgressResponse
	for _, item := range req.Items {
		resp, err := s.CreateProgress(&item, mentorID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (s *progressService) UpdateProgress(id uuid.UUID, req *UpdateProgressRequest) (*ProgressResponse, error) {
	progress, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("progress not found")
	}

	progress.Surah = req.Surah
	progress.Status = req.Status
	progress.AyatStart = req.AyatStart
	progress.AyatEnd = req.AyatEnd
	progress.Notes = req.Notes

	if err := s.repo.Update(progress); err != nil {
		return nil, err
	}

	return s.toResponse(progress), nil
}

func (s *progressService) ListProgress(studentID, date string, page, limit int) ([]ProgressResponse, int, error) {
	progressList, total, err := s.repo.List(studentID, date, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []ProgressResponse
	for _, p := range progressList {
		responses = append(responses, *s.toResponse(&p))
	}
	return responses, total, nil
}

func (s *progressService) GetDashboardSummary() (*DashboardSummaryResponse, error) {
	summary, err := s.repo.GetDashboardSummary()
	if err != nil {
		return nil, err
	}

	weeklyActivity, err := s.repo.GetWeeklyActivity()
	if err != nil {
		return nil, err
	}

	recentProgress, err := s.repo.GetRecentProgress(5)
	if err != nil {
		return nil, err
	}

	return &DashboardSummaryResponse{
		TotalStudents:            summary.TotalStudents,
		ActiveToday:              summary.ActiveToday,
		WeeklyProgressPercentage: summary.WeeklyProgressPercentage,
		WeeklyActivity:           weeklyActivity,
		RecentProgress:           recentProgress,
	}, nil
}

func (s *progressService) toResponse(progress *Progress) *ProgressResponse {
	return &ProgressResponse{
		ID:           progress.ID,
		StudentID:    progress.StudentID,
		MentorID:     progress.MentorID,
		MentorName:   progress.MentorName,
		Surah:        progress.Surah,
		Status:       progress.Status,
		AyatStart:    progress.AyatStart,
		AyatEnd:      progress.AyatEnd,
		Notes:        progress.Notes,
		ProgressDate: progress.ProgressDate.Format("2006-01-02"),
	}
}
