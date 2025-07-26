package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllJadwals godoc
// @Summary Get all jadwal
// @Tags Jadwal
// @Produce json
// @Success 200 {array} models.Jadwal
// @Failure 500 {object} map[string]interface{}
// @Router /api/jadwals [get]
func GetAllJadwals(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	var jadwals []models.Jadwal
	cursor, err := jadwalCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &jadwals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(jadwals)
}

func GetAllJadwalsWithFilm(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	cursor, err := jadwalCollection.Aggregate(context.TODO(), bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "films",
				"localField":   "film_id",
				"foreignField": "_id",
				"as":           "film",
			},
		},
		bson.M{
			"$unwind": "$film",
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

// GetJadwalByID godoc
// @Summary Get jadwal by ID
// @Tags Jadwal
// @Produce json
// @Param id path string true "Jadwal ID"
// @Success 200 {object} models.Jadwal
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/jadwals/{id} [get]
func GetJadwalByID(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&jadwal)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}
	return c.JSON(jadwal)
}

func GetJadwalsByFilmID(c *fiber.Ctx) error {
	filmID := c.Params("filmId")

	// Convert string FilmID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FilmID format"})
	}

	jadwalCollection := config.DB.Collection("jadwals")

	var jadwals []models.Jadwal
	cursor, err := jadwalCollection.Find(context.TODO(), bson.M{"film_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &jadwals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(jadwals)
}

// CreateJadwal godoc
// @Summary Create new jadwal
// @Tags Jadwal
// @Accept json
// @Produce json
// @Param jadwal body models.Jadwal true "Jadwal object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/jadwals [post]
func CreateJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")

	var input struct {
		FilmID  string  `json:"film_id"`
		Tanggal string  `json:"tanggal"`
		Waktu   string  `json:"waktu"`
		Ruangan string  `json:"ruangan"`
		Harga   float64 `json:"harga"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi field
	if input.FilmID == "" || input.Tanggal == "" || input.Waktu == "" || input.Ruangan == "" || input.Harga <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required and harga must be > 0"})
	}

	// Validasi format tanggal (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", input.Tanggal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tanggal harus format YYYY-MM-DD"})
	}
	// Validasi format waktu (HH:MM)
	if _, err := time.Parse("15:04", input.Waktu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Waktu harus format HH:MM"})
	}
	// Validasi ruangan
	validStudio := false
	for _, studio := range models.ValidStudios {
		if input.Ruangan == studio {
			validStudio = true
			break
		}
	}
	if !validStudio {
		return c.Status(400).JSON(fiber.Map{"error": "Ruangan harus salah satu dari Studio 1-5"})
	}

	// Convert FilmID string to ObjectID
	filmID, err := primitive.ObjectIDFromHex(input.FilmID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FilmID format"})
	}

	// Ambil judul film
	filmCollection := config.DB.Collection("films")
	var film models.Film
	err = filmCollection.FindOne(context.TODO(), bson.M{"_id": filmID}).Decode(&film)
	if err != nil {
		// Coba cari berdasarkan title jika gagal dengan ObjectID
		errTitle := filmCollection.FindOne(context.TODO(), bson.M{"title": input.FilmID}).Decode(&film)
		if errTitle != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Film not found"})
		}
	}

	now := time.Now()
	jadwal := models.Jadwal{
		ID:        primitive.NewObjectID(),
		FilmID:    filmID,
		FilmTitle: film.Title,
		Tanggal:   input.Tanggal,
		Waktu:     input.Waktu,
		Ruangan:   input.Ruangan,
		Harga:     input.Harga,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := jadwalCollection.InsertOne(context.TODO(), jadwal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	jadwal.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Jadwal created successfully",
		"jadwal":  jadwal,
	})
}

// UpdateJadwal godoc
// @Summary Update jadwal by ID
// @Tags Jadwal
// @Accept json
// @Produce json
// @Param id path string true "Jadwal ID"
// @Param jadwal body models.Jadwal true "Jadwal object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/jadwals/{id} [put]
func UpdateJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var input struct {
		FilmID  string  `json:"film_id"`
		Tanggal string  `json:"tanggal"`
		Waktu   string  `json:"waktu"`
		Ruangan string  `json:"ruangan"`
		Harga   float64 `json:"harga"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi field
	if input.FilmID == "" || input.Tanggal == "" || input.Waktu == "" || input.Ruangan == "" || input.Harga <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required and harga must be > 0"})
	}
	// Validasi format tanggal (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", input.Tanggal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tanggal harus format YYYY-MM-DD"})
	}
	// Validasi format waktu (HH:MM)
	if _, err := time.Parse("15:04", input.Waktu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Waktu harus format HH:MM"})
	}
	// Validasi ruangan
	validStudio := false
	for _, studio := range models.ValidStudios {
		if input.Ruangan == studio {
			validStudio = true
			break
		}
	}
	if !validStudio {
		return c.Status(400).JSON(fiber.Map{"error": "Ruangan harus salah satu dari Studio 1-5"})
	}
	// Convert FilmID string to ObjectID
	filmID, err := primitive.ObjectIDFromHex(input.FilmID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FilmID format"})
	}
	// Ambil judul film
	filmCollection := config.DB.Collection("films")
	var film models.Film
	err = filmCollection.FindOne(context.TODO(), bson.M{"_id": filmID}).Decode(&film)
	if err != nil {
		// Coba cari berdasarkan title jika gagal dengan ObjectID
		errTitle := filmCollection.FindOne(context.TODO(), bson.M{"title": input.FilmID}).Decode(&film)
		if errTitle != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Film not found"})
		}
	}
	update := bson.M{
		"$set": bson.M{
			"film_id":    filmID,
			"film_title": film.Title,
			"tanggal":    input.Tanggal,
			"waktu":      input.Waktu,
			"ruangan":    input.Ruangan,
			"harga":      input.Harga,
			"updated_at": time.Now(),
		},
	}
	res, err := jadwalCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Jadwal updated"})
}

