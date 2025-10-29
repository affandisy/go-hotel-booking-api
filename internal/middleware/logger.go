package middleware

import (
	"hotel-booking-api/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			requestID := uuid.New().String()

			c.Set("requestID", requestID)
			c.Response().Header().Set("X-Request-ID", requestID)

			err := next(c)

			duration := time.Since(start)
			logger.Info("Request processed",
				"request_id", requestID,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration_ms", duration.Milliseconds(),
				"ip", c.RealIP(),
			)

			return err
		}
	}
}
