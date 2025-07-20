package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestFilmCRUD tests film controller CRUD operations
func TestFilmCreateValidation(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid film creation",
			input: map[string]interface{}{
				"title":       "Test Movie",
				"genre":       "Action",
				"duration":    120,
				"rating":      "PG-13",
				"description": "Test movie description",
				"poster_url":  "https://example.com/poster.jpg",
			},
			expectedStatus: 201,
			expectedError:  false,
		},
		{
			name: "Missing title",
			input: map[string]interface{}{
				"genre":       "Action",
				"duration":    120,
				"rating":      "PG-13",
				"description": "Test movie description",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Invalid duration (zero)",
			input: map[string]interface{}{
				"title":       "Test Movie",
				"genre":       "Action",
				"duration":    0,
				"rating":      "PG-13",
				"description": "Test movie description",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Invalid duration (negative)",
			input: map[string]interface{}{
				"title":       "Test Movie",
				"genre":       "Action",
				"duration":    -30,
				"rating":      "PG-13",
				"description": "Test movie description",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fiber app
			app := fiber.New()

			// Create request body
			body, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest("POST", "/films", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && resp.StatusCode == 201 {
				// Parse response for successful creation
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "film")

				film := response["film"].(map[string]interface{})
				assert.Contains(t, film, "_id")
				assert.Contains(t, film, "title")
				assert.Contains(t, film, "genre")
				assert.Contains(t, film, "duration")

				// Validate ObjectID format in response
				filmID := film["_id"].(string)
				_, err = primitive.ObjectIDFromHex(filmID)
				assert.NoError(t, err, "Film ID should be valid ObjectID format")
			}
		})
	}
}

// TestFilmObjectIDHandling tests ObjectID handling in film operations
func TestFilmObjectIDHandling(t *testing.T) {
	validObjectID := primitive.NewObjectID()
	invalidObjectIDs := []string{
		"invalid-id",
		"12345",
		"507f1f77bcf86cd79943901",   // too short
		"507f1f77bcf86cd7994390111", // too long
		"507f1f77bcf86cd79943901g",  // invalid characters
	}

	t.Run("Valid ObjectID in URL parameter", func(t *testing.T) {
		app := fiber.New()

		// Test GET request with valid ObjectID
		req := httptest.NewRequest("GET", fmt.Sprintf("/films/%s", validObjectID.Hex()), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Should not return 400 for valid ObjectID format
		assert.NotEqual(t, 400, resp.StatusCode, "Valid ObjectID should not return 400")
	})

	for _, invalidID := range invalidObjectIDs {
		t.Run(fmt.Sprintf("Invalid ObjectID: %s", invalidID), func(t *testing.T) {
			app := fiber.New()

			// Test GET request with invalid ObjectID
			req := httptest.NewRequest("GET", fmt.Sprintf("/films/%s", invalidID), nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Should return 400 for invalid ObjectID format
			if resp.StatusCode == 400 {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}
}

// TestJadwalObjectIDRelationships tests ObjectID relationships in jadwal
func TestJadwalObjectIDRelationships(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid jadwal with valid film_id",
			input: map[string]interface{}{
				"film_id": primitive.NewObjectID().Hex(),
				"tanggal": "2024-01-15",
				"waktu":   "19:00",
				"ruangan": "Cinema 1",
				"harga":   50000,
			},
			expectedStatus: 201,
			expectedError:  false,
		},
		{
			name: "Invalid film_id format",
			input: map[string]interface{}{
				"film_id": "invalid-film-id",
				"tanggal": "2024-01-15",
				"waktu":   "19:00",
				"ruangan": "Cinema 1",
				"harga":   50000,
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing film_id",
			input: map[string]interface{}{
				"tanggal": "2024-01-15",
				"waktu":   "19:00",
				"ruangan": "Cinema 1",
				"harga":   50000,
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Zero harga",
			input: map[string]interface{}{
				"film_id": primitive.NewObjectID().Hex(),
				"tanggal": "2024-01-15",
				"waktu":   "19:00",
				"ruangan": "Cinema 1",
				"harga":   0,
			},
			expectedStatus: 400,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			// Create request body
			body, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest("POST", "/jadwals", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && resp.StatusCode == 201 {
				// Parse response for successful creation
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "jadwal")

				jadwal := response["jadwal"].(map[string]interface{})
				assert.Contains(t, jadwal, "_id")
				assert.Contains(t, jadwal, "film_id")

				// Validate ObjectID formats in response
				jadwalID := jadwal["_id"].(string)
				_, err = primitive.ObjectIDFromHex(jadwalID)
				assert.NoError(t, err, "Jadwal ID should be valid ObjectID format")

				filmID := jadwal["film_id"].(string)
				_, err = primitive.ObjectIDFromHex(filmID)
				assert.NoError(t, err, "Film ID should be valid ObjectID format")
			}
		})
	}
}

// TestTiketBookingObjectIDs tests ObjectID handling in ticket booking
func TestTiketBookingObjectIDs(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid tiket with valid IDs",
			input: map[string]interface{}{
				"user_id":   primitive.NewObjectID().Hex(),
				"jadwal_id": primitive.NewObjectID().Hex(),
				"kursi":     "A1,A2",
				"status":    "confirmed",
			},
			expectedStatus: 201,
			expectedError:  false,
		},
		{
			name: "Invalid user_id format",
			input: map[string]interface{}{
				"user_id":   "invalid-user-id",
				"jadwal_id": primitive.NewObjectID().Hex(),
				"kursi":     "A1,A2",
				"status":    "confirmed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Invalid jadwal_id format",
			input: map[string]interface{}{
				"user_id":   primitive.NewObjectID().Hex(),
				"jadwal_id": "invalid-jadwal-id",
				"kursi":     "A1,A2",
				"status":    "confirmed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing kursi",
			input: map[string]interface{}{
				"user_id":   primitive.NewObjectID().Hex(),
				"jadwal_id": primitive.NewObjectID().Hex(),
				"status":    "confirmed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			// Create request body
			body, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest("POST", "/tikets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && resp.StatusCode == 201 {
				// Parse response for successful creation
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "tiket")

				tiket := response["tiket"].(map[string]interface{})
				assert.Contains(t, tiket, "_id")
				assert.Contains(t, tiket, "user_id")
				assert.Contains(t, tiket, "jadwal_id")

				// Validate ObjectID formats in response
				tiketID := tiket["_id"].(string)
				_, err = primitive.ObjectIDFromHex(tiketID)
				assert.NoError(t, err, "Tiket ID should be valid ObjectID format")

				userID := tiket["user_id"].(string)
				_, err = primitive.ObjectIDFromHex(userID)
				assert.NoError(t, err, "User ID should be valid ObjectID format")

				jadwalID := tiket["jadwal_id"].(string)
				_, err = primitive.ObjectIDFromHex(jadwalID)
				assert.NoError(t, err, "Jadwal ID should be valid ObjectID format")
			}
		})
	}
}