// DeleteJadwal godoc
// @Summary Delete jadwal by ID
// @Tags Jadwal
// @Produce json
// @Param id path string true "Jadwal ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/jadwals/{id} [delete]
func DeleteJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Cascade delete: hapus semua tiket yang referensi ke jadwal ini
	tiketCollection := config.DB.Collection("tikets")
	tiketRes, err := tiketCollection.DeleteMany(context.TODO(), bson.M{"jadwal_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus tiket terkait jadwal", "details": err.Error()})
	}

	// Hapus jadwal
	res, err := jadwalCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}
	return c.JSON(fiber.Map{
		"message":        "Jadwal deleted beserta tiket terkait",
		"deleted_tikets": tiketRes.DeletedCount,
	})
}

// ValidateJadwalID godoc
// @Summary Validasi ID jadwal untuk debugging
// @Tags Jadwal
// @Produce json
// @Param id path string true "Jadwal ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/jadwals/check/{id} [get]
func ValidateJadwalID(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Coba konversi string ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid ID format",
			"message": err.Error(),
			"id":      id,
		})
	}

	// Coba cari jadwal dengan _id
	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&jadwal)
	if err != nil {
		// Jika tidak ditemukan dengan _id, coba cari dengan properti lain untuk debug
		return c.Status(404).JSON(fiber.Map{
			"error":   "Jadwal not found with _id",
			"id":      id,
			"message": err.Error(),
		})
	}

	// Kembalikan informasi jadwal untuk debugging
	return c.JSON(fiber.Map{
		"message":       "Jadwal found",
		"jadwal_id":     id,
		"jadwal_object": jadwal,
		"film_id":       jadwal.FilmID.Hex(),
		"film_title":    jadwal.FilmTitle,
	})
}
