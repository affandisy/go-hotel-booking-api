package router

import (
	"hotel-booking-api/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupAuthRoutes(api *echo.Group, handler *handler.AuthHandler) {
	auth := api.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
}

func SetupHotelRoutes(api *echo.Group, handler *handler.HotelHandler, auth echo.MiddlewareFunc) {
	hotels := api.Group("/hotels")

	// Public routes
	hotels.GET("", handler.ListHotels)
	hotels.GET("/:id", handler.GetHotel)

	// Protected routes
	hotels.POST("", handler.CreateHotel, auth)
	hotels.PUT("/:id", handler.UpdateHotel, auth)
	hotels.DELETE("/:id", handler.DeleteHotel, auth)
}

func SetupRoomRoutes(api *echo.Group, handler *handler.RoomHandler, auth echo.MiddlewareFunc) {
	rooms := api.Group("/rooms")

	// Public routes
	rooms.GET("/hotel/:hotelId", handler.ListRoomsByHotel)
	rooms.GET("/:id", handler.GetRoom)

	// Protected routes
	rooms.POST("", handler.CreateRoom, auth)
	rooms.PUT("/:id", handler.UpdateRoom, auth)
}

func SetupBookingRoutes(api *echo.Group, handler *handler.BookingHandler, auth echo.MiddlewareFunc) {
	bookings := api.Group("/bookings", auth)

	// Protected routes
	bookings.POST("", handler.CreateBooking)
	bookings.GET("", handler.GetUserBookings)
	bookings.PATCH("/:id/cancel", handler.CancelBooking)
}

func SetupPaymentRoutes(api *echo.Group, handler *handler.PaymentHandler) {
	payments := api.Group("/payments")
	payments.POST("/webhook", handler.HandleWebhook)
}
