package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllFilms(c *fiber.Ctx) error {
	collection := config.DB.Collection("films")
	var films []models.Film
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err = cursor.All(context.TODO(), &films); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(films)
}

func GetFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var film models.Film
	err := config.DB.Collection("films").FindOne(context.TODO(), bson.M{"id": id}).Decode(&film)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}
	return c.JSON(film)
}

func CreateFilm(c *fiber.Ctx) error {
	var film models.Film
	if err := c.BodyParser(&film); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if film.ID == "" || film.Title == "" || film.Genre == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}
	if film.Duration <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be positive"})
	}

	collection := config.DB.Collection("films")
	count, _ := collection.CountDocuments(context.TODO(), bson.M{"id": film.ID})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID already exists"})
	}

	_, err := collection.InsertOne(context.TODO(), film)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Film created"})
}

func UpdateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	var film models.Film
	if err := c.BodyParser(&film); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if film.Duration <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be positive"})
	}

	update := bson.M{
		"$set": bson.M{
			"title":    film.Title,
			"genre":    film.Genre,
			"duration": film.Duration,
		},
	}

	collection := config.DB.Collection("films")
	res, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Film updated"})
}

func DeleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("films")
	res, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}
	return c.JSON(fiber.Map{"message": "Film deleted"})
}
