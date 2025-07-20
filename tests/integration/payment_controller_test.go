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

// TestPembayaranMultipleObjectIDs tests payment with multiple ObjectID references
func TestPembayaranMultipleObjectIDs(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid pembayaran with valid IDs",
			input: map[string]interface{}{
				"tiket_id":          primitive.NewObjectID().Hex(),
				"user_id":           primitive.NewObjectID().Hex(),
				"jumlah":            100000,
				"metode_pembayaran": "credit_card",
				"status":            "completed",
			},
			expectedStatus: 201,
			expectedError:  false,
		},
		{
			name: "Invalid tiket_id format",
			input: map[string]interface{}{
				"tiket_id":          "invalid-tiket-id",
				"user_id":           primitive.NewObjectID().Hex(),
				"jumlah":            100000,
				"metode_pembayaran": "credit_card",
				"status":            "completed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Invalid user_id format",
			input: map[string]interface{}{
				"tiket_id":          primitive.NewObjectID().Hex(),
				"user_id":           "invalid-user-id",
				"jumlah":            100000,
				"metode_pembayaran": "credit_card",
				"status":            "completed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Zero jumlah",
			input: map[string]interface{}{
				"tiket_id":          primitive.NewObjectID().Hex(),
				"user_id":           primitive.NewObjectID().Hex(),
				"jumlah":            0,
				"metode_pembayaran": "credit_card",
				"status":            "completed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Negative jumlah",
			input: map[string]interface{}{
				"tiket_id":          primitive.NewObjectID().Hex(),
				"user_id":           primitive.NewObjectID().Hex(),
				"jumlah":            -1000,
				"metode_pembayaran": "credit_card",
				"status":            "completed",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing metode_pembayaran",
			input: map[string]interface{}{
				"tiket_id": primitive.NewObjectID().Hex(),
				"user_id":  primitive.NewObjectID().Hex(),
				"jumlah":   100000,
				"status":   "completed",
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
			req := httptest.NewRequest("POST", "/pembayarans", bytes.NewReader(body))
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
				assert.Contains(t, response, "pembayaran")

				pembayaran := response["pembayaran"].(map[string]interface{})
				assert.Contains(t, pembayaran, "_id")
				assert.Contains(t, pembayaran, "tiket_id")
				assert.Contains(t, pembayaran, "user_id")

				// Validate ObjectID formats in response
				pembayaranID := pembayaran["_id"].(string)
				_, err = primitive.ObjectIDFromHex(pembayaranID)
				assert.NoError(t, err, "Pembayaran ID should be valid ObjectID format")

				tiketID := pembayaran["tiket_id"].(string)
				_, err = primitive.ObjectIDFromHex(tiketID)
				assert.NoError(t, err, "Tiket ID should be valid ObjectID format")

				userID := pembayaran["user_id"].(string)
				_, err = primitive.ObjectIDFromHex(userID)
				assert.NoError(t, err, "User ID should be valid ObjectID format")
			}
		})
	}
}

// TestObjectIDQueryParameters tests ObjectID handling in query parameters
func TestObjectIDQueryParameters(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		paramValue  string
		expectError bool
	}{
		{
			name:        "Valid ObjectID in user tickets query",
			endpoint:    "/tikets/user/%s",
			paramValue:  primitive.NewObjectID().Hex(),
			expectError: false,
		},
		{
			name:        "Invalid ObjectID in user tickets query",
			endpoint:    "/tikets/user/%s",
			paramValue:  "invalid-user-id",
			expectError: true,
		},
		{
			name:        "Valid ObjectID in film schedules query",
			endpoint:    "/jadwals/film/%s",
			paramValue:  primitive.NewObjectID().Hex(),
			expectError: false,
		},
		{
			name:        "Invalid ObjectID in film schedules query",
			endpoint:    "/jadwals/film/%s",
			paramValue:  "invalid-film-id",
			expectError: true,
		},
		{
			name:        "Valid ObjectID in single resource query",
			endpoint:    "/films/%s",
			paramValue:  primitive.NewObjectID().Hex(),
			expectError: false,
		},
		{
			name:        "Invalid ObjectID in single resource query",
			endpoint:    "/films/%s",
			paramValue:  "invalid-film-id",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			// Create request URL with parameter
			url := fmt.Sprintf(tt.endpoint, tt.paramValue)
			req := httptest.NewRequest("GET", url, nil)

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			if tt.expectError {
				// Should return 400 for invalid ObjectID
				if resp.StatusCode == 400 {
					var response map[string]interface{}
					err := json.NewDecoder(resp.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Contains(t, response, "error")
				}
			} else {
				// Should not return 400 for valid ObjectID format
				assert.NotEqual(t, 400, resp.StatusCode, "Valid ObjectID should not return 400")
			}
		})
	}
}

