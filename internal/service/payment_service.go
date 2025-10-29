package service

import (
	"errors"
	"hotel-booking-api/internal/repository"
)

type PaymentService interface {
	HandlePaymentCallback(bookingID, transactionID, status string) error
}

type paymentService struct {
	bookingRepo repository.BookingRepository
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(bookingRepo repository.BookingRepository, paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *paymentService) HandlePaymentCallback(bookingID, transactionID, status string) error {
	payment, err := s.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		return errors.New("data pembayaran tidak ditemukan")
	}

	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	if status == "SUCCESS" {
		payment.Status = "SUCCESS"
		booking.Status = "CONFIRMED"
	} else {
		payment.Status = "FAILED"
		booking.Status = "CANCELLED"
	}

	payment.TransactionID = transactionID

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	return s.bookingRepo.Update(booking)
}
