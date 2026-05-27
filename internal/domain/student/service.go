package student

import (
	"errors"
	"time"

	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/google/uuid"
)

type StudentService interface {
	CreateStudent(req *CreateStudentRequest, mentorID uuid.UUID) (*StudentResponse, error)
	GetStudent(id uuid.UUID) (*StudentResponse, error)
	UpdateStudent(id uuid.UUID, req *UpdateStudentRequest) (*StudentResponse, error)
	DeleteStudent(id uuid.UUID) error
	ListStudents(search, status, learningLevel string, page, limit int) ([]StudentResponse, int, error)
}

type studentService struct {
	repo StudentRepository
}

func NewStudentService(repo StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) CreateStudent(req *CreateStudentRequest, mentorID uuid.UUID) (*StudentResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	student := &Student{
		ID:            uuid.New(),
		MentorID:      mentorID,
		Name:          req.Name,
		Username:      req.Username,
		Password:      hashedPassword,
		LearningLevel: req.LearningLevel,
		Age:           req.Age,
		Status:        req.Status,
		Contact:       req.Contact,
		JoinDate:      time.Now(),
	}

	if err := s.repo.Create(student); err != nil {
		return nil, err
	}

	return s.toResponse(student), nil
}

func (s *studentService) GetStudent(id uuid.UUID) (*StudentResponse, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("student not found")
	}
	return s.toResponse(student), nil
}

func (s *studentService) UpdateStudent(id uuid.UUID, req *UpdateStudentRequest) (*StudentResponse, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Username != "" {
		student.Username = req.Username
	}
	if req.Age != 0 {
		student.Age = req.Age
	}
	if req.LearningLevel != "" {
		student.LearningLevel = req.LearningLevel
	}
	if req.Status != "" {
		student.Status = req.Status
	}
	if req.Contact != "" {
		student.Contact = req.Contact
	}

	if err := s.repo.Update(student); err != nil {
		return nil, err
	}

	return s.toResponse(student), nil
}

func (s *studentService) DeleteStudent(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *studentService) ListStudents(search, status, learningLevel string, page, limit int) ([]StudentResponse, int, error) {
	students, total, err := s.repo.List(search, status, learningLevel, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []StudentResponse
	for _, st := range students {
		responses = append(responses, *s.toResponse(&st))
	}
	return responses, total, nil
}

func (s *studentService) toResponse(student *Student) *StudentResponse {
	return &StudentResponse{
		ID:            student.ID,
		MentorID:      student.MentorID,
		MentorName:    student.MentorName,
		Name:          student.Name,
		Username:      student.Username,
		ProfileImg:    student.ProfileImg,
		CoverImg:      student.CoverImg,
		Age:           student.Age,
		LearningLevel: student.LearningLevel,
		Fluency:       student.Fluency,
		Status:        student.Status,
		Contact:       student.Contact,
		JoinDate:      student.JoinDate.Format("2006-01-02"),
		LastProgress:  formatTime(student.LastProgress),
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
