package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAPIHealthCheck tests basic API functionality
func TestAPIHealthCheck(t *testing.T) {
	app := fiber.New()

	// Add a simple health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Cinema API is running",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestObjectIDInAPI tests ObjectID handling in API responses
func TestObjectIDInAPI(t *testing.T) {
	app := fiber.New()

	// Mock endpoint that returns ObjectID
	app.Get("/test-objectid", func(c *fiber.Ctx) error {
		testID := primitive.NewObjectID()
		return c.JSON(fiber.Map{
			"_id": testID,
			"hex": testID.Hex(),
		})
	})

	req := httptest.NewRequest("GET", "/test-objectid", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify ObjectID fields exist
	assert.Contains(t, response, "_id")
	assert.Contains(t, response, "hex")

	// Verify hex string format
	hexStr, ok := response["hex"].(string)
	assert.True(t, ok)
	assert.Len(t, hexStr, 24, "ObjectID hex should be 24 characters")
}

// TestJSONPayloadValidation tests JSON payload validation
func TestJSONPayloadValidation(t *testing.T) {
	app := fiber.New()

	app.Post("/test-payload", func(c *fiber.Ctx) error {
		var payload struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		// Validate ObjectID format
		if payload.ID != "" {
			_, err := primitive.ObjectIDFromHex(payload.ID)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
			}
		}

		return c.JSON(fiber.Map{
			"message": "Payload valid",
			"data":    payload,
		})
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name: "Valid payload with ObjectID",
			payload: map[string]interface{}{
				"id":    primitive.NewObjectID().Hex(),
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatus: 200,
			shouldContain:  "Payload valid",
		},
		{
			name: "Invalid ObjectID format",
			payload: map[string]interface{}{
				"id":    "invalid-objectid",
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatus: 400,
			shouldContain:  "Invalid ObjectID format",
		},
		{
			name: "Valid payload without ObjectID",
			payload: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatus: 200,
			shouldContain:  "Payload valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/test-payload", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			responseStr, _ := json.Marshal(response)
			assert.Contains(t, string(responseStr), tt.shouldContain)
		})
	}
}

// TestCORSHeaders tests CORS configuration
func TestCORSHeaders(t *testing.T) {
	app := fiber.New()

	// Add CORS middleware simulation
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		return c.Next()
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "CORS test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Authorization")
}
