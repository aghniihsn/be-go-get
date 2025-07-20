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
	requestedUserIDStr := c.Params("user_id")
	currentUser := c.Locals("user").(*models.JWTPayload)

	// Convert string to ObjectID for comparison
	requestedUserID, err := primitive.ObjectIDFromHex(requestedUserIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	// Allow users to see only their own tickets, admin can see all
	if currentUser.Role != "admin" && currentUser.ID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{"user_id": requestedUserID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tikets)
}

func CreateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")

	var input struct {
		JadwalID string `json:"jadwal_id"`
		UserID   string `json:"user_id"`
		Kursi    string `json:"kursi"`
		Status   string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.JadwalID == "" || input.UserID == "" || input.Kursi == "" {
		return c.Status(400).JSON(fiber.Map{"error": "JadwalID, UserID, and Kursi are required"})
	}

	// Convert string IDs to ObjectIDs
	jadwalID, err := primitive.ObjectIDFromHex(input.JadwalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JadwalID format"})
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UserID format"})
	}

	// Verify jadwal exists
	jadwalCollection := config.DB.Collection("jadwals")
	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": jadwalID}).Decode(&jadwal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Jadwal not found"})
	}

	// Verify user exists
	userCollection := config.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	// Set default status if not provided
	status := input.Status
	if status == "" {
		status = "confirmed"
	}

	// Create tiket with ObjectID
	now := time.Now()
	tiket := models.Tiket{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		JadwalID:         jadwalID,
		Kursi:            input.Kursi,
		Status:           status,
		TanggalPembelian: now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	result, err := tiketCollection.InsertOne(context.TODO(), tiket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the inserted ID
	tiket.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Tiket created successfully",
		"tiket":   tiket,
	})
}

func UpdateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	var tiket models.Tiket
	if err := c.BodyParser(&tiket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	update := bson.M{
		"$set": bson.M{
			"user_id":           tiket.UserID,
			"jadwal_id":         tiket.JadwalID,
			"kursi":             tiket.Kursi,
			"status":            tiket.Status,
			"tanggal_pembelian": tiket.TanggalPembelian,
			"updated_at":        time.Now(),
		},
	}

	res, err := tiketCollection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Tiket updated successfully"})
}

func DeleteTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	res, err := tiketCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}
	return c.JSON(fiber.Map{"message": "Tiket deleted successfully"})
}
