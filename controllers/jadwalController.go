package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
				"foreignField": "id",
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
	var jadwal models.Jadwal
	err := jadwalCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&jadwal)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}
	return c.JSON(jadwal)
}

func GetJadwalsByFilmID(c *fiber.Ctx) error {
	filmID := c.Params("filmId")
	jadwalCollection := config.DB.Collection("jadwals")

	var jadwals []models.Jadwal
	cursor, err := jadwalCollection.Find(context.TODO(), bson.M{"film_id": filmID})
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
	var jadwal models.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if jadwal.ID == "" || jadwal.FilmID == "" || jadwal.Tanggal == "" || jadwal.Waktu == "" || jadwal.Ruangan == "" || jadwal.Harga <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required and harga must be > 0"})
	}

	count, _ := jadwalCollection.CountDocuments(context.TODO(), bson.M{"id": jadwal.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID already exists"})
	}

	_, err := jadwalCollection.InsertOne(context.TODO(), jadwal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Jadwal created"})
}

func UpdateJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")
	var jadwal models.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	update := bson.M{
		"$set": bson.M{
			"film_id": jadwal.FilmID,
			"tanggal": jadwal.Tanggal,
			"waktu":   jadwal.Waktu,
			"ruangan": jadwal.Ruangan,
			"harga":   jadwal.Harga,
		},
	}
	res, err := jadwalCollection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Jadwal updated"})
}

func DeleteJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")
	res, err := jadwalCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}
	return c.JSON(fiber.Map{"message": "Jadwal deleted"})
}
