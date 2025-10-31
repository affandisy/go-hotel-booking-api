package handler

import (
	"hotel-booking-api/internal/domain"
	"hotel-booking-api/internal/dto/request"
	dto "hotel-booking-api/internal/dto/response"
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/jsonres"
	"hotel-booking-api/pkg/util"
	"hotel-booking-api/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req request.CreateRoomRequest
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

	room := &domain.Room{
		HotelID:       util.ParseUUID(req.HotelID),
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		Availability:  req.Availability,
	}

	if err := h.roomService.CreateRoom(room); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"CREATE_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Room created successfully", dto.ToRoomResponse(room),
	))
}

func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	roomID := c.Param("id")

	var req request.UpdateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if errs := validator.Validate(&req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation error", errs,
		))
	}

	room, err := h.roomService.GetRoomByID(roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonres.Error(
			"NOT_FOUND", "Room not found", nil,
		))
	}

	room.RoomType = req.RoomType
	room.PricePerNight = req.PricePerNight
	room.Availability = req.Availability

	if err := h.roomService.UpdateRoom(room); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"UPDATE_FAILED", "Failed to update room", err.Error(),
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Room updated successfully", dto.ToRoomResponse(room),
	))
}

func (h *RoomHandler) ListRoomsByHotel(c echo.Context) error {
	hotelId := c.Param("hotelId")

	rooms, err := h.roomService.GetRoomsByHotel(hotelId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"FETCH_FAILED", "Failed to fetch room", err.Error(),
		))
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = dto.ToRoomResponse(&room)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Room retrieved successfully", roomResponses,
	))
}

func (h *RoomHandler) GetRoom(c echo.Context) error {
	roomID := c.Param("id")

	room, err := h.roomService.GetRoomByID(roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonres.Error(
			"NOT_FOUND", "Room not found", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Room retrieved successfully", dto.ToRoomResponse(room),
	))
}
