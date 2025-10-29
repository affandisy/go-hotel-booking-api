package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RoomID     uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	CheckIn    time.Time `gorm:"not null" json:"check_in"`
	CheckOut   time.Time `gorm:"not null" json:"check_out"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	Status     string    `gorm:"not null;default:'PENDING'" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Room    Room     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;" json:"room"`
	Payment *Payment `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;" json:"payment,omitempty"`
}
