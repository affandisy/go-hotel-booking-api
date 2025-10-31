package service

import (
	"errors"
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/repository"
)

type PaymentService interface {
	HandlePaymentCallback(bookingID, transactionID, status string) error
}

type paymentService struct {
	bookingRepo repository.BookingRepository
	paymentRepo repository.PaymentRepository
	roomRepo    repository.RoomRepository
}

func NewPaymentService(bookingRepo repository.BookingRepository, paymentRepo repository.PaymentRepository, roomRepo repository.RoomRepository) PaymentService {
	return &paymentService{
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
		roomRepo:    roomRepo,
	}
}

func (s *paymentService) HandlePaymentCallback(bookingID, transactionID, status string) error {
	payment, err := s.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		return errors.New("payment record not found")
	}

	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if status == "SUCCESS" {
		payment.Status = domain.PaymentStatusSuccess
		payment.TransactionID = transactionID
		booking.Status = domain.BookingStatusConfirmed
	} else {
		payment.Status = domain.PaymentStatusFailed
		booking.Status = domain.BookingStatusCancelled

		room, err := s.roomRepo.FindByID(booking.RoomID.String())
		if err == nil {
			_ = s.roomRepo.UpdateAvailability(booking.RoomID.String(), room.Availability+1)
		}
	}

	payment.TransactionID = transactionID

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	return s.bookingRepo.Update(booking)
}
