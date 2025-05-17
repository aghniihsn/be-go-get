package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
	var pembayaran models.Pembayaran
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&pembayaran)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found"})
	}
	return c.JSON(pembayaran)
}

func CreatePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	var pembayaran models.Pembayaran
	if err := c.BodyParser(&pembayaran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if pembayaran.ID == "" || pembayaran.TiketID == "" || pembayaran.Metode == "" || pembayaran.Status == "" || pembayaran.Total <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields required and total > 0"})
	}

	count, _ := collection.CountDocuments(context.TODO(), bson.M{"id": pembayaran.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID already exists"})
	}

	_, err := collection.InsertOne(context.TODO(), pembayaran)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Pembayaran created"})
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
			"tiket_id": pembayaran.TiketID,
			"metode":   pembayaran.Metode,
			"status":   pembayaran.Status,
			"total":    pembayaran.Total,
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
