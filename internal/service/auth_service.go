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
	Login(email, password string) (string, error)
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
		return nil, errors.New("email tidak valid")
	}

	if err := s.validate.Var(password, "required,min=6"); err != nil {
		return nil, errors.New("password minimal 6 karakter")
	}

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hash, _ := util.HashPassword(password)
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

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !util.CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	token, err := util.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
