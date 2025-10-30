package service

import (
	"errors"
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/repository"
)

type RoomService interface {
	CreateRoom(room *domain.Room) error
	UpdateRoom(room *domain.Room) error
	GetRoomsByHotel(hotelID string) ([]domain.Room, error)
	GetRoomByID(id string) (*domain.Room, error)
	SearchAvailableRooms(hotelID string, checkIn, checkOut string) ([]domain.Room, error)
	CheckAvailability(roomID string, checkIn, checkOut string) (bool, error)
}

type roomService struct {
	roomRepo    repository.RoomRepository
	hotelRepo   repository.HotelRepository
	bookingRepo repository.BookingRepository
}

func NewRoomService(roomRepo repository.RoomRepository, hotelRepo repository.HotelRepository) RoomService {
	return &roomService{
		roomRepo:  roomRepo,
		hotelRepo: hotelRepo,
	}
}

func (s *roomService) CreateRoom(room *domain.Room) error {
	if _, err := s.hotelRepo.FindByID(room.HotelID.String()); err != nil {
		return errors.New("hotel not found")
	}

	if room.PricePerNight <= 0 {
		return errors.New("price per night must be greater than 0")
	}

	if room.Availability < 0 {
		return errors.New("availability cannot be negative")
	}

	if room.RoomType == "" {
		return errors.New("room type is required")
	}

	return s.roomRepo.Create(room)
}

func (s *roomService) UpdateRoom(room *domain.Room) error {
	existingRoom, err := s.roomRepo.FindByID(room.ID.String())
	if err != nil {
		return errors.New("room not found")
	}

	if room.PricePerNight <= 0 {
		return errors.New("price per night must be greater than 0")
	}

	if room.Availability < 0 {
		return errors.New("availability cannot be negative")
	}

	if room.RoomType == "" {
		return errors.New("room type is required")
	}

	room.HotelID = existingRoom.HotelID

	return s.roomRepo.Update(room)
}

func (s *roomService) GetRoomsByHotel(hotelID string) ([]domain.Room, error) {
	if _, err := s.hotelRepo.FindByID(hotelID); err != nil {
		return nil, errors.New("hotel not found")
	}

	return s.roomRepo.FindByHotel(hotelID)
}

func (s *roomService) GetRoomByID(id string) (*domain.Room, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (s *roomService) SearchAvailableRooms(hotelID string, checkIn, checkOut string) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindByHotel(hotelID)
	if err != nil {
		return nil, err
	}

	var availableRooms []domain.Room
	for _, room := range rooms {
		if room.Availability > 0 {
			availableRooms = append(availableRooms, room)
		}
	}

	return availableRooms, nil
}

func (s *roomService) CheckAvailability(roomID string, checkIn, checkOut string) (bool, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return false, errors.New("room not found")
	}

	if room.Availability <= 0 {
		return false, nil
	}

	return true, nil
}
