package middleware

import (
	"hotel-booking-api/pkg/jsonres"
	"hotel-booking-api/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	requestID := ""
	if rid := c.Get("requestID"); rid != nil {
		requestID = rid.(string)
	}

	logger.Error("Request error",
		"path", c.Request().URL.Path,
		"method", c.Request().Method,
		"error", err.Error(),
		"request_id", requestID,
	)

	if !c.Response().Committed {
		c.JSON(code, jsonres.ErrorWithRequestID("ERROR", message, nil, requestID))
	}
}
