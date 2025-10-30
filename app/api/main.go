package api

import (
	"context"
	"fmt"
	"hotel-booking-api/app/router"
	"hotel-booking-api/internal/handler"
	"hotel-booking-api/internal/middleware"
	"hotel-booking-api/internal/repository"
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/config"
	"hotel-booking-api/pkg/db"
	"hotel-booking-api/pkg/logger"
	"hotel-booking-api/pkg/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.App.Environment)
	logger.Info("Starting Hotel Booking API", "version", cfg.App.Version)

	db, err := db.InitPostgres(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	logger.Info("Database connected successfully")

	if err := db.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to migrate database", "error", err)
	}

	// Init validator
	validate := validator.New()

	// Init repo
	userRepo := repository.NewUserRepository(db)
	hotelRepo := repository.NewHotelRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Init service
	authService := service.NewAuthService(userRepo, validate)
	hotelService := service.NewHotelService(hotelRepo)
	roomService := service.NewRoomService(roomRepo, hotelRepo)
	bookingService := service.NewBookingService(db, bookingRepo, roomRepo, paymentRepo)
	paymentService := service.NewPaymentService(bookingRepo, paymentRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)
	hotelHandler := handler.NewHotelHandler(hotelService)
	roomHandler := handler.NewRoomHandler(roomService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Init echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Custom error handler
	e.HTTPErrorHandler = middleware.ErrorHandler

	// Global middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.RequestLogger())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// Setup routes
	api := e.Group("/api/v1")
	router.SetupAuthRoutes(api, authHandler)
	router.SetupHotelRoutes(api, hotelHandler, middleware.AuthMiddleware())
	router.SetupRoomRoutes(api, roomHandler, middleware.AuthMiddleware())
	router.SetupBookingRoutes(api, bookingHandler, middleware.AuthMiddleware())
	router.SetupPaymentRoutes(api, paymentHandler)

	// goroutine server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Server starting", "address", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped")
}
