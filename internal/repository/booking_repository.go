package repository

import (
	"hotel-booking-api/internal/domain"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *domain.Booking) error
	Update(booking *domain.Booking) error
	FindByUser(userID string) ([]domain.Booking, error)
	FindByID(id string) (*domain.Booking, error)
	FindActiveByRoom(roomID string, checkIn, checkOut string) ([]domain.Booking, error)
}

type bookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{DB: db}
}

func (r *bookingRepository) Create(booking *domain.Booking) error {
	return r.DB.Create(booking).Error
}

func (r *bookingRepository) Update(booking *domain.Booking) error {
	return r.DB.Save(booking).Error
}

func (r *bookingRepository) FindByUser(userID string) ([]domain.Booking, error) {
	var bookings []domain.Booking

	err := r.DB.Preload("Room.Hotel").Preload("Payment").
		Where("user_id = ?", userID).Order("created_at desc").Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) FindByID(id string) (*domain.Booking, error) {
	var booking domain.Booking

	err := r.DB.Preload("Room.Hotel").Preload("User").Preload("Payment").
		First(&booking, "id = ?", id).Error
	return &booking, err
}

func (r *bookingRepository) FindActiveByRoom(roomID string, checkIn, checkOut string) ([]domain.Booking, error) {
	var bookings []domain.Booking

	err := r.DB.Where("room_id = ? AND status IN ? AND check_in < ? AND check_out > ?",
		roomID, []string{"PENDING", "CONFIRMED"}, checkOut, checkIn).
		Find(&bookings).Error
	return bookings, err
}
