package handler

import (
	"hotel-booking-api/internal/service"
	"hotel-booking-api/pkg/jsonres"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

type WebhookRequest struct {
	BookingID     string `json:"booking_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func (h *PaymentHandler) HandleWebhook(c echo.Context) error {
	var req WebhookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if err := h.paymentService.HandlePaymentCallback(req.BookingID, req.TransactionID, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"WEBHOOK_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Payment webhook successfully", nil,
	))
}
