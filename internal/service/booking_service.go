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
	activeBookings, _ := s.bookingRepo.FindActiveByRoom(roomID, checkIn.Format(time.RFC3339), checkOut.Format(time.RFC3339))
	if len(activeBookings) > 0 {
		return nil, errors.New("kamar sudah dibooking untuk tanggal tersebut")
	}

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("room tidak ditemukan")
	}
	if room.Availability <= 0 {
		return nil, errors.New("kamar tidak tersedia")
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
		Status:     "PENDING",
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
			Status:    "PENDING",
		}
		if err := s.paymentRepo.Create(payment); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return booking, nil
}

func (s *bookingService) CancelBooking(userID, bookingID string) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return err
	}

	if booking.UserID.String() != userID {
		return errors.New("tidak diizinkan membatalkan booking")
	}

	booking.Status = "CANCELLED"
	return s.bookingRepo.Update(booking)
}

func (s *bookingService) GetUserBookings(userID string) ([]domain.Booking, error) {
	return s.bookingRepo.FindByUser(userID)
}
