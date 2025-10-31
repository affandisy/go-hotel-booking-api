package handler

import (
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/dto/request"
	dto "hotel-booking-api/internal/dto/response"
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/jsonres"
	"hotel-booking-api/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HotelHandler struct {
	hotelService service.HotelService
}

func NewHotelHandler(hotelService service.HotelService) *HotelHandler {
	return &HotelHandler{
		hotelService: hotelService,
	}
}

func (h *HotelHandler) CreateHotel(c echo.Context) error {
	var req request.CreateHotelRequest

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

	hotel := &domain.Hotel{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	if err := h.hotelService.CreateHotel(hotel); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"CREATE_FAILED", "Failed to create hotel", err.Error(),
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Hotel created successfully", dto.ToHotelResponse(hotel),
	))
}

func (h *HotelHandler) UpdateHotel(c echo.Context) error {
	hotelID := c.Param("id")

	var req request.UpdateHotelRequest
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

	hotel, err := h.hotelService.GetHotelDetail(hotelID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonres.Error(
			"NOT_FOUND", "Hotel not found", nil,
		))
	}

	hotel.Name = req.Name
	hotel.Location = req.Location
	hotel.Description = req.Description

	if err := h.hotelService.UpdateHotel(hotel); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"UPDATE_FAILED", "Failed to update hotel", err.Error(),
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Hotel updated successfully", dto.ToHotelResponse(hotel),
	))
}

func (h *HotelHandler) ListHotels(c echo.Context) error {
	hotels, err := h.hotelService.ListHotel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"FETCH_FAILED", "Failed to fetch hotels", err.Error(),
		))
	}

	hotelResponse := make([]dto.HotelResponse, len(hotels))
	for i, hotel := range hotels {
		hotelResponse[i] = dto.ToHotelResponse(&hotel)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Hotels retrieved successfully", hotelResponse,
	))
}

func (h *HotelHandler) GetHotel(c echo.Context) error {
	hotelID := c.Param("id")

	hotel, err := h.hotelService.GetHotelDetail(hotelID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonres.Error(
			"NOT_FOUND", "Hotel not found", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Hotel retrieved successfully", dto.ToHotelResponse(hotel),
	))
}

func (h *HotelHandler) DeleteHotel(c echo.Context) error {
	hotelID := c.Param("id")

	if err := h.hotelService.DeleteHotel(hotelID); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"DELETE_FAILED", "Failed to delete hotel", err.Error(),
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Hotel deleted successfully", nil,
	))
}
