package main_test

import (
	"fmt"
	"go-get-backend/config/middleware"
	"go-get-backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestJWTEncoderDecoder tests JWT encoding and decoding with ObjectID
func TestJWTEncoderDecoder(t *testing.T) {
	// Create test payload with ObjectID
	testUserID := primitive.NewObjectID()
	originalPayload := models.JWTPayload{
		ID:       testUserID,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	// Test encoding
	token, err := middleware.Encoder(originalPayload)
	assert.NoError(t, err, "JWT encoding should not return error")
	assert.NotEmpty(t, token, "JWT token should not be empty")

	// Test decoding
	decodedPayload, err := middleware.Decoder(token)
	assert.NoError(t, err, "JWT decoding should not return error")
	assert.NotNil(t, decodedPayload, "Decoded payload should not be nil")

	// Test ObjectID conversion accuracy
	assert.Equal(t, originalPayload.ID, decodedPayload.ID, "ObjectID should match after encoding/decoding")
	assert.Equal(t, originalPayload.Username, decodedPayload.Username, "Username should match")
	assert.Equal(t, originalPayload.Email, decodedPayload.Email, "Email should match")
	assert.Equal(t, originalPayload.Role, decodedPayload.Role, "Role should match")
}

// TestJWTWithInvalidToken tests JWT decoding with invalid tokens
func TestJWTWithInvalidToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid format token",
			token: "invalid.token.format",
		},
		{
			name:  "Malformed JWT",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := middleware.Decoder(tt.token)
			assert.Error(t, err, "Should return error for invalid token: %s", tt.name)
			assert.Nil(t, payload, "Payload should be nil for invalid token")
		})
	}
}

// TestJWTMultipleObjectIDs tests JWT with multiple different ObjectIDs
func TestJWTMultipleObjectIDs(t *testing.T) {
	// Test with multiple different ObjectIDs
	objectIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	for i, objectID := range objectIDs {
		t.Run(fmt.Sprintf("ObjectID_%d", i), func(t *testing.T) {
			payload := models.JWTPayload{
				ID:       objectID,
				Username: fmt.Sprintf("user_%d", i),
				Email:    fmt.Sprintf("user%d@example.com", i),
				Role:     "user",
			}

			// Encode
			token, err := middleware.Encoder(payload)
			assert.NoError(t, err)

			// Decode
			decodedPayload, err := middleware.Decoder(token)
			assert.NoError(t, err)

			// Verify ObjectID accuracy
			assert.Equal(t, objectID, decodedPayload.ID)
			assert.Equal(t, objectID.Hex(), decodedPayload.ID.Hex())
		})
	}
}