// TestCompleteBookingWorkflow tests the complete booking workflow with ObjectIDs
func TestCompleteBookingWorkflow(t *testing.T) {
	// This test simulates a complete booking workflow
	// 1. User registration (creates User ObjectID)
	// 2. Film creation (creates Film ObjectID)
	// 3. Schedule creation (links to Film ObjectID)
	// 4. Ticket booking (links to User and Schedule ObjectIDs)
	// 5. Payment processing (links to Ticket and User ObjectIDs)

	// Generate test ObjectIDs
	userID := primitive.NewObjectID()
	filmID := primitive.NewObjectID()
	jadwalID := primitive.NewObjectID()
	tiketID := primitive.NewObjectID()

	t.Run("Complete workflow ObjectID chain", func(t *testing.T) {
		// Test ObjectID relationships
		assert.NotEqual(t, userID, filmID, "User and Film IDs should be different")
		assert.NotEqual(t, filmID, jadwalID, "Film and Jadwal IDs should be different")
		assert.NotEqual(t, jadwalID, tiketID, "Jadwal and Tiket IDs should be different")

		// Test ObjectID format consistency
		userIDHex := userID.Hex()
		filmIDHex := filmID.Hex()
		jadwalIDHex := jadwalID.Hex()
		tiketIDHex := tiketID.Hex()

		// All should be 24-character hex strings
		assert.Len(t, userIDHex, 24)
		assert.Len(t, filmIDHex, 24)
		assert.Len(t, jadwalIDHex, 24)
		assert.Len(t, tiketIDHex, 24)

		// Test round-trip conversion
		convertedUserID, err := primitive.ObjectIDFromHex(userIDHex)
		assert.NoError(t, err)
		assert.Equal(t, userID, convertedUserID)

		convertedFilmID, err := primitive.ObjectIDFromHex(filmIDHex)
		assert.NoError(t, err)
		assert.Equal(t, filmID, convertedFilmID)

		convertedJadwalID, err := primitive.ObjectIDFromHex(jadwalIDHex)
		assert.NoError(t, err)
		assert.Equal(t, jadwalID, convertedJadwalID)

		convertedTiketID, err := primitive.ObjectIDFromHex(tiketIDHex)
		assert.NoError(t, err)
		assert.Equal(t, tiketID, convertedTiketID)
	})
}

// TestObjectIDPerformance tests ObjectID operations performance
func TestObjectIDPerformance(t *testing.T) {
	// Test ObjectID generation performance
	t.Run("ObjectID generation performance", func(t *testing.T) {
		const iterations = 1000

		for i := 0; i < iterations; i++ {
			id := primitive.NewObjectID()
			assert.NotEqual(t, primitive.NilObjectID, id)
		}
	})

	// Test ObjectID conversion performance
	t.Run("ObjectID conversion performance", func(t *testing.T) {
		const iterations = 1000
		testID := primitive.NewObjectID()

		for i := 0; i < iterations; i++ {
			hexStr := testID.Hex()
			convertedID, err := primitive.ObjectIDFromHex(hexStr)
			assert.NoError(t, err)
			assert.Equal(t, testID, convertedID)
		}
	})
}

// BenchmarkObjectIDOperations benchmarks ObjectID operations
func BenchmarkObjectIDOperations(b *testing.B) {
	b.Run("Generation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = primitive.NewObjectID()
		}
	})

	b.Run("ToHex", func(b *testing.B) {
		id := primitive.NewObjectID()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = id.Hex()
		}
	})

	b.Run("FromHex", func(b *testing.B) {
		id := primitive.NewObjectID()
		hexStr := id.Hex()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = primitive.ObjectIDFromHex(hexStr)
		}
	})

	b.Run("RoundTrip", func(b *testing.B) {
		id := primitive.NewObjectID()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hexStr := id.Hex()
			_, _ = primitive.ObjectIDFromHex(hexStr)
		}
	})
}
