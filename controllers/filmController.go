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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}
	var film models.Film
	err = config.DB.Collection("films").FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&film)
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
// @example {"title": "Laskar Pelangi", "genre": ["Drama", "Family"], "duration": 120, "rating": "Remaja", "description": "Film inspiratif tentang anak-anak di Belitung.", "poster_url": "https://example.com/poster.jpg"}
//
//	@Example {
//	  "title": "Laskar Pelangi",
//	  "genre": ["Drama", "Family"],
//	  "duration": 120,
//	  "rating": "Remaja",
//	  "description": "Film inspiratif tentang anak-anak di Belitung.",
//	  "poster_url": "https://example.com/poster.jpg"
//	}
//
// @Success 201 {object} map[string]interface{} "{ 'message': 'Film created successfully', 'film': { ...film fields... } }"
// @Failure 400 {object} map[string]interface{} "{ 'error': 'Invalid input' }"
// @Failure 500 {object} map[string]interface{} "{ 'error': 'Internal server error' }"
// @Security BearerAuth
// @Router /api/films [post]
func CreateFilm(c *fiber.Ctx) error {
	// Define allowed genres and ratings
	allowedGenres := map[string]bool{
		"Action": true, "Adventure": true, "Animation": true, "Biography": true, "Comedy": true, "Crime": true, "Documentary": true, "Drama": true, "Family": true, "Fantasy": true, "History": true, "Horror": true, "Music": true, "Musical": true, "Mystery": true, "Romance": true, "Sci-Fi": true, "Sport": true, "Thriller": true, "War": true, "Western": true,
	}
	allowedRatings := map[string]bool{
		models.RatingSemua: true, models.RatingAnak: true, models.RatingRemaja: true, models.RatingDewasa: true,
	}

	var input struct {
		Title       string   `json:"title"`
		Genre       []string `json:"genre"`
		Duration    int      `json:"duration"`
		Rating      string   `json:"rating"`
		Description string   `json:"description"`
		PosterURL   string   `json:"poster_url"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Title validation
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	// Genre validation: must be array, at least one, all valid
	if len(input.Genre) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one genre is required"})
	}
	for _, g := range input.Genre {
		if !allowedGenres[g] {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid genre: " + g})
		}
	}

	// Duration validation
	if input.Duration <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be positive"})
	}

	// Rating validation
	if !allowedRatings[input.Rating] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rating"})
	}

	// PosterURL validation (simple)
	if input.PosterURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Poster URL is required"})
	}

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
// @example {"title": "Laskar Pelangi", "genre": ["Drama", "Family"], "duration": 120, "rating": "Remaja", "description": "Film inspiratif tentang anak-anak di Belitung.", "poster_url": "https://example.com/poster.jpg"}
//
//	@Example {
//	  "title": "Laskar Pelangi",
//	  "genre": ["Drama", "Family"],
//	  "duration": 120,
//	  "rating": "Remaja",
//	  "description": "Film inspiratif tentang anak-anak di Belitung.",
//	  "poster_url": "https://example.com/poster.jpg"
//	}
//
// @Success 200 {object} map[string]interface{} "{ 'message': 'Film updated successfully' }"
// @Failure 400 {object} map[string]interface{} "{ 'error': 'Invalid input' }"
// @Failure 404 {object} map[string]interface{} "{ 'error': 'Film not found or update failed' }"
// @Security BearerAuth
// @Router /api/films/{id} [put]
func UpdateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	// Define allowed genres and ratings
	allowedGenres := map[string]bool{
		"Action": true, "Adventure": true, "Animation": true, "Biography": true, "Comedy": true, "Crime": true, "Documentary": true, "Drama": true, "Family": true, "Fantasy": true, "History": true, "Horror": true, "Music": true, "Musical": true, "Mystery": true, "Romance": true, "Sci-Fi": true, "Sport": true, "Thriller": true, "War": true, "Western": true,
	}
	allowedRatings := map[string]bool{
		models.RatingSemua: true, models.RatingAnak: true, models.RatingRemaja: true, models.RatingDewasa: true,
	}

	var input struct {
		Title       string   `json:"title"`
		Genre       []string `json:"genre"`
		Duration    int      `json:"duration"`
		Rating      string   `json:"rating"`
		Description string   `json:"description"`
		PosterURL   string   `json:"poster_url"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}
	if len(input.Genre) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one genre is required"})
	}
	for _, g := range input.Genre {
		if !allowedGenres[g] {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid genre: " + g})
		}
	}
	if input.Duration <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be positive"})
	}
	if !allowedRatings[input.Rating] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rating"})
	}
	if input.PosterURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Poster URL is required"})
	}

	update := bson.M{
		"$set": bson.M{
			"title":       input.Title,
			"genre":       input.Genre,
			"duration":    input.Duration,
			"rating":      input.Rating,
			"description": input.Description,
			"poster_url":  input.PosterURL,
			"updated_at":  time.Now(),
		},
	}

	collection := config.DB.Collection("films")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}
	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
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
// @Security BearerAuth
// @Router /api/films/{id} [delete]
func DeleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("films")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}
	return c.JSON(fiber.Map{"message": "Film deleted"})
}
