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

func TestHotelOperations_Integration(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	var token string
	var hotelID string

	// Register admin user
	t.Run("Register admin user", func(t *testing.T) {
		reqBody := request.RegisterRequest{
			Name:     "Admin User",
			Email:    "admin@test.com",
			Password: "password123",
			Role:     "ADMIN",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	// Login to get token
	t.Run("Login admin", func(t *testing.T) {
		reqBody := request.LoginRequest{
			Email:    "admin@test.com",
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

		// Extract token from response
		data := resp.Data.(map[string]interface{})
		token = data["token"].(string)
		assert.NotEmpty(t, token)
	})

	// Create hotel
	t.Run("Create hotel", func(t *testing.T) {
		reqBody := request.CreateHotelRequest{
			Name:        "Test Hotel",
			Location:    "Test City",
			Description: "A test hotel",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/hotels", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp jsonres.SuccessResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)

		// Extract hotel ID
		data := resp.Data.(map[string]interface{})
		hotelID = data["id"].(string)
		assert.NotEmpty(t, hotelID)
	})

	// List hotels
	t.Run("List hotels", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hotels", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp jsonres.SuccessResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})

	// Get hotel by ID
	t.Run("Get hotel by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hotels/"+hotelID, nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp jsonres.SuccessResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})
}
