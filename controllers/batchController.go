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

// CreateBatchTikets godoc
// @Summary Create multiple tickets at once
// @Tags Tiket
// @Accept json
// @Produce json
// @Param batchRequest body models.CreateBatchTicketsRequest true "Batch tickets request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/tikets/batch [post]
// @Security BearerAuth
func CreateBatchTikets(c *fiber.Ctx) error {
	tiketCollection := config.DB.Collection("tikets")

	// Parse request
	var input models.CreateBatchTicketsRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.JadwalID == "" || input.UserID == "" || len(input.Kursis) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "JadwalID, UserID, and at least one Kursi are required"})
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

	// Check for valid and available seats
	var invalidSeats []string
	var alreadyBookedSeats []string

	// Check validity of all seats
	for _, kursi := range input.Kursis {
		kursiValid := false
		for _, seat := range models.StudioSeats {
			if kursi == seat {
				kursiValid = true
				break
			}
		}
		if !kursiValid {
			invalidSeats = append(invalidSeats, kursi)
		}
	}

	// If any invalid seats, return error
	if len(invalidSeats) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":         "Some seats are invalid",
			"invalid_seats": invalidSeats,
		})
	}

	// Check availability of all seats
	for _, kursi := range input.Kursis {
		count, err := tiketCollection.CountDocuments(context.TODO(), bson.M{
			"jadwal_id": jadwalID,
			"kursi":     kursi,
			"status":    bson.M{"$in": []string{models.TiketStatusConfirmed, models.TiketStatusWaitingPayment}},
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if count > 0 {
			alreadyBookedSeats = append(alreadyBookedSeats, kursi)
		}
	}

	// If any seats already booked, return error
	if len(alreadyBookedSeats) > 0 {
		return c.Status(409).JSON(fiber.Map{
			"error":        "Some seats are already booked",
			"booked_seats": alreadyBookedSeats,
		})
	}

	// Create tickets
	now := time.Now()
	var tikets []interface{}
	var createdTikets []models.Tiket

	for _, kursi := range input.Kursis {
		tiket := models.Tiket{
			ID:               primitive.NewObjectID(),
			UserID:           userID,
			JadwalID:         jadwalID,
			Kursi:            kursi,
			Status:           models.TiketStatusWaitingPayment, // Default status untuk tiket baru
			TanggalPembelian: now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		tikets = append(tikets, tiket)
		createdTikets = append(createdTikets, tiket)
	}

	// Insert all tickets
	_, err = tiketCollection.InsertMany(context.TODO(), tikets)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tickets: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Tickets created successfully",
		"tickets": createdTikets,
		"count":   len(createdTikets),
	})
}

// GetTicketSummary godoc
// @Summary Get detailed ticket information with film and jadwal details
// @Tags Tiket
// @Produce json
// @Param id path string true "Tiket ID"
// @Success 200 {object} models.TicketSummary
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/tikets/{id}/summary [get]
// @Security BearerAuth
func GetTicketSummary(c *fiber.Ctx) error {
	tiketID := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(tiketID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Tiket ID"})
	}

	// Get ticket
	tiketCollection := config.DB.Collection("tikets")
	var tiket models.Tiket
	err = tiketCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&tiket)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tiket not found"})
	}

	// Get jadwal
	jadwalCollection := config.DB.Collection("jadwals")
	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": tiket.JadwalID}).Decode(&jadwal)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}

	// Get film
	filmCollection := config.DB.Collection("films")
	var film models.Film
	err = filmCollection.FindOne(context.TODO(), bson.M{"_id": jadwal.FilmID}).Decode(&film)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Film not found"})
	}

	// Get payment if exists
	pembayaranCollection := config.DB.Collection("pembayarans")
	var pembayaran models.Pembayaran
	statusPembayaran := "belum_dibayar"
	var tanggalPembayaran time.Time

	err = pembayaranCollection.FindOne(context.TODO(), bson.M{"tiket_id": objectID}).Decode(&pembayaran)
	if err == nil {
		statusPembayaran = pembayaran.Status
		tanggalPembayaran = pembayaran.TanggalPembayaran
	}

	// Create ticket summary
	summary := models.TicketSummary{
		Tiket:             tiket,
		FilmTitle:         film.Title,
		FilmPoster:        film.PosterURL,
		JadwalWaktu:       jadwal.Waktu,
		JadwalTanggal:     jadwal.Tanggal,
		JadwalRuangan:     jadwal.Ruangan,
		HargaTiket:        jadwal.Harga,
		StatusPembayaran:  statusPembayaran,
		TanggalPembayaran: tanggalPembayaran,
	}

	return c.JSON(summary)
}

// GetPaymentMethods godoc
// @Summary Get available payment methods
// @Tags Pembayaran
// @Produce json
// @Success 200 {array} models.MetodePembayaran
// @Router /api/payment-methods [get]
func GetPaymentMethods(c *fiber.Ctx) error {
	return c.JSON(models.AvailablePaymentMethods)
}
