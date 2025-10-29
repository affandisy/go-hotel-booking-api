package service

import (
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/repository"
)

type HotelService interface {
	CreateHotel(hotel *domain.Hotel) error
	UpdateHotel(hotel *domain.Hotel) error
	ListHotel() ([]domain.Hotel, error)
	GetHotelDetail(id string) (*domain.Hotel, error)
	DeleteHotel(id string) error
}

type hotelService struct {
	hotelRepo repository.HotelRepository
}

func NewHotelService(repo repository.HotelRepository) HotelService {
	return &hotelService{hotelRepo: repo}
}

func (s *hotelService) CreateHotel(hotel *domain.Hotel) error {
	return s.hotelRepo.Create(hotel)
}

func (s *hotelService) UpdateHotel(hotel *domain.Hotel) error {
	return s.hotelRepo.Update(hotel)
}

func (s *hotelService) ListHotel() ([]domain.Hotel, error) {
	return s.hotelRepo.FindAll()
}

func (s *hotelService) GetHotelDetail(id string) (*domain.Hotel, error) {
	return s.hotelRepo.FindByID(id)
}

func (s *hotelService) DeleteHotel(id string) error {
	return s.hotelRepo.Delete(id)
}
