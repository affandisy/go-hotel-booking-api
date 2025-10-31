package service

import (
	"errors"
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/repository"
	"hotel-booking-api/pkg/util"
	"time"

	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(userID, roomID string, checkIn, checkOut time.Time) (*domain.Booking, error)
	CancelBooking(userID, bookingID string) error
	GetUserBookings(userID string) ([]domain.Booking, error)
}

type bookingService struct {
	DB          *gorm.DB
	bookingRepo repository.BookingRepository
	roomRepo    repository.RoomRepository
	paymentRepo repository.PaymentRepository
}

func NewBookingService(db *gorm.DB, bookingRepo repository.BookingRepository, roomRepo repository.RoomRepository, paymentRepo repository.PaymentRepository) BookingService {
	return &bookingService{
		DB:          db,
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *bookingService) CreateBooking(userID, roomID string, checkIn, checkOut time.Time) (*domain.Booking, error) {
	now := time.Now()
	if checkIn.Before(now) {
		return nil, errors.New("check-in date cannot be in the past")
	}
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.New("check-out date must be after check-in date")
	}

	maxDuration := 30 * 24 * time.Hour
	if checkOut.Sub(checkIn) > maxDuration {
		return nil, errors.New("maximum booking duration is 30 days")
	}

	activeBookings, _ := s.bookingRepo.FindActiveByRoom(roomID, checkIn.Format(time.RFC3339), checkOut.Format(time.RFC3339))
	if len(activeBookings) > 0 {
		return nil, errors.New("room is already booked for selected dates")
	}

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if room.Availability <= 0 {
		return nil, errors.New("room not available")
	}

	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	if nights < 1 {
		nights = 1
	}
	totalPrice := float64(nights) * room.PricePerNight

	booking := &domain.Booking{
		UserID:     util.ParseUUID(userID),
		RoomID:     util.ParseUUID(roomID),
		CheckIn:    checkIn,
		CheckOut:   checkOut,
		TotalPrice: totalPrice,
		Status:     domain.BookingStatusPending,
	}

	txErr := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.bookingRepo.Create(booking); err != nil {
			return err
		}

		if err := s.roomRepo.UpdateAvailability(roomID, room.Availability-1); err != nil {
			return err
		}

		payment := &domain.Payment{
			BookingID: booking.ID,
			Amount:    totalPrice,
			Status:    domain.PaymentStatusPending,
		}
		if err := s.paymentRepo.Create(payment); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	booking, _ = s.bookingRepo.FindByID(booking.ID.String())

	return booking, nil
}

func (s *bookingService) CancelBooking(userID, bookingID string) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.UserID.String() != userID {
		return errors.New("unauthorized to cancel this booking")
	}

	if booking.Status == domain.BookingStatusCancelled {
		return errors.New("booking already cancelled")
	}

	if booking.Status == domain.BookingStatusCompleted {
		return errors.New("cannot cancel completed booking")
	}

	room, err := s.roomRepo.FindByID(booking.RoomID.String())
	if err == nil {
		_ = s.roomRepo.UpdateAvailability(booking.RoomID.String(), room.Availability+1)
	}

	booking.Status = domain.BookingStatusCancelled
	return s.bookingRepo.Update(booking)
}

func (s *bookingService) GetUserBookings(userID string) ([]domain.Booking, error) {
	return s.bookingRepo.FindByUser(userID)
}
