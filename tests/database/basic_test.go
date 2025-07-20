package database_test

import (
	"go-get-backend/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestBasicObjectID test basic ObjectID functionality without database
func TestBasicObjectID(t *testing.T) {
	// Test ObjectID generation
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	// IDs should be different
	if id1 == id2 {
		t.Error("Generated ObjectIDs should be unique")
	}

	// Test ObjectID in struct
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	// ID should not be zero
	if user.ID.IsZero() {
		t.Error("User ID should not be zero")
	}

	// Test ObjectID string conversion
	idString := user.ID.Hex()
	if len(idString) != 24 {
		t.Errorf("ObjectID hex string should be 24 characters, got %d", len(idString))
	}

	// Test ObjectID from string
	parsedID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		t.Errorf("Failed to parse ObjectID from hex: %v", err)
	}

	if parsedID != user.ID {
		t.Error("Parsed ObjectID should match original")
	}
}

// TestModelStructures test that all models have ObjectID fields
func TestModelStructures(t *testing.T) {
	// Test User model
	user := models.User{
		ID: primitive.NewObjectID(),
	}
	if user.ID.IsZero() {
		t.Error("User should have ObjectID field")
	}

	// Test Film model
	film := models.Film{
		ID: primitive.NewObjectID(),
	}
	if film.ID.IsZero() {
		t.Error("Film should have ObjectID field")
	}

	// Test Jadwal model
	jadwal := models.Jadwal{
		ID:     primitive.NewObjectID(),
		FilmID: primitive.NewObjectID(),
	}
	if jadwal.ID.IsZero() || jadwal.FilmID.IsZero() {
		t.Error("Jadwal should have ObjectID fields")
	}

	// Test Tiket model
	tiket := models.Tiket{
		ID:       primitive.NewObjectID(),
		UserID:   primitive.NewObjectID(),
		JadwalID: primitive.NewObjectID(),
	}
	if tiket.ID.IsZero() || tiket.UserID.IsZero() || tiket.JadwalID.IsZero() {
		t.Error("Tiket should have ObjectID fields")
	}

	// Test Pembayaran model
	pembayaran := models.Pembayaran{
		ID:      primitive.NewObjectID(),
		TiketID: primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
	}
	if pembayaran.ID.IsZero() || pembayaran.TiketID.IsZero() || pembayaran.UserID.IsZero() {
		t.Error("Pembayaran should have ObjectID fields")
	}
}
