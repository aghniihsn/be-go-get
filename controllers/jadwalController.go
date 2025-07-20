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

	if input.FilmID == "" || input.Tanggal == "" || input.Waktu == "" || input.Ruangan == "" || input.Harga <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required and harga must be > 0"})
	}

	// Convert FilmID string to ObjectID
	filmID, err := primitive.ObjectIDFromHex(input.FilmID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FilmID format"})
	}

	// Verify film exists
	filmCollection := config.DB.Collection("films")
	var film models.Film
	err = filmCollection.FindOne(context.TODO(), bson.M{"_id": filmID}).Decode(&film)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Film not found"})
	}

	// Create jadwal with ObjectID
	now := time.Now()
	jadwal := models.Jadwal{
		ID:        primitive.NewObjectID(),
		FilmID:    filmID,
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

	// Set the inserted ID
	jadwal.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Jadwal created successfully",
		"jadwal":  jadwal,
	})
}

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

	// Convert FilmID string to ObjectID
	filmID, err := primitive.ObjectIDFromHex(input.FilmID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FilmID format"})
	}

	update := bson.M{
		"$set": bson.M{
			"film_id":    filmID,
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

func DeleteJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	res, err := jadwalCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}
	return c.JSON(fiber.Map{"message": "Jadwal deleted"})
}
