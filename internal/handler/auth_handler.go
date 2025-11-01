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

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Registration details"
// @Success 201 {object} jsonres.SuccessResponse{data=response.UserResponse}
// @Failure 400 {object} jsonres.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req request.RegisterRequest

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

	if req.Role == "" {
		req.Role = domain.RoleCustomer
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"REGISTER_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Registration successfull", dto.ToUserResponse(user),
	))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} jsonres.SuccessResponse{data=response.LoginResponse}
// @Failure 401 {object} jsonres.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest

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

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, jsonres.Error(
			"LOGIN_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success("Login Successful", dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}))
}
