package request

import "time"

type CreateBookingRequest struct {
	RoomID   string    `json:"room_id" validate:"required,uuid4"`
	CheckIn  time.Time `json:"check_in" validate:"required"`
	CheckOut time.Time `json:"check_out" validate:"required,gtfield=CheckIn"`
}
