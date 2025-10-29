package middleware

import (
	"hotel-booking-api/pkg/response"
	"hotel-booking-api/pkg/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echo echo.Context) error {
			authHeader := echo.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.JSON(http.StatusUnauthorized, response.Error(
					"UNAUTHORIZED", "Missing authorization header", nil,
				))
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return echo.JSON(http.StatusUnauthorized, response.Error(
					"UNAUTHORIZED", "Invalid authorization format", nil,
				))
			}

			claims, err := util.ParseJWT(tokenParts[1])
			if err != nil {
				return echo.JSON(http.StatusUnauthorized, response.Error(
					"UNAUTHORIZED", "Invalid token", nil,
				))
			}

			echo.Set("userID", claims.UserID)
			echo.Set("role", claims.Role)

			return next(echo)
		}
	}
}

func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echo echo.Context) error {
			role := echo.Get("role")
			if role != "ADMIN" {
				return echo.JSON(http.StatusForbidden, response.Error(
					"FORBIDDEN", "Admin access required", nil,
				))
			}
			return next(echo)
		}
	}
}
