package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/config/middleware"
	"go-get-backend/models"
	"go-get-backend/pkg/password"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user (user or admin)
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/register [post]
// Register handles user registration
func Register(c *fiber.Ctx) error {
	userCollection := config.DB.Collection("users")
	var input models.UserRegister

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	// Gender validation (only 'male' or 'female' allowed)
	if input.Gender != "" && input.Gender != "male" && input.Gender != "female" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Gender must be either 'male' or 'female'",
		})
	}

	// Validation
	if input.Email == "" || input.Username == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email, username, and password are required",
		})
	}

	// Set default role if not provided
	if input.Role == "" {
		input.Role = "user"
	}

	// Check if user already exists (by email or username)
	count, _ := userCollection.CountDocuments(context.TODO(), bson.M{
		"$or": []bson.M{
			{"email": input.Email},
			{"username": input.Username},
		},
	})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "User with this email or username already exists",
		})
	}

	// Hash password
	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user with ObjectID
	now := time.Now()
	user := models.User{
		ID:                primitive.NewObjectID(),
		Username:          input.Username,
		Email:             input.Email,
		Password:          hashedPassword,
		Role:              input.Role,
		Firstname:         input.Firstname,
		Lastname:          input.Lastname,
		Gender:            input.Gender,
		PhoneNumber:       input.PhoneNumber,
		ProfilePictureURL: input.ProfilePictureURL,
		Address:           input.Address,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set the inserted ID
	user.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"_id":                 user.ID,
			"username":            user.Username,
			"email":               user.Email,
			"role":                user.Role,
			"firstname":           user.Firstname,
			"lastname":            user.Lastname,
			"gender":              user.Gender,
			"phone_number":        user.PhoneNumber,
			"profile_picture_url": user.ProfilePictureURL,
			"address":             user.Address,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/login [post]
// Login handles user authentication
func Login(c *fiber.Ctx) error {
	userCollection := config.DB.Collection("users")
	var input models.UserLogin

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Find user by email
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check password using VerifyPassword
	if password.VerifyPassword(input.Password, user.Password) != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token using Encoder
	token, err := middleware.Encoder(models.JWTPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"_id":      user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// GetProfile handles getting user profile
func GetProfile(c *fiber.Ctx) error {
	currentUser := c.Locals("user").(*models.JWTPayload)

	userCollection := config.DB.Collection("users")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": currentUser.ID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"_id":                 user.ID,
		"username":            user.Username,
		"email":               user.Email,
		"role":                user.Role,
		"firstname":           user.Firstname,
		"lastname":            user.Lastname,
		"gender":              user.Gender,
		"phone_number":        user.PhoneNumber,
		"profile_picture_url": user.ProfilePictureURL,
		"address":             user.Address,
		"created_at":          user.CreatedAt,
		"updated_at":          user.UpdatedAt,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	currentUser := c.Locals("user").(*models.JWTPayload)
	userCollection := config.DB.Collection("users")

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	// Gender validation (only 'male' or 'female' allowed)
	if input.Gender != "" && input.Gender != "male" && input.Gender != "female" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Gender must be either 'male' or 'female'",
		})
	}

	// Build update document
	updateDoc := bson.M{
		"updated_at": time.Now(),
	}
	if input.Username != "" {
		updateDoc["username"] = input.Username
	}
	if input.Email != "" {
		updateDoc["email"] = input.Email
	}
	if input.Firstname != "" {
		updateDoc["firstname"] = input.Firstname
	}
	if input.Lastname != "" {
		updateDoc["lastname"] = input.Lastname
	}
	if input.Gender != "" {
		updateDoc["gender"] = input.Gender
	}
	if input.PhoneNumber != "" {
		updateDoc["phone_number"] = input.PhoneNumber
	}
	if input.ProfilePictureURL != "" {
		updateDoc["profile_picture_url"] = input.ProfilePictureURL
	}
	if input.Address != "" {
		updateDoc["address"] = input.Address
	}

	// Update user
	_, err := userCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": currentUser.ID},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}
