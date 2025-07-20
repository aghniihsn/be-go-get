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

// GetAllFilms godoc
// @Summary Get all films
// @Description Get list of all films in the cinema
// @Tags Films
// @Accept json
// @Produce json
// @Success 200 {array} models.Film
// @Failure 500 {object} map[string]interface{}
// @Router /api/films [get]
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

// GetFilmByID godoc
// @Summary Get film by ID
// @Description Get a single film by its ID
// @Tags Films
// @Accept json
// @Produce json
// @Param id path string true "Film ID"
// @Success 200 {object} models.Film
// @Failure 404 {object} map[string]interface{}
// @Router /api/films/{id} [get]
func GetFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var film models.Film
	err := config.DB.Collection("films").FindOne(context.TODO(), bson.M{"id": id}).Decode(&film)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}
	return c.JSON(film)
}

// CreateFilm godoc
// @Summary Create a new film
// @Description Create a new film in the cinema
// @Tags Films
// @Accept json
// @Produce json
// @Param film body models.Film true "Film object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/films [post]
func CreateFilm(c *fiber.Ctx) error {
	var input struct {
		Title       string `json:"title"`
		Genre       string `json:"genre"`
		Duration    int    `json:"duration"`
		Rating      string `json:"rating"`
		Description string `json:"description"`
		PosterURL   string `json:"poster_url"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.Title == "" || input.Genre == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title and genre are required"})
	}
	if input.Duration <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be positive"})
	}

	// Create film with ObjectID
	now := time.Now()
	film := models.Film{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Genre:       input.Genre,
		Duration:    input.Duration,
		Rating:      input.Rating,
		Description: input.Description,
		PosterURL:   input.PosterURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	collection := config.DB.Collection("films")

	result, err := collection.InsertOne(context.TODO(), film)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the inserted ID
	film.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Film created successfully",
		"film":    film,
	})
}

// UpdateFilm godoc
// @Summary Update a film
// @Description Update film by ID
// @Tags Films
// @Accept json
// @Produce json
// @Param id path string true "Film ID"
// @Param film body models.Film true "Film object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/films/{id} [put]
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
			"title":       film.Title,
			"genre":       film.Genre,
			"duration":    film.Duration,
			"rating":      film.Rating,
			"description": film.Description,
			"poster_url":  film.PosterURL,
			"updated_at":  time.Now(),
		},
	}

	collection := config.DB.Collection("films")
	res, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found or update failed"})
	}
	return c.JSON(fiber.Map{"message": "Film updated successfully"})
}

// DeleteFilm godoc
// @Summary Delete a film
// @Description Delete film by ID
// @Tags Films
// @Accept json
// @Produce json
// @Param id path string true "Film ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/films/{id} [delete]
func DeleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("films")
	res, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}
	return c.JSON(fiber.Map{"message": "Film deleted"})
}
