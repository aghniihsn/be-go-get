package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllTikets(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tikets)
}

func GetTiketByID(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	var tiket models.Tiket
	err := tiketCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&tiket)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}
	return c.JSON(tiket)
}

func CreateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	jadwalCollection := config.DB.Collection("jadwals")

	var tiket models.Tiket
	if err := c.BodyParser(&tiket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if tiket.ID == "" || tiket.JadwalID == "" || tiket.Nama == "" || tiket.Email == "" || tiket.Jumlah <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required and jumlah must be > 0"})
	}

	// Cek ID unik
	count, _ := tiketCollection.CountDocuments(context.TODO(), bson.M{"id": tiket.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID already exists"})
	}

	// Ambil harga dari koleksi jadwals
	var jadwal models.Jadwal
	err := jadwalCollection.FindOne(context.TODO(), bson.M{"id": tiket.JadwalID}).Decode(&jadwal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JadwalID"})
	}

	tiket.TotalHarga = float64(tiket.Jumlah) * jadwal.Harga

	_, err = tiketCollection.InsertOne(context.TODO(), tiket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Tiket created", "total_harga": tiket.TotalHarga})
}

func UpdateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	var tiket models.Tiket
	if err := c.BodyParser(&tiket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var jadwal models.Jadwal
	err := jadwalCollection.FindOne(context.TODO(), bson.M{"id": tiket.JadwalID}).Decode(&jadwal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid jadwal_id"})
	}

	tiket.TotalHarga = float64(tiket.Jumlah) * jadwal.Harga

	update := bson.M{
		"$set": bson.M{
			"jadwal_id":   tiket.JadwalID,
			"nama":        tiket.Nama,
			"email":       tiket.Email,
			"jumlah":      tiket.Jumlah,
			"total_harga": tiket.TotalHarga,
		},
	}

	res, err := tiketCollection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Tiket updated", "total_harga": tiket.TotalHarga})
}

func DeleteTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	res, err := tiketCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}
	return c.JSON(fiber.Map{"message": "Tiket deleted"})
}
