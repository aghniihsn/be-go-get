package controllers

import (
	"context"
	"fmt"
	"go-get-backend/config"
	"go-get-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllPembayaran godoc
// @Summary Get all payments (admin only)
// @Description Get a list of all payment records
// @Tags Pembayaran
// @Produce json
// @Success 200 {array} models.Pembayaran
// @Failure 500 {object} map[string]interface{}
// @Router /api/pembayarans [get]
// @Security BearerAuth
func GetAllPembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	var data []models.Pembayaran
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cursor.All(context.TODO(), &data); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// GetPembayaranByID godoc
// @Summary Get payment details by ID
// @Description Get detailed information about a specific payment
// @Tags Pembayaran
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.Pembayaran
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/pembayarans/{id} [get]
// @Security BearerAuth
func GetPembayaranByID(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}

	var pembayaran models.Pembayaran
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&pembayaran)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found"})
	}
	return c.JSON(pembayaran)
}

// GetPembayaranByUserID godoc
// @Summary Get all payments for a specific user
// @Description Retrieve payment history for a user
// @Tags Pembayaran
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Pembayaran
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/pembayarans/user/{user_id} [get]
// @Security BearerAuth
func GetPembayaranByUserID(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	userID := c.Params("user_id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid User ObjectID format"})
	}

	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user payments"})
	}
	defer cursor.Close(context.TODO())

	var pembayarans []models.Pembayaran
	if err = cursor.All(context.TODO(), &pembayarans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode user payments"})
	}

	return c.JSON(pembayarans)
}

// CreatePembayaran godoc
// @Summary Create a new payment
// @Description Create a new payment for a ticket
// @Tags Pembayaran
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Payment details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/pembayarans [post]
// @Security BearerAuth
func CreatePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")

	var input models.PaymentRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.TiketID == "" || input.UserID == "" || input.MetodePembayaran == "" || input.Jumlah <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "TiketID, UserID, MetodePembayaran, and Jumlah are required and Jumlah > 0"})
	}

	// Convert string IDs to ObjectIDs
	tiketID, err := primitive.ObjectIDFromHex(input.TiketID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid TiketID format"})
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UserID format"})
	}

	// Verify tiket exists and not already paid
	tiketCollection := config.DB.Collection("tikets")
	var tiket models.Tiket
	err = tiketCollection.FindOne(context.TODO(), bson.M{"_id": tiketID}).Decode(&tiket)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tiket not found"})
	}

	// Check if payment already exists for this ticket
	paymentExists, err := collection.CountDocuments(context.TODO(), bson.M{"tiket_id": tiketID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check existing payments"})
	}
	if paymentExists > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Payment already exists for this ticket"})
	}

	// Verify user exists and matches the ticket
	userCollection := config.DB.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	// Verify user owns the ticket
	if tiket.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "User does not own this ticket"})
	}

	// Create pembayaran with ObjectID
	now := time.Now()
	pembayaran := models.Pembayaran{
		ID:                primitive.NewObjectID(),
		TiketID:           tiketID,
		Jumlah:            input.Jumlah,
		MetodePembayaran:  input.MetodePembayaran,
		Status:            models.PembayaranStatusPending, // Initial status is pending
		TanggalPembayaran: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result, err := collection.InsertOne(context.TODO(), pembayaran)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the inserted ID
	pembayaran.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message":    "Pembayaran created successfully",
		"pembayaran": pembayaran,
	})
	return c.Status(201).JSON(fiber.Map{
		"message":    "Pembayaran created successfully",
		"pembayaran": pembayaran,
	})
}

// UpdatePembayaran godoc
// @Summary Update payment status
// @Description Update payment information including status and payment proof
// @Tags Pembayaran
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param payment body object{status=string,bukti_pembayaran=string} true "Payment update information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/pembayarans/{id} [put]
// @Security BearerAuth
func UpdatePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")

	// Parse ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment ID format"})
	} // Get existing payment
	var existingPayment models.Pembayaran
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingPayment)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}

	// Parse update data
	var input struct {
		Status          string `json:"status"`
		BuktiPembayaran string `json:"bukti_pembayaran,omitempty"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate status
	if input.Status != "" {
		validStatus := map[string]bool{
			models.PembayaranStatusPending:   true,
			models.PembayaranStatusCompleted: true,
			models.PembayaranStatusFailed:    true,
		}

		if !validStatus[input.Status] {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid payment status"})
		}
	} else {
		input.Status = existingPayment.Status
	}

	// Prepare update document
	update := bson.M{
		"$set": bson.M{
			"status":     input.Status,
			"updated_at": time.Now(),
		},
	}

	// Add bukti_pembayaran if provided
	if input.BuktiPembayaran != "" {
		update["$set"].(bson.M)["bukti_pembayaran"] = input.BuktiPembayaran
	}

	// Update payment
	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil || res.MatchedCount == 0 {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update payment"})
	}

	// If payment status is completed, update ticket status
	if input.Status == models.PembayaranStatusCompleted {
		tiketCollection := config.DB.Collection("tikets")
		_, err := tiketCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": existingPayment.TiketID},
			bson.M{"$set": bson.M{"status": models.TiketStatusConfirmed, "updated_at": time.Now()}},
		)

		if err != nil {
			// Log error but don't fail the response
			// In production, you might want to use a proper logger
			fmt.Println("Failed to update ticket status:", err)
		}
	}

	return c.JSON(fiber.Map{"message": "Payment updated successfully"})
}

// DeletePembayaran godoc
// @Summary Delete a payment (admin only)
// @Description Permanently delete a payment record
// @Tags Pembayaran
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/pembayarans/{id} [delete]
// @Security BearerAuth
func DeletePembayaran(c *fiber.Ctx) error {
	collection := config.DB.Collection("pembayarans")
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment ID format"})
	}

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pembayaran not found"})
	}
	return c.JSON(fiber.Map{"message": "Pembayaran deleted"})
}
