package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/config/middleware"
	"go-get-backend/models"
	"go-get-backend/pkg/password"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Register godoc
// @Summary Register new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/register [post]
func Register(c *fiber.Ctx) error {
	userCollection := config.DB.Collection("users")
	var input models.UserRegister

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validation
	if input.ID == "" || input.Email == "" || input.Nama == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Set default role if not provided
	if input.Role == "" {
		input.Role = "user"
	}

	// Check if user already exists (by email or ID)
	count, _ := userCollection.CountDocuments(context.TODO(), bson.M{
		"$or": []bson.M{
			{"email": input.Email},
			{"id": input.ID},
		},
	})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "User with this email or ID already exists",
		})
	}

	// Hash password
	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := models.User{
		ID:       input.ID,
		Nama:     input.Nama,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     input.Role,
	}

	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func Login(c *fiber.Ctx) error {
	userCollection := config.DB.Collection("users")
	var input models.UserLogin

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validation
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

	// Verify password
	err = password.VerifyPassword(input.Password, user.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token
	jwtPayload := models.JWTPayload{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	token, err := middleware.Encoder(jwtPayload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile from JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/profile [get]
func GetProfile(c *fiber.Ctx) error {
	// Get user from middleware context
	userData := c.Locals("user").(*models.JWTPayload)

	userCollection := config.DB.Collection("users")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id": userData.ID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.JSON(fiber.Map{
		"user": user,
	})
}
