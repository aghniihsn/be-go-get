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

func GetTiketByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	tiketCollection := config.DB.Collection("tikets")

	cursor, err := tiketCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	var tikets []models.Tiket
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Parse error"})
	}
	return c.JSON(tikets)
}

func CreateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	jadwalCollection := config.DB.Collection("jadwals")
	userCollection := config.DB.Collection("users") // ⬅️ Tambah ini

	var tiket models.Tiket
	if err := c.BodyParser(&tiket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if tiket.ID == "" || tiket.JadwalID == "" || tiket.Nama == "" || tiket.Email == "" || tiket.Jumlah <= 0 || tiket.UserID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Cek ID unik tiket
	count, _ := tiketCollection.CountDocuments(context.TODO(), bson.M{"id": tiket.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID already exists"})
	}

	// Simpan user jika belum ada
	userExist, _ := userCollection.CountDocuments(context.TODO(), bson.M{"id": tiket.UserID})
	if userExist == 0 {
		_, _ = userCollection.InsertOne(context.TODO(), bson.M{
			"id":    tiket.UserID,
			"nama":  tiket.Nama,
			"email": tiket.Email,
		})
	}

	// Ambil harga dari jadwal
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
