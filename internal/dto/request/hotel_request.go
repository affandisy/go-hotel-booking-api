package request

type CreateHotelRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Location    string `json:"location" validate:"required"`
	Description string `json:"description"`
}

type UpdateHotelRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Location    string `json:"location" validate:"required"`
	Description string `json:"description"`
}

type CreateRoomRequest struct {
	HotelID       string  `json:"hotel_id" validate:"required,uuid4"`
	RoomType      string  `json:"room_type" validate:"required"`
	PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	Availability  int     `json:"availability" validate:"required,gte=0"`
}

type UpdateRoomRequest struct {
	RoomType      string  `json:"room_type" validate:"required"`
	PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	Availability  int     `json:"availability" validate:"required,gte=0"`
}
