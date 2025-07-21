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

func GetAllPembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	var data []models.Pembayaran
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &data); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func GetPembayaranByID(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}

	var pembayaran models.Pembayaran
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&pembayaran)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found"})
	}
	return c.JSON(pembayaran)
}

func GetPembayaranByUserID(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	userID := c.Params("user_id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid User ObjectID format"})
	}

	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user payments"})
	}
	defer cursor.Close(context.TODO())

	var pembayarans []models.Pembayaran
	if err = cursor.All(context.TODO(), &pembayarans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode user payments"})
	}

	return c.JSON(pembayarans)
}

func CreatePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")

	var input struct {
		TiketID          string  `json:"tiket_id"`
		UserID           string  `json:"user_id"`
		Jumlah           float64 `json:"jumlah"`
		MetodePembayaran string  `json:"metode_pembayaran"`
		Status           string  `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.TiketID == "" || input.UserID == "" || input.MetodePembayaran == "" || input.Status == "" || input.Jumlah <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "TiketID, UserID, MetodePembayaran, Status, and Jumlah are required and Jumlah > 0"})
	}

	// Convert string IDs to ObjectIDs
	tiketID, err := primitive.ObjectIDFromHex(input.TiketID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid TiketID format"})
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UserID format"})
	}

	// Verify tiket exists
	tiketCollection := config.DB.Collection("tikets")
	var tiket models.Tiket
	err = tiketCollection.FindOne(context.TODO(), bson.M{"_id": tiketID}).Decode(&tiket)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tiket not found"})
	}

	// Verify user exists
	userCollection := config.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	// Create pembayaran with ObjectID
	now := time.Now()
	pembayaran := models.Pembayaran{
		ID:                primitive.NewObjectID(),
		TiketID:           tiketID,
		UserID:            userID,
		Jumlah:            input.Jumlah,
		MetodePembayaran:  input.MetodePembayaran,
		Status:            input.Status,
		TanggalPembayaran: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result, err := collection.InsertOne(context.TODO(), pembayaran)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the inserted ID
	pembayaran.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message":    "Pembayaran created successfully",
		"pembayaran": pembayaran,
	})
	return c.Status(201).JSON(fiber.Map{
		"message":    "Pembayaran created successfully",
		"pembayaran": pembayaran,
	})
}

func UpdatePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")
	var pembayaran models.Pembayaran
	if err := c.BodyParser(&pembayaran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	update := bson.M{
		"$set": bson.M{
			"tiket_id":           pembayaran.TiketID,
			"user_id":            pembayaran.UserID,
			"jumlah":             pembayaran.Jumlah,
			"metode_pembayaran":  pembayaran.MetodePembayaran,
			"status":             pembayaran.Status,
			"tanggal_pembayaran": pembayaran.TanggalPembayaran,
			"updated_at":         time.Now(),
		},
	}

	res, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Pembayaran updated"})
}

func DeletePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")
	res, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found"})
	}
	return c.JSON(fiber.Map{"message": "Pembayaran deleted"})
}
