package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestControllerObjectIDValidation tests ObjectID validation logic used in controllers
func TestControllerObjectIDValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Valid ObjectID string",
			input:     primitive.NewObjectID().Hex(),
			expectErr: false,
		},
		{
			name:      "Invalid ObjectID - too short",
			input:     "507f1f77bcf86cd79943901",
			expectErr: true,
		},
		{
			name:      "Invalid ObjectID - too long",
			input:     "507f1f77bcf86cd7994390111",
			expectErr: true,
		},
		{
			name:      "Invalid ObjectID - non-hex",
			input:     "507f1f77bcf86cd79943901g",
			expectErr: true,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is the validation logic used in controllers
			objectID, err := primitive.ObjectIDFromHex(tt.input)

			if tt.expectErr {
				assert.Error(t, err, "Should return error for invalid ObjectID")
				assert.Equal(t, primitive.NilObjectID, objectID)
			} else {
				assert.NoError(t, err, "Should not return error for valid ObjectID")
				assert.NotEqual(t, primitive.NilObjectID, objectID)
			}
		})
	}
}

// TestFilmInputValidation tests film creation input validation logic
func TestFilmInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		genre    string
		duration int
		isValid  bool
	}{
		{
			name:     "Valid film input",
			title:    "Test Movie",
			genre:    "Action",
			duration: 120,
			isValid:  true,
		},
		{
			name:     "Empty title",
			title:    "",
			genre:    "Action",
			duration: 120,
			isValid:  false,
		},
		{
			name:     "Empty genre",
			title:    "Test Movie",
			genre:    "",
			duration: 120,
			isValid:  false,
		},
		{
			name:     "Zero duration",
			title:    "Test Movie",
			genre:    "Action",
			duration: 0,
			isValid:  false,
		},
		{
			name:     "Negative duration",
			title:    "Test Movie",
			genre:    "Action",
			duration: -30,
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the validation logic in CreateFilm controller
			isValid := true

			if tt.title == "" || tt.genre == "" {
				isValid = false
			}
			if tt.duration <= 0 {
				isValid = false
			}

			assert.Equal(t, tt.isValid, isValid, "Validation result should match expected")
		})
	}
}

// TestJadwalInputValidation tests jadwal creation input validation logic
func TestJadwalInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		filmID  string
		tanggal string
		waktu   string
		ruangan string
		harga   float64
		isValid bool
	}{
		{
			name:    "Valid jadwal input",
			filmID:  primitive.NewObjectID().Hex(),
			tanggal: "2024-01-15",
			waktu:   "19:00",
			ruangan: "Cinema 1",
			harga:   50000,
			isValid: true,
		},
		{
			name:    "Invalid FilmID format",
			filmID:  "invalid-film-id",
			tanggal: "2024-01-15",
			waktu:   "19:00",
			ruangan: "Cinema 1",
			harga:   50000,
			isValid: false,
		},
		{
			name:    "Empty FilmID",
			filmID:  "",
			tanggal: "2024-01-15",
			waktu:   "19:00",
			ruangan: "Cinema 1",
			harga:   50000,
			isValid: false,
		},
		{
			name:    "Zero harga",
			filmID:  primitive.NewObjectID().Hex(),
			tanggal: "2024-01-15",
			waktu:   "19:00",
			ruangan: "Cinema 1",
			harga:   0,
			isValid: false,
		},
		{
			name:    "Empty required fields",
			filmID:  primitive.NewObjectID().Hex(),
			tanggal: "",
			waktu:   "",
			ruangan: "",
			harga:   50000,
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the validation logic in CreateJadwal controller
			isValid := true

			// Check required fields
			if tt.filmID == "" || tt.tanggal == "" || tt.waktu == "" || tt.ruangan == "" || tt.harga <= 0 {
				isValid = false
			}

			// Check FilmID ObjectID format
			if tt.filmID != "" {
				_, err := primitive.ObjectIDFromHex(tt.filmID)
				if err != nil {
					isValid = false
				}
			}

			assert.Equal(t, tt.isValid, isValid, "Validation result should match expected")
		})
	}
}

// TestTiketInputValidation tests tiket creation input validation logic
func TestTiketInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		jadwalID string
		kursi    string
		status   string
		isValid  bool
	}{
		{
			name:     "Valid tiket input",
			userID:   primitive.NewObjectID().Hex(),
			jadwalID: primitive.NewObjectID().Hex(),
			kursi:    "A1,A2",
			status:   "confirmed",
			isValid:  true,
		},
		{
			name:     "Invalid UserID format",
			userID:   "invalid-user-id",
			jadwalID: primitive.NewObjectID().Hex(),
			kursi:    "A1,A2",
			status:   "confirmed",
			isValid:  false,
		},
		{
			name:     "Invalid JadwalID format",
			userID:   primitive.NewObjectID().Hex(),
			jadwalID: "invalid-jadwal-id",
			kursi:    "A1,A2",
			status:   "confirmed",
			isValid:  false,
		},
		{
			name:     "Empty kursi",
			userID:   primitive.NewObjectID().Hex(),
			jadwalID: primitive.NewObjectID().Hex(),
			kursi:    "",
			status:   "confirmed",
			isValid:  false,
		},
		{
			name:     "Empty status",
			userID:   primitive.NewObjectID().Hex(),
			jadwalID: primitive.NewObjectID().Hex(),
			kursi:    "A1,A2",
			status:   "",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the validation logic in CreateTiket controller
			isValid := true

			// Check required fields
			if tt.userID == "" || tt.jadwalID == "" || tt.kursi == "" || tt.status == "" {
				isValid = false
			}

			// Check ObjectID formats
			if tt.userID != "" {
				_, err := primitive.ObjectIDFromHex(tt.userID)
				if err != nil {
					isValid = false
				}
			}

			if tt.jadwalID != "" {
				_, err := primitive.ObjectIDFromHex(tt.jadwalID)
				if err != nil {
					isValid = false
				}
			}

			assert.Equal(t, tt.isValid, isValid, "Validation result should match expected")
		})
	}
}

