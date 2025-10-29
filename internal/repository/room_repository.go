package repository

import (
	"hotel-booking-api/internal/domain"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *domain.Room) error
	Update(room *domain.Room) error
	FindByHotel(hotelID string) ([]domain.Room, error)
	FindByID(id string) (*domain.Room, error)
	UpdateAvailability(id string, availability int) error
}

type roomRepository struct {
	DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{DB: db}
}

func (r *roomRepository) Create(room *domain.Room) error {
	return r.DB.Create(room).Error
}

func (r *roomRepository) Update(room *domain.Room) error {
	return r.DB.Save(room).Error
}

func (r *roomRepository) FindByHotel(hotelID string) ([]domain.Room, error) {
	var rooms []domain.Room
	err := r.DB.Where("hotel_id = ?", hotelID).Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) FindByID(id string) (*domain.Room, error) {
	var room domain.Room
	err := r.DB.First(&room, "id = ?", id).Error

	return &room, err
}

func (r *roomRepository) UpdateAvailability(id string, availability int) error {
	return r.DB.Model(&domain.Room{}).Where("id = ?", id).Update("availability", availability).Error
}
