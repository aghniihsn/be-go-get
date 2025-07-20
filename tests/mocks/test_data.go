package mocks

import (
	"go-get-backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper functions for creating test data
func CreateTestUser() models.User {
	return models.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Role:     "user",
	}
}

func CreateTestUserRegister() models.UserRegister {
	return models.UserRegister{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
}

func CreateTestUserLogin() models.UserLogin {
	return models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}
}

func CreateTestObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func CreateValidObjectIDString() string {
	return primitive.NewObjectID().Hex()
}

func CreateInvalidObjectIDString() string {
	return "invalid-object-id"
}

func CreateTestFilm() models.Film {
	return models.Film{
		ID:          primitive.NewObjectID(),
		Title:       "Test Movie",
		Genre:       "Action",
		Duration:    120,
		Rating:      "PG-13",
		Description: "Test movie description",
		PosterURL:   "https://example.com/poster.jpg",
	}
}

func CreateTestJadwal() models.Jadwal {
	return models.Jadwal{
		ID:      primitive.NewObjectID(),
		FilmID:  primitive.NewObjectID(),
		Tanggal: "2024-01-15",
		Waktu:   "19:00",
		Ruangan: "Cinema 1",
		Harga:   50000,
	}
}
