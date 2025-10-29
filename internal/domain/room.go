package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4;primaryKey" json:"id"`
	HotelID       uuid.UUID `gorm:"type:uuid;not null" json:"hotel_id"`
	RoomType      string    `gorm:"not null" json:"room_type"`
	PricePerNight float64   `gorm:"not null" json:"price_per_night"`
	Availability  int       `gorm:"not null;default:1" json:"availability"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Hotel    Hotel     `gorm:"foreignKey:HotelID;constraint:OnDelete:CASCADE;" json:"hotel"`
	Bookings []Booking `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;" json:"bookings,omitempty"`
}
