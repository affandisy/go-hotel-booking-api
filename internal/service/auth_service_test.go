package service

import (
	"errors"
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/pkg/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register("Test User", "test@example.com", "password123", "CUSTOMER")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "CUSTOMER", user.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	existingUser := &domain.User{Email: "test@example.com"}
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	user, err := service.Register("Test User", "test@example.com", "password123", "CUSTOMER")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	user, err := service.Register("Test User", "invalid-email", "password123", "CUSTOMER")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid email format", err.Error())
}

func TestAuthService_Register_ShortPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	user, err := service.Register("Test User", "test@example.com", "12345", "CUSTOMER")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "password must be at least 6 characters", err.Error())
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	hashedPassword := "$2a$14$..." // Mock hash
	existingUser := &domain.User{
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     "CUSTOMER",
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	token, user, err := service.Login("test@example.com", "wrongpassword")

	assert.Error(t, err) // Will error due to password mismatch
	assert.Empty(t, token)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	validate := validator.New()
	service := NewAuthService(mockRepo, validate)

	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

	token, user, err := service.Login("notfound@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Empty(t, token)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
