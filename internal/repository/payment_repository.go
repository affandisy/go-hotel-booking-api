package repository

import (
	"hotel-booking-api/internal/domain"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	Update(payment *domain.Payment) error
	FindByBookingID(bookingID string) (*domain.Payment, error)
}

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{DB: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	return r.DB.Save(payment).Error
}

func (r *paymentRepository) FindByBookingID(bookingID string) (*domain.Payment, error) {
	var payment domain.Payment

	err := r.DB.First(&payment, "booking_id = ?", bookingID).Error

	return &payment, err
}