// TestPembayaranInputValidation tests pembayaran creation input validation logic
func TestPembayaranInputValidation(t *testing.T) {
	tests := []struct {
		name             string
		tiketID          string
		userID           string
		jumlah           float64
		metodePembayaran string
		status           string
		isValid          bool
	}{
		{
			name:             "Valid pembayaran input",
			tiketID:          primitive.NewObjectID().Hex(),
			userID:           primitive.NewObjectID().Hex(),
			jumlah:           100000,
			metodePembayaran: "credit_card",
			status:           "completed",
			isValid:          true,
		},
		{
			name:             "Invalid TiketID format",
			tiketID:          "invalid-tiket-id",
			userID:           primitive.NewObjectID().Hex(),
			jumlah:           100000,
			metodePembayaran: "credit_card",
			status:           "completed",
			isValid:          false,
		},
		{
			name:             "Invalid UserID format",
			tiketID:          primitive.NewObjectID().Hex(),
			userID:           "invalid-user-id",
			jumlah:           100000,
			metodePembayaran: "credit_card",
			status:           "completed",
			isValid:          false,
		},
		{
			name:             "Zero jumlah",
			tiketID:          primitive.NewObjectID().Hex(),
			userID:           primitive.NewObjectID().Hex(),
			jumlah:           0,
			metodePembayaran: "credit_card",
			status:           "completed",
			isValid:          false,
		},
		{
			name:             "Negative jumlah",
			tiketID:          primitive.NewObjectID().Hex(),
			userID:           primitive.NewObjectID().Hex(),
			jumlah:           -1000,
			metodePembayaran: "credit_card",
			status:           "completed",
			isValid:          false,
		},
		{
			name:             "Empty metode pembayaran",
			tiketID:          primitive.NewObjectID().Hex(),
			userID:           primitive.NewObjectID().Hex(),
			jumlah:           100000,
			metodePembayaran: "",
			status:           "completed",
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the validation logic in CreatePembayaran controller
			isValid := true

			// Check required fields
			if tt.tiketID == "" || tt.userID == "" || tt.metodePembayaran == "" || tt.status == "" || tt.jumlah <= 0 {
				isValid = false
			}

			// Check ObjectID formats
			if tt.tiketID != "" {
				_, err := primitive.ObjectIDFromHex(tt.tiketID)
				if err != nil {
					isValid = false
				}
			}

			if tt.userID != "" {
				_, err := primitive.ObjectIDFromHex(tt.userID)
				if err != nil {
					isValid = false
				}
			}

			assert.Equal(t, tt.isValid, isValid, "Validation result should match expected")
		})
	}
}

// TestWorkflowObjectIDChain tests ObjectID relationships in complete workflow
func TestWorkflowObjectIDChain(t *testing.T) {
	// Simulate complete booking workflow ObjectID chain
	userID := primitive.NewObjectID()
	filmID := primitive.NewObjectID()
	jadwalID := primitive.NewObjectID()
	tiketID := primitive.NewObjectID()
	pembayaranID := primitive.NewObjectID()

	t.Run("ObjectID uniqueness in workflow", func(t *testing.T) {
		// All ObjectIDs should be unique
		ids := []primitive.ObjectID{userID, filmID, jadwalID, tiketID, pembayaranID}

		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				assert.NotEqual(t, ids[i], ids[j], "All ObjectIDs in workflow should be unique")
			}
		}
	})

	t.Run("ObjectID format consistency", func(t *testing.T) {
		// All ObjectIDs should convert to proper hex format
		userHex := userID.Hex()
		filmHex := filmID.Hex()
		jadwalHex := jadwalID.Hex()
		tiketHex := tiketID.Hex()
		pembayaranHex := pembayaranID.Hex()

		assert.Len(t, userHex, 24)
		assert.Len(t, filmHex, 24)
		assert.Len(t, jadwalHex, 24)
		assert.Len(t, tiketHex, 24)
		assert.Len(t, pembayaranHex, 24)
	})

	t.Run("Round-trip conversion accuracy", func(t *testing.T) {
		// Test conversion accuracy for workflow IDs
		originalIDs := []primitive.ObjectID{userID, filmID, jadwalID, tiketID, pembayaranID}

		for _, originalID := range originalIDs {
			hexStr := originalID.Hex()
			convertedID, err := primitive.ObjectIDFromHex(hexStr)
			assert.NoError(t, err)
			assert.Equal(t, originalID, convertedID, "Round-trip conversion should be accurate")
		}
	})
}
