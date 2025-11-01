package service

import (
	"errors"
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/repository"
	"hotel-booking-api/pkg/util"

	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Register(name, email, password, role string) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	validate *validator.Validate
}

func NewAuthService(userRepo repository.UserRepository, validate *validator.Validate) AuthService {
	return &authService{
		userRepo: userRepo,
		validate: validate,
	}
}

func (s *authService) Register(name, email, password, role string) (*domain.User, error) {
	if err := s.validate.Var(email, "required,email"); err != nil {
		return nil, errors.New("invalid email format")
	}

	if err := s.validate.Var(password, "required,min=6"); err != nil {
		return nil, errors.New("password must be at least 6 characters")
	}

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	if !util.CheckPassword(password, user.Password) {
		return "", nil, errors.New("incorrect password")
	}

	token, err := util.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
