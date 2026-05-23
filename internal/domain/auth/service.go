package auth

import (
	"errors"
	"strings"

	"backend-mosque-tahfidz-management/pkg/token"
	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	Login(email, password string) (*LoginResponse, error)
	CreateUser(req *CreateUserRequest) (*UserResponse, error)
	GetUser(id uuid.UUID) (*UserResponse, error)
	UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(id uuid.UUID) error
	ListUsers(search, role string, page, limit int) ([]UserResponse, int, error)
	UpdateProfile(id uuid.UUID, req *UpdateProfileRequest) (*UserResponse, error)
	UpdatePassword(id uuid.UUID, req *UpdatePasswordRequest) error
	GetMentorDetail(id uuid.UUID) (*MentorDetailResponse, error)
}

type authService struct {
	repo       UserRepository
	tokenMaker token.Maker
}

func NewAuthService(repo UserRepository, tokenMaker token.Maker) AuthService {
	return &authService{repo: repo, tokenMaker: tokenMaker}
}

func (s *authService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("user not found")
		return nil, errors.New("invalid email or password")
	}

	match, err := utils.ComparePassword(user.Password, password)
	if err != nil || !match {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.tokenMaker.CreateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, Role: user.Role}, nil
}

func (s *authService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	return &UserResponse{
		ID:     user.ID,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *authService) GetUser(id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &UserResponse{
		ID:     user.ID,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *authService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:     user.ID,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *authService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *authService) ListUsers(search, role string, page, limit int) ([]UserResponse, int, error) {
	users, total, err := s.repo.List(search, role, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []UserResponse
	for _, u := range users {
		responses = append(responses, UserResponse{
			ID:     u.ID,
			UserID: u.ID,
			Name:   u.Name,
			Email:  u.Email,
			Role:   u.Role,
		})
	}
	return responses, total, nil
}

func (s *authService) UpdateProfile(id uuid.UUID, req *UpdateProfileRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:     user.ID,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *authService) UpdatePassword(id uuid.UUID, req *UpdatePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	match, err := utils.ComparePassword(user.Password, req.OldPassword)
	if err != nil || !match {
		return errors.New("old password does not match")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(id, hashedPassword)
}

func (s *authService) GetMentorDetail(id uuid.UUID) (*MentorDetailResponse, error) {
	user, err := s.repo.DetailMentor(id)
	if err != nil {
		return nil, errors.New("mentor not found")
	}

	students, err := s.repo.GetMentorStudents(id)
	if err != nil {
		log.Error().Err(err).Str("mentor_id", id.String()).Msg("failed to fetch mentor students")
	}

	studentResponses := make([]MentorStudentResponse, 0)
	for _, st := range students {
		lastProgress := ""
		if st.LastProgress != nil {
			lastProgress = st.LastProgress.Format("2006-01-02 15:04:05")
		}

		studentResponses = append(studentResponses, MentorStudentResponse{
			ID:            st.ID,
			Name:          st.Name,
			ProfileImg:    st.ProfileImg,
			CoverImg:      st.CoverImg,
			Age:           st.Age,
			LearningLevel: st.LearningLevel,
			Fluency:       st.Fluency,
			Status:        st.Status,
			Contact:       st.Contact,
			JoinDate:      st.JoinDate.Format("2006-01-02"),
			LastProgress:  lastProgress,
		})
	}

	return &MentorDetailResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Students: studentResponses,
	}, nil
}
