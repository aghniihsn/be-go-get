package controllers

import (
	"context"
	"fmt"
	"go-get-backend/config"
	"go-get-backend/models"
	"go-get-backend/pkg/storage"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateProfileWithImage godoc
// @Summary Update user profile with image
// @Description Update user profile information with optional profile picture upload
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param firstname formData string false "First name"
// @Param lastname formData string false "Last name"
// @Param gender formData string false "Gender (male or female)"
// @Param phone_number formData string false "Phone number"
// @Param address formData string false "Address"
// @Param profile_picture formData file false "Profile picture"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/profile [put]
// @Security BearerAuth
func UpdateProfileWithImage(c *fiber.Ctx) error {
	// Get user from context
	userData := c.Locals("user").(*models.JWTPayload)
	userID := userData.ID // JWTPayload already has ObjectID, no need to convert

	// Get user from database
	collection := config.DB.Collection("users")
	var user models.User
	if err := collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Expected multipart form"})
	}

	// Prepare update document
	update := bson.M{
		"updatedat": time.Now(),
	}

	// Update fields if provided
	if firstname := c.FormValue("firstname"); firstname != "" {
		update["firstname"] = firstname
	}

	if lastname := c.FormValue("lastname"); lastname != "" {
		update["lastname"] = lastname
	}

	if gender := c.FormValue("gender"); gender != "" {
		if gender != "male" && gender != "female" {
			return c.Status(400).JSON(fiber.Map{"error": "Gender must be 'male' or 'female'"})
		}
		update["gender"] = gender
	}

	if phoneNumber := c.FormValue("phone_number"); phoneNumber != "" {
		update["phone_number"] = phoneNumber
	}

	if address := c.FormValue("address"); address != "" {
		update["address"] = address
	}

	// Handle profile picture upload if provided
	profilePicFiles := form.File["profile_picture"]
	if len(profilePicFiles) > 0 {
		// Get storage service
		storageService, err := storage.GetStorageService()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Storage initialization failed: %v", err),
			})
		}

		// Upload profile picture
		profilePicURL, err := storageService.Upload(
			profilePicFiles[0],
			storage.ProfilePicture,
			profilePicFiles[0].Filename,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to upload profile picture: %v", err),
			})
		}

		// Delete old profile picture if exists
		if user.ProfilePictureURL != "" {
			_ = storageService.Delete(user.ProfilePictureURL)
		}

		update["profile_picture_url"] = profilePicURL
	}

	// Update user in database
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Get updated user
	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve updated user"})
	}

	// Don't return password in response
	updatedUser.Password = ""

	return c.JSON(updatedUser)
}
