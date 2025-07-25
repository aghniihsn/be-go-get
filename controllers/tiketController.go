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

// GetAllTikets godoc
// @Summary Get all tickets (admin only)
// @Tags Tiket
// @Produce json
// @Success 200 {array} models.Tiket
// @Failure 403 {object} map[string]interface{}
// @Router /api/tikets [get]
// @Security BearerAuth
func GetAllTikets(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tikets)
}

// GetTiketByID godoc
// @Summary Get ticket by ID
// @Tags Tiket
// @Produce json
// @Param id path string true "Tiket ID"
// @Success 200 {object} models.Tiket
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/tikets/{id} [get]
// @Security BearerAuth
func GetTiketByID(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Tiket ID"})
	}
	var tiket models.Tiket
	err = tiketCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&tiket)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}
	return c.JSON(tiket)
}

// GetTiketByUserID godoc
// @Summary Get tickets by user ID (admin or self)
// @Tags Tiket
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Tiket
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/tikets/user/{user_id} [get]
// @Security BearerAuth
func GetTiketByUserID(c *fiber.Ctx) error {
	requestedUserIDStr := c.Params("user_id")
	currentUser := c.Locals("user").(*models.JWTPayload)

	// Convert string to ObjectID for comparison
	requestedUserID, err := primitive.ObjectIDFromHex(requestedUserIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	// Allow users to see only their own tickets, admin can see all
	if currentUser.Role != "admin" && currentUser.ID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{"user_id": requestedUserID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tikets)
}

// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/tikets [post]
// @Security BearerAuth
func CreateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")

	var input struct {
		JadwalID string `json:"jadwal_id"`
		UserID   string `json:"user_id"`
		Kursi    string `json:"kursi"`
		Status   string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.JadwalID == "" || input.UserID == "" || input.Kursi == "" {
		return c.Status(400).JSON(fiber.Map{"error": "JadwalID, UserID, and Kursi are required"})
	}

	// Convert string IDs to ObjectIDs
	jadwalID, err := primitive.ObjectIDFromHex(input.JadwalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JadwalID format"})
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UserID format"})
	}

	// Verify jadwal exists
	jadwalCollection := config.DB.Collection("jadwals")
	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": jadwalID}).Decode(&jadwal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Jadwal not found"})
	}

	// Verify user exists
	userCollection := config.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	// Validasi kursi: harus ada di StudioSeats
	kursiValid := false
	for _, seat := range models.StudioSeats {
		if input.Kursi == seat {
			kursiValid = true
			break
		}
	}
	if !kursiValid {
		return c.Status(400).JSON(fiber.Map{"error": "Kursi tidak valid"})
	}

	// Cek kursi sudah dibooking di jadwal yang sama (status confirmed atau waiting_for_payment)
	count, err := tiketCollection.CountDocuments(context.TODO(), bson.M{
		"jadwal_id": jadwalID,
		"kursi":     input.Kursi,
		"status":    bson.M{"$in": []string{models.TiketStatusConfirmed, models.TiketStatusWaitingPayment}},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Kursi sudah dibooking"})
	}

	// Validasi status
	status := input.Status
	validStatus := map[string]bool{
		models.TiketStatusWaitingPayment: true,
		models.TiketStatusConfirmed:      true,
		models.TiketStatusCancelled:      true,
		models.TiketStatusUsed:           true,
	}
	if status == "" {
		status = models.TiketStatusWaitingPayment // Default status untuk tiket baru
	}
	if !validStatus[status] {
		return c.Status(400).JSON(fiber.Map{"error": "Status tiket tidak valid"})
	}

	// Create tiket with ObjectID
	now := time.Now()
	tiket := models.Tiket{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		JadwalID:         jadwalID,
		Kursi:            input.Kursi,
		Status:           status,
		TanggalPembelian: now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	result, err := tiketCollection.InsertOne(context.TODO(), tiket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the inserted ID
	tiket.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Tiket created successfully",
		"tiket":   tiket,
	})
}

// UpdateTiket godoc
// @Summary Update ticket (seat/status)
// @Tags Tiket
// @Accept json
// @Produce json
// @Param id path string true "Tiket ID"
// @Param tiket body models.Tiket true "Tiket object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/tikets/{id} [put]
// @Security BearerAuth
func UpdateTiket(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Tiket ID"})
	}

	var input struct {
		Kursi  string `json:"kursi"`
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi status
	validStatus := map[string]bool{
		models.TiketStatusConfirmed: true,
		models.TiketStatusCancelled: true,
		models.TiketStatusUsed:      true,
	}
	if input.Status != "" && !validStatus[input.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Status tiket tidak valid"})
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if input.Kursi != "" {
		update["$set"].(bson.M)["kursi"] = input.Kursi
	}
	if input.Status != "" {
		update["$set"].(bson.M)["status"] = input.Status
	}

	res, err := tiketCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}

	return c.JSON(fiber.Map{"message": "Tiket updated successfully"})
}

// GetKursiKosong godoc
// @Summary Get available seats for a schedule
// @Tags Tiket
// @Produce json
// @Param jadwalId path string true "Jadwal ID"
// @Success 200 {array} string
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/tikets/kursi-kosong/{jadwalId} [get]
func GetKursiKosong(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwal_id")
	objectID, err := primitive.ObjectIDFromHex(jadwalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Jadwal ID"})
	}
	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{
		"jadwal_id": objectID,
		"status":    bson.M{"$in": []string{models.TiketStatusConfirmed, models.TiketStatusWaitingPayment}},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	booked := make(map[string]bool)
	for _, t := range tikets {
		booked[t.Kursi] = true
	}
	var available []string
	for _, seat := range models.StudioSeats {
		if !booked[seat] {
			available = append(available, seat)
		}
	}
	return c.JSON(available)
}

// GetTiketUser godoc
// @Summary Get current user's ticket history
// @Tags Tiket
// @Produce json
// @Success 200 {array} models.Tiket
// @Failure 401 {object} map[string]interface{}
// @Router /api/tikets/me [get]
// @Security BearerAuth
func GetTiketUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.JWTPayload)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	tiketCollection := config.DB.Collection("tikets")
	var tikets []models.Tiket
	cursor, err := tiketCollection.Find(context.TODO(), bson.M{"user_id": user.ID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &tikets); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tikets)
}

// CancelTiket godoc
// @Summary Cancel a ticket
// @Tags Tiket
// @Param id path string true "Tiket ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/tikets/batal/{id} [put]
func CancelTiket(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Tiket ID"})
	}
	tiketCollection := config.DB.Collection("tikets")
	update := bson.M{"$set": bson.M{"status": models.TiketStatusCancelled, "updated_at": time.Now()}}
	res, err := tiketCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found or cancel failed"})
	}
	return c.JSON(fiber.Map{"message": "Tiket berhasil dibatalkan"})
}
