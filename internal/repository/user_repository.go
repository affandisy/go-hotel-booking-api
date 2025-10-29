package repository

import (
	"hotel-booking-api/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
