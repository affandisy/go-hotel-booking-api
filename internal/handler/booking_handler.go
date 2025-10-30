package handler

import (
	"hotel-booking-api/internal/dto/request"
	dto "hotel-booking-api/internal/dto/response"
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/jsonres"
	"hotel-booking-api/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(echo echo.Context) error {
	userID := echo.Get("userID").(string)

	var req request.CreateBookingRequest
	if err := echo.Bind(&req); err != nil {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if errs := validator.Validate(&req); len(errs) > 0 {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation failed", errs,
		))
	}

	booking, err := h.bookingService.CreateBooking(userID, req.RoomID, req.CheckIn, req.CheckOut)
	if err != nil {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"BOOKING_FAILED", err.Error(), nil,
		))
	}

	return echo.JSON(http.StatusCreated, jsonres.Success(
		"Booking created successfully", dto.ToBookingResponse(booking),
	))
}

func (h *BookingHandler) CancelBooking(echo echo.Context) error {
	userID := echo.Get("userID").(string)
	bookingID := echo.Param("id")

	if err := h.bookingService.CancelBooking(userID, bookingID); err != nil {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"CANCEL_FAILED", err.Error(), nil,
		))
	}

	return echo.JSON(http.StatusOK, jsonres.Success(
		"Booking cancelled successfully", nil,
	))
}

func (h *BookingHandler) GetUserBookings(echo echo.Context) error {
	userID := echo.Get("userID").(string)

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"FETCH_FAILED", "Failed to fetch bookings", err.Error(),
		))
	}

	bookingResponses := make([]dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		bookingResponses[i] = dto.ToBookingResponse(&booking)
	}

	return echo.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully", bookingResponses,
	))
}
