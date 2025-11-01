package integration

import (
	"hotel-booking-api/internal/handler"
	"hotel-booking-api/internal/middleware"
	"hotel-booking-api/internal/repository"
	"hotel-booking-api/internal/router"
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/config"
	"hotel-booking-api/pkg/database"
	"hotel-booking-api/pkg/logger"
	"hotel-booking-api/pkg/validator"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
	testE  *echo.Echo
)

func setupTestServer(t *testing.T) (*echo.Echo, func()) {
	os.Setenv("APP_ENV", "test")
	os.Setenv("DB_NAME", "hotel_booking_test")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	logger.Init("test")

	db, err := database.InitPostgres(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	testDB = db

	validate := validator.New()
	userRepo := repository.NewUserRepository(db)
	hotelRepo := repository.NewHotelRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	authService := service.NewAuthService(userRepo, validate)
	hotelService := service.NewHotelService(hotelRepo)
	roomService := service.NewRoomService(roomRepo, hotelRepo)
	bookingService := service.NewBookingService(db, bookingRepo, roomRepo, paymentRepo)
	paymentService := service.NewPaymentService(bookingRepo, paymentRepo, roomRepo)

	authHandler := handler.NewAuthHandler(authService)
	hotelHandler := handler.NewHotelHandler(hotelService)
	roomHandler := handler.NewRoomHandler(roomService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	e := echo.New()
	e.HTTPErrorHandler = middleware.ErrorHandler

	api := e.Group("/api/v1")
	router.SetupAuthRoutes(api, authHandler)
	router.SetupHotelRoutes(api, hotelHandler, middleware.AuthMiddleware())
	router.SetupRoomRoutes(api, roomHandler, middleware.AuthMiddleware())
	router.SetupBookingRoutes(api, bookingHandler, middleware.AuthMiddleware())
	router.SetupPaymentRoutes(api, paymentHandler)

	testE = e

	// Cleanup function
	cleanup := func() {
		// Clean test data
		db.Exec("TRUNCATE TABLE payments CASCADE")
		db.Exec("TRUNCATE TABLE bookings CASCADE")
		db.Exec("TRUNCATE TABLE rooms CASCADE")
		db.Exec("TRUNCATE TABLE hotels CASCADE")
		db.Exec("TRUNCATE TABLE users CASCADE")

		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return e, cleanup
}
