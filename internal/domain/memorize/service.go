package memorize

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MemorizeService interface {
	CreateMemorize(req *CreateMemorizeRequest, mentorID uuid.UUID) (*MemorizeResponse, error)
	GetByID(id uuid.UUID) (*MemorizeResponse, error)
	UpdateStatus(id uuid.UUID, req *UpdateMemorizeStatusRequest, mentorID uuid.UUID) (*MemorizeResponse, error)
	BulkUpdateStatus(req *BulkUpdateStatusRequest, mentorID uuid.UUID) error
	ListByStudent(studentID uuid.UUID, surah, status string, page, limit int) ([]MemorizeResponse, int, error)
	GetStudentSurahDetail(studentID uuid.UUID, surahNumber int) ([]MemorizeResponse, error)
	DeleteMemorize(id uuid.UUID) error
}

type memorizeService struct {
	repo MemorizeRepository
}

func NewMemorizeService(repo MemorizeRepository) MemorizeService {
	return &memorizeService{repo: repo}
}

func (s *memorizeService) CreateMemorize(req *CreateMemorizeRequest, mentorID uuid.UUID) (*MemorizeResponse, error) {
	id := uuid.New()
	now := time.Now()

	var memorizedAt *time.Time
	if req.Status == "memorized" {
		memorizedAt = &now
	}

	var lastReviewedAt *time.Time
	if req.Status == "murojaah" {
		lastReviewedAt = &now
	}

	m := &Memorize{
		ID:             id,
		StudentID:      req.StudentID,
		VerifiedBy:     &mentorID,
		Surah:          req.Surah,
		SurahNumber:    req.SurahNumber,
		AyatStart:      req.AyatStart,
		AyatEnd:        req.AyatEnd,
		Status:         req.Status,
		Notes:          req.Notes,
		MemorizedAt:    memorizedAt,
		LastReviewedAt: lastReviewedAt,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}

	// Fetch full record to get verifier_name
	fullRecord, err := s.repo.GetByID(id)
	if err != nil {
		return s.toResponse(m), nil
	}

	return s.toResponse(fullRecord), nil
}

func (s *memorizeService) GetByID(id uuid.UUID) (*MemorizeResponse, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("memorize record not found")
	}
	return s.toResponse(m), nil
}

func (s *memorizeService) UpdateStatus(id uuid.UUID, req *UpdateMemorizeStatusRequest, mentorID uuid.UUID) (*MemorizeResponse, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("memorize record not found")
	}

	now := time.Now()
	if req.Status == "memorized" && m.MemorizedAt == nil {
		m.MemorizedAt = &now
	}
	if req.Status == "murojaah" {
		m.LastReviewedAt = &now
	}

	m.VerifiedBy = &mentorID
	m.Notes = req.Notes
	m.Status = req.Status

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}

	// Fetch full record to get updated verifier_name
	fullRecord, err := s.repo.GetByID(id)
	if err != nil {
		return s.toResponse(m), nil
	}

	return s.toResponse(fullRecord), nil
}

func (s *memorizeService) BulkUpdateStatus(req *BulkUpdateStatusRequest, mentorID uuid.UUID) error {
	return s.repo.BulkUpdateStatus(req.IDs, req.Status, mentorID)
}

func (s *memorizeService) ListByStudent(studentID uuid.UUID, surah, status string, page, limit int) ([]MemorizeResponse, int, error) {
	list, total, err := s.repo.ListByStudent(studentID, surah, status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []MemorizeResponse
	for _, item := range list {
		responses = append(responses, *s.toResponse(&item))
	}
	return responses, total, nil
}

func (s *memorizeService) GetStudentSurahDetail(studentID uuid.UUID, surahNumber int) ([]MemorizeResponse, error) {
	list, err := s.repo.GetStudentSurahDetail(studentID, surahNumber)
	if err != nil {
		return nil, err
	}

	var responses []MemorizeResponse
	for _, item := range list {
		responses = append(responses, *s.toResponse(&item))
	}
	return responses, nil
}

func (s *memorizeService) DeleteMemorize(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("memorize record not found")
	}
	return s.repo.Delete(id)
}

func (s *memorizeService) toResponse(m *Memorize) *MemorizeResponse {
	var memorizedAtStr *string
	if m.MemorizedAt != nil {
		sStr := m.MemorizedAt.Format("2006-01-02T15:04:05Z")
		memorizedAtStr = &sStr
	}

	var lastReviewedAtStr *string
	if m.LastReviewedAt != nil {
		sStr := m.LastReviewedAt.Format("2006-01-02T15:04:05Z")
		lastReviewedAtStr = &sStr
	}

	return &MemorizeResponse{
		ID:             m.ID,
		StudentID:      m.StudentID,
		VerifiedBy:     m.VerifiedBy,
		VerifierName:   m.VerifierName,
		Surah:          m.Surah,
		SurahNumber:    m.SurahNumber,
		AyatStart:      m.AyatStart,
		AyatEnd:        m.AyatEnd,
		Status:         m.Status,
		Notes:          m.Notes,
		MemorizedAt:    memorizedAtStr,
		LastReviewedAt: lastReviewedAtStr,
	}
}
