package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BookingID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"booking_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Status        string    `gorm:"not null;default:'PENDING'" json:"status"`
	TransactionID string    `gorm:"type:varchar(100)" json:"transaction_id"`
	PaymentMethod string    `gorm:"type:varchar(50)" json:"payment_method"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Booking Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;" json:"booking"`
}
