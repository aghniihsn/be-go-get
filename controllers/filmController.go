// Helper function to validate film rating}

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-get-backend/config"
	"go-get-backend/models"
	"go-get-backend/pkg/storage"

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

// Helper function to validate film rating
func isValidRating(rating string) bool {
	validRatings := []string{models.RatingSemua, models.RatingAnak, models.RatingRemaja, models.RatingDewasa}
	for _, r := range validRatings {
		if r == rating {
			return true
		}
	}
	return false
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
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Expected multipart form"})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	durationStr := c.FormValue("duration")
	ratingStr := c.FormValue("rating")
	genresStr := c.FormValue("genre")

	if title == "" || durationStr == "" || ratingStr == "" || genresStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title, duration, rating, and genre are required"})
	}

	// Konversi duration ke int
	var duration int
	if _, err := fmt.Sscanf(durationStr, "%d", &duration); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Duration must be a number"})
	}

	// Validasi rating
	if !isValidRating(ratingStr) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rating. Must be one of: Semua Umur, Anak-anak, Remaja, Dewasa"})
	}

	// Parse genres
	genres := strings.Split(genresStr, ",")
	for i, genre := range genres {
		genres[i] = strings.TrimSpace(genre)
	}

	// Handle poster upload
	posterFiles := form.File["poster"]
	if len(posterFiles) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Poster file is required"})
	}

	// Get storage service
	storageService, err := storage.GetStorageService()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Storage initialization failed: %v", err)})
	}

	// Upload poster to storage
	posterURL, err := storageService.Upload(posterFiles[0], storage.FilmPoster, posterFiles[0].Filename)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to upload poster: %v", err)})
	}

	now := time.Now()
	film := models.Film{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Genre:       genres,
		Duration:    duration,
		Rating:      ratingStr,
		Description: description,
		PosterURL:   posterURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	collection := config.DB.Collection("films")
	result, err := collection.InsertOne(context.TODO(), film)
	if err != nil {
		_ = storageService.Delete(posterURL)
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

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}

	// Get existing film
	collection := config.DB.Collection("films")
	var existingFilm models.Film
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingFilm)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Expected multipart form"})
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	if title := c.FormValue("title"); title != "" {
		update["title"] = title
	}
	if description := c.FormValue("description"); description != "" {
		update["description"] = description
	}
	if durationStr := c.FormValue("duration"); durationStr != "" {
		var duration int
		if _, err := fmt.Sscanf(durationStr, "%d", &duration); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Duration must be a number"})
		}
		update["duration"] = duration
	}
	if ratingStr := c.FormValue("rating"); ratingStr != "" {
		if !isValidRating(ratingStr) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid rating. Must be one of: Semua Umur, Anak-anak, Remaja, Dewasa"})
		}
		update["rating"] = ratingStr
	}
	if genresStr := c.FormValue("genre"); genresStr != "" {
		genres := strings.Split(genresStr, ",")
		for i, genre := range genres {
			genres[i] = strings.TrimSpace(genre)
		}
		update["genre"] = genres
	}

	// Handle poster update if provided
	posterFiles := form.File["poster"]
	if len(posterFiles) > 0 {
		storageService, err := storage.GetStorageService()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Storage initialization failed: %v", err)})
		}
		posterURL, err := storageService.Upload(posterFiles[0], storage.FilmPoster, posterFiles[0].Filename)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to upload poster: %v", err)})
		}
		if existingFilm.PosterURL != "" {
			_ = storageService.Delete(existingFilm.PosterURL)
		}
		update["poster_url"] = posterURL
	}

	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}

	var updatedFilm models.Film
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&updatedFilm)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve updated film"})
	}

	return c.JSON(updatedFilm)
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
