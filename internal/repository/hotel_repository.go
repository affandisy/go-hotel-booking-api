package repository

import (
	"hotel-booking-api/internal/domain"

	"gorm.io/gorm"
)

type HotelRepository interface {
	Create(hotel *domain.Hotel) error
	Update(hotel *domain.Hotel) error
	FindAll() ([]domain.Hotel, error)
	FindByID(id string) (*domain.Hotel, error)
	Delete(id string) error
}

type hotelRepository struct {
	DB *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{DB: db}
}

func (r *hotelRepository) Create(hotel *domain.Hotel) error {
	return r.DB.Create(hotel).Error
}

func (r *hotelRepository) Update(hotel *domain.Hotel) error {
	return r.DB.Save(hotel).Error
}

func (r *hotelRepository) FindAll() ([]domain.Hotel, error) {
	var hotels []domain.Hotel
	err := r.DB.Preload("Rooms").Find(&hotels).Error

	return hotels, err
}

func (r *hotelRepository) FindByID(id string) (*domain.Hotel, error) {
	var hotel domain.Hotel
	err := r.DB.Preload("Rooms").First(&hotel, "id = ?", id).Error

	return &hotel, err
}

func (r *hotelRepository) Delete(id string) error {
	return r.DB.Delete(&domain.Hotel{}, "id = ?", id).Error
}
