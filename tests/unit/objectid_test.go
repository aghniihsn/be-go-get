package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestObjectIDValidation tests ObjectID string conversion
func TestObjectIDValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid ObjectID string",
			input:   "507f1f77bcf86cd799439011",
			wantErr: false,
		},
		{
			name:    "Invalid ObjectID string - too short",
			input:   "507f1f77bcf86cd79943901",
			wantErr: true,
		},
		{
			name:    "Invalid ObjectID string - non-hex characters",
			input:   "507f1f77bcf86cd79943901g",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectID, err := primitive.ObjectIDFromHex(tt.input)

			if tt.wantErr {
				assert.Error(t, err, "Expected error for test case: %s", tt.name)
				assert.Equal(t, primitive.NilObjectID, objectID)
			} else {
				assert.NoError(t, err, "Expected no error for test case: %s", tt.name)
				assert.NotEqual(t, primitive.NilObjectID, objectID)

				// Test round-trip conversion
				hexString := objectID.Hex()
				assert.Equal(t, tt.input, hexString)
			}
		})
	}
}

// TestObjectIDGeneration tests ObjectID generation and uniqueness
func TestObjectIDGeneration(t *testing.T) {
	// Generate multiple ObjectIDs
	ids := make([]primitive.ObjectID, 10)
	for i := 0; i < 10; i++ {
		ids[i] = primitive.NewObjectID()
	}

	// Test uniqueness
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			assert.NotEqual(t, ids[i], ids[j], "ObjectIDs should be unique")
		}
	}

	// Test hex string format
	for _, id := range ids {
		hexStr := id.Hex()
		assert.Len(t, hexStr, 24, "ObjectID hex string should be 24 characters")

		// Test conversion back
		convertedID, err := primitive.ObjectIDFromHex(hexStr)
		assert.NoError(t, err)
		assert.Equal(t, id, convertedID)
	}
}

// TestObjectIDComparison tests ObjectID comparison operations
func TestObjectIDComparison(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id1Copy := id1

	// Test equality
	assert.Equal(t, id1, id1Copy, "Same ObjectIDs should be equal")
	assert.NotEqual(t, id1, id2, "Different ObjectIDs should not be equal")

	// Test with primitive.NilObjectID
	var nilID primitive.ObjectID
	assert.Equal(t, primitive.NilObjectID, nilID)
	assert.NotEqual(t, id1, nilID)
}
