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

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(echo echo.Context) error {
	var req request.RegisterRequest

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

	if req.Role == "" {
		req.Role = "CUSTOMER"
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return echo.JSON(http.StatusBadRequest, jsonres.Error(
			"REGISTER_FAILED", err.Error(), nil,
		))
	}

	return echo.JSON(http.StatusCreated, jsonres.Success(
		"Registration successfull", dto.ToUserResponse(user),
	))
}

func (h *AuthHandler) Login(echo echo.Context) error {
	var req request.LoginRequest

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

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return echo.JSON(http.StatusUnauthorized, jsonres.Error(
			"LOGIN_FAILED", err.Error(), nil,
		))
	}

	return echo.JSON(http.StatusOK, jsonres.Success("Login Successful", dto.LoginResponse{
		Token: token,
	}))
}
