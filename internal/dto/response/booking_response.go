package response

import (
	"hotel-booking-api/internal/domain"
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID         uuid.UUID        `json:"id"`
	UserID     uuid.UUID        `json:"user_id"`
	Room       RoomResponse     `json:"room"`
	CheckIn    time.Time        `json:"check_in"`
	CheckOut   time.Time        `json:"check_out"`
	TotalPrice float64          `json:"total_price"`
	Status     string           `json:"status"`
	Payment    *PaymentResponse `json:"payment,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
}

type PaymentResponse struct {
	ID            uuid.UUID `json:"id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id,omitempty"`
	PaymentMethod string    `json:"payment_method,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToBookingResponse(booking *domain.Booking) BookingResponse {
	resp := BookingResponse{
		ID:         booking.ID,
		UserID:     booking.UserID,
		Room:       ToRoomResponse(&booking.Room),
		CheckIn:    booking.CheckIn,
		CheckOut:   booking.CheckOut,
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		CreatedAt:  booking.CreatedAt,
	}

	if booking.Payment != nil {
		payment := ToPaymentResponse(booking.Payment)
		resp.Payment = &payment
	}

	return resp
}

func ToPaymentResponse(payment *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:            payment.ID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
	}
}
