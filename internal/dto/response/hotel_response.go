package response

import (
	"hotel-booking-api/internal/domain"
	"time"

	"github.com/google/uuid"
)

type HotelResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Location    string         `json:"location"`
	Description string         `json:"description"`
	Rooms       []RoomResponse `json:"rooms,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

type RoomResponse struct {
	ID            uuid.UUID     `json:"id"`
	HotelID       uuid.UUID     `json:"hotel_id"`
	Hotel         *HotelSummary `json:"hotel,omitempty"`
	RoomType      string        `json:"room_type"`
	PricePerNight float64       `json:"price_per_night"`
	Availability  int           `json:"availability"`
	CreatedAt     time.Time     `json:"created_at"`
}

type HotelSummary struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
}

func ToHotelResponse(hotel *domain.Hotel) HotelResponse {
	resp := HotelResponse{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Location:    hotel.Location,
		Description: hotel.Description,
		CreatedAt:   hotel.CreatedAt,
	}

	if len(hotel.Room) > 0 {
		rooms := make([]RoomResponse, len(hotel.Room))
		for i, room := range hotel.Room {
			rooms[i] = RoomResponse{
				ID:            room.ID,
				HotelID:       room.HotelID,
				RoomType:      room.RoomType,
				PricePerNight: room.PricePerNight,
				Availability:  room.Availability,
				CreatedAt:     room.CreatedAt,
			}
		}
		resp.Rooms = rooms
	}

	return resp
}

func ToRoomResponse(room *domain.Room) RoomResponse {
	resp := RoomResponse{
		ID:            room.ID,
		HotelID:       room.HotelID,
		RoomType:      room.RoomType,
		PricePerNight: room.PricePerNight,
		Availability:  room.Availability,
		CreatedAt:     room.CreatedAt,
	}

	if room.Hotel.ID != uuid.Nil {
		resp.Hotel = &HotelSummary{
			ID:       room.Hotel.ID,
			Name:     room.Hotel.Name,
			Location: room.Hotel.Location,
		}
	}

	return resp
}
