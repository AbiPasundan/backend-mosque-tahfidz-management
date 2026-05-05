package progress

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ProgressService interface {
	CreateProgress(req *CreateProgressRequest, mentorID uuid.UUID) (*ProgressResponse, error)
	UpdateProgress(id uuid.UUID, req *UpdateProgressRequest) (*ProgressResponse, error)
	ListProgress(studentID, date string) ([]ProgressResponse, error)
	GetDashboardSummary() (*DashboardSummaryResponse, error)
}

type progressService struct {
	repo ProgressRepository
}

func NewProgressService(repo ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) CreateProgress(req *CreateProgressRequest, mentorID uuid.UUID) (*ProgressResponse, error) {
	progress := &Progress{
		ID:           uuid.New(),
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

	return s.toResponse(progress), nil
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

func (s *progressService) ListProgress(studentID, date string) ([]ProgressResponse, error) {
	progressList, err := s.repo.List(studentID, date)
	if err != nil {
		return nil, err
	}

	var responses []ProgressResponse
	for _, p := range progressList {
		responses = append(responses, *s.toResponse(&p))
	}
	return responses, nil
}

func (s *progressService) GetDashboardSummary() (*DashboardSummaryResponse, error) {
	summary, err := s.repo.GetDashboardSummary()
	if err != nil {
		return nil, err
	}

	return &DashboardSummaryResponse{
		TotalStudents:            summary.TotalStudents,
		ActiveToday:              summary.ActiveToday,
		WeeklyProgressPercentage: summary.WeeklyProgressPercentage,
	}, nil
}

func (s *progressService) toResponse(progress *Progress) *ProgressResponse {
	return &ProgressResponse{
		ID:           progress.ID,
		StudentID:    progress.StudentID,
		MentorID:     progress.MentorID,
		Surah:        progress.Surah,
		Status:       progress.Status,
		AyatStart:    progress.AyatStart,
		AyatEnd:      progress.AyatEnd,
		Notes:        progress.Notes,
		ProgressDate: progress.ProgressDate.Format("2006-01-02"),
	}
}
