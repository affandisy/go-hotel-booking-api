package integration

import (
	"bytes"
	"encoding/json"
	"hotel-booking-api/internal/dto/request"
	"hotel-booking-api/pkg/jsonres"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin_Integration(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test Register
	t.Run("Register new user", func(t *testing.T) {
		reqBody := request.RegisterRequest{
			Name:     "Integration Test User",
			Email:    "integration@test.com",
			Password: "password123",
			Role:     "CUSTOMER",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp jsonres.SuccessResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})

	// Test Login
	t.Run("Login with registered user", func(t *testing.T) {
		reqBody := request.LoginRequest{
			Email:    "integration@test.com",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp jsonres.SuccessResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})
}
