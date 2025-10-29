package middleware

import (
	"hotel-booking-api/pkg/logger"
	"hotel-booking-api/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	logger.Error("Request error",
		"path", c.Request().URL.Path,
		"method", c.Request().Method,
		"error", err.Error(),
	)

	if !c.Response().Committed {
		c.JSON(code, response.Error("ERROR", message, nil))
	}
}
