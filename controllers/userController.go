package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(c *fiber.Ctx) error {
	userCollection := config.DB.Collection("users")
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if user.ID == "" || user.Email == "" || user.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields required"})
	}

	count, _ := userCollection.CountDocuments(context.TODO(), bson.M{"id": user.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "User ID already exists"})
	}

	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "User created"})
}

func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	userCollection := config.DB.Collection("users")

	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	userCollection := config.DB.Collection("users")

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	update := bson.M{
		"$set": bson.M{
			"nama":  input.Nama,
			"email": input.Email,
		},
	}

	res, err := userCollection.UpdateOne(context.TODO(), bson.M{"id": userID}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Update failed"})
	}

	return c.JSON(fiber.Map{"message": "User updated"})
}
