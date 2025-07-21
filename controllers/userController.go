package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"
	"go-get-backend/pkg/password"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser is now deprecated, use Register in authController instead
func CreateUser(c *fiber.Ctx) error {
	return c.Status(400).JSON(fiber.Map{
		"error": "Please use /api/auth/register endpoint for user registration",
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by ID (admin only or own profile)
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	requestedUserIDStr := c.Params("id")
	currentUser := c.Locals("user").(*models.JWTPayload)

	// Convert string to ObjectID for comparison
	requestedUserID, err := primitive.ObjectIDFromHex(requestedUserIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	// Allow users to see their own profile, or admin to see any profile
	if currentUser.Role != "admin" && currentUser.ID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You can only access your own profile",
		})
	}

	userCollection := config.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": requestedUserID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Remove password from response
	user.Password = ""
	return c.JSON(user)
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Update user information (own profile or admin)
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	currentUser := c.Locals("user").(*models.JWTPayload)

	// Convert string to ObjectID for comparison
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	// Allow users to update their own profile, or admin to update any profile
	if currentUser.Role != "admin" && currentUser.ID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You can only update your own profile",
		})
	}

	userCollection := config.DB.Collection("users")

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Prepare update data
	updateData := bson.M{}
	if input.Username != "" {
		updateData["username"] = input.Username
	}
	if input.Email != "" {
		updateData["email"] = input.Email
	}
	if input.Password != "" {
		// Hash new password
		hashedPassword, err := password.HashPassword(input.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		updateData["password"] = hashedPassword
	}
	// Only admin can change roles
	if input.Role != "" && currentUser.Role == "admin" {
		updateData["role"] = input.Role
	}
	if input.Firstname != "" {
		updateData["firstname"] = input.Firstname
	}
	if input.Lastname != "" {
		updateData["lastname"] = input.Lastname
	}
	if input.Gender != "" {
		updateData["gender"] = input.Gender
	}
	if input.PhoneNumber != "" {
		updateData["phone_number"] = input.PhoneNumber
	}
	if input.ProfilePictureURL != "" {
		updateData["profile_picture_url"] = input.ProfilePictureURL
	}
	if input.Address != "" {
		updateData["address"] = input.Address
	}

	// Add updated timestamp
	updateData["updated_at"] = time.Now()

	update := bson.M{"$set": updateData}

	res, err := userCollection.UpdateOne(context.TODO(), bson.M{"_id": userID}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Update failed"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}
