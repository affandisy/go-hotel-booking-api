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

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new hotel room booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param request body request.CreateBookingRequest true "Booking details"
// @Success 201 {object} jsonres.SuccessResponse{data=response.BookingResponse}
// @Failure 400 {object} jsonres.ErrorResponse
// @Security BearerAuth
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	userID := c.Get("userID").(string)

	var req request.CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if errs := validator.Validate(&req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation failed", errs,
		))
	}

	booking, err := h.bookingService.CreateBooking(userID, req.RoomID, req.CheckIn, req.CheckOut)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BOOKING_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Booking created successfully", dto.ToBookingResponse(booking),
	))
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} jsonres.SuccessResponse
// @Failure 400 {object} jsonres.ErrorResponse
// @Security BearerAuth
// @Router /bookings/{id}/cancel [patch]
func (h *BookingHandler) CancelBooking(c echo.Context) error {
	userID := c.Get("userID").(string)
	bookingID := c.Param("id")

	if err := h.bookingService.CancelBooking(userID, bookingID); err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"CANCEL_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Booking cancelled successfully", nil,
	))
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get all bookings for the authenticated user
// @Tags bookings
// @Accept json
// @Produce json
// @Success 200 {object} jsonres.SuccessResponse{data=[]response.BookingResponse}
// @Failure 400 {object} jsonres.ErrorResponse
// @Security BearerAuth
// @Router /bookings [get]
func (h *BookingHandler) GetUserBookings(c echo.Context) error {
	userID := c.Get("userID").(string)

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"FETCH_FAILED", "Failed to fetch bookings", err.Error(),
		))
	}

	bookingResponses := make([]dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		bookingResponses[i] = dto.ToBookingResponse(&booking)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully", bookingResponses,
	))
}
