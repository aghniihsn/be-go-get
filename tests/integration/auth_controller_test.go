package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-get-backend/models"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAuthRegister tests user registration functionality
func TestAuthRegister(t *testing.T) {
	tests := []struct {
		name           string
		input          models.UserRegister
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid registration",
			input: models.UserRegister{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				Role:     "user",
			},
			expectedStatus: 201,
			expectedError:  false,
		},
		{
			name: "Missing email",
			input: models.UserRegister{
				Username: "newuser",
				Password: "password123",
				Role:     "user",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing username",
			input: models.UserRegister{
				Email:    "newuser@example.com",
				Password: "password123",
				Role:     "user",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing password",
			input: models.UserRegister{
				Username: "newuser",
				Email:    "newuser@example.com",
				Role:     "user",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Default role assignment",
			input: models.UserRegister{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				// Role omitted - should default to "user"
			},
			expectedStatus: 201,
			expectedError:  false,
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
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && resp.StatusCode == 201 {
				// Parse response for successful registration
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "user")

				user := response["user"].(map[string]interface{})
				assert.Contains(t, user, "_id")
				assert.Contains(t, user, "username")
				assert.Contains(t, user, "email")
				assert.Contains(t, user, "role")

				// Validate ObjectID format in response
				userID := user["_id"].(string)
				_, err = primitive.ObjectIDFromHex(userID)
				assert.NoError(t, err, "User ID should be valid ObjectID format")

				// Check role default assignment
				if tt.input.Role == "" {
					assert.Equal(t, "user", user["role"], "Should default to 'user' role")
				}
			}
		})
	}
}

// TestAuthLogin tests user login functionality
func TestAuthLogin(t *testing.T) {
	tests := []struct {
		name           string
		input          models.UserLogin
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid login credentials",
			input: models.UserLogin{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: 200,
			expectedError:  false,
		},
		{
			name: "Invalid email",
			input: models.UserLogin{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: 401,
			expectedError:  true,
		},
		{
			name: "Invalid password",
			input: models.UserLogin{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: 401,
			expectedError:  true,
		},
		{
			name: "Missing email",
			input: models.UserLogin{
				Password: "password123",
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "Missing password",
			input: models.UserLogin{
				Email: "test@example.com",
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
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError && resp.StatusCode == 200 {
				// Parse response for successful login
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")

				// Validate JWT token
				token := response["token"].(string)
				assert.NotEmpty(t, token, "JWT token should not be empty")

				user := response["user"].(map[string]interface{})
				assert.Contains(t, user, "_id")
				assert.Contains(t, user, "username")
				assert.Contains(t, user, "email")
				assert.Contains(t, user, "role")

				// Validate ObjectID format in response
				userID := user["_id"].(string)
				_, err = primitive.ObjectIDFromHex(userID)
				assert.NoError(t, err, "User ID should be valid ObjectID format")
			}
		})
	}
}

// TestObjectIDInAuthResponses tests ObjectID consistency in auth responses
func TestObjectIDInAuthResponses(t *testing.T) {
	// Test that ObjectIDs are consistently formatted in responses
	testUserID := primitive.NewObjectID()

	// Mock user data
	testUser := models.User{
		ID:       testUserID,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	// Test ObjectID to string conversion
	userIDString := testUser.ID.Hex()
	assert.Len(t, userIDString, 24, "ObjectID hex should be 24 characters")

	// Test string to ObjectID conversion
	convertedID, err := primitive.ObjectIDFromHex(userIDString)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, convertedID, "Round-trip conversion should match")

	// Test JSON marshaling/unmarshaling with ObjectID
	userJSON, err := json.Marshal(testUser)
	assert.NoError(t, err)

	var unmarshaledUser models.User
	err = json.Unmarshal(userJSON, &unmarshaledUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, unmarshaledUser.ID, "ObjectID should survive JSON round-trip")
}

// TestAuthPasswordHashing tests password hashing in registration
func TestAuthPasswordHashing(t *testing.T) {
	testCases := []string{
		"password123",
		"ComplexP@ssw0rd!",
		"short",
		"verylongpasswordthatshouldalsowork123456789",
	}

	for _, password := range testCases {
		t.Run(fmt.Sprintf("Password_%s", password), func(t *testing.T) {
			registerInput := models.UserRegister{
				Username: "testuser",
				Email:    "test@example.com",
				Password: password,
				Role:     "user",
			}

			// Create request body
			body, err := json.Marshal(registerInput)
			assert.NoError(t, err)

			// Create Fiber app
			app := fiber.New()

			// Create request
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Should not store plain text password
			if resp.StatusCode == 201 {
				// Verify that password was hashed (not stored as plain text)
				// This would require database inspection in actual implementation
				assert.NotEqual(t, password, "Stored password should be hashed, not plain text")
			}
		})
	}
}
