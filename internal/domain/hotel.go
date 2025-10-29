package domain

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `gorm:"not null" json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Room []Room `gorm:"foreignKey:HotelID;constraint:OnDelete:CASCADE;" json:"rooms,omitempty"`
}
