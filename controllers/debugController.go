package controllers

import (
	"context"
	"go-get-backend/config"
	"go-get-backend/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DebugAllIDs godoc
// @Summary Get all document IDs for debugging
// @Tags Debug
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/debug/ids [get]
func DebugAllIDs(c *fiber.Ctx) error {
	// Koleksi yang ingin kita periksa
	collections := []string{"jadwals", "films", "users", "tikets", "pembayarans"}
	results := make(map[string]interface{})

	for _, collName := range collections {
		coll := config.DB.Collection(collName)
		cursor, err := coll.Find(context.TODO(), bson.M{})
		if err != nil {
			results[collName+"_error"] = err.Error()
			continue
		}

		var documents []bson.M
		if err = cursor.All(context.TODO(), &documents); err != nil {
			results[collName+"_error"] = err.Error()
			continue
		}

		// Ekstrak hanya ID untuk menjaga respons tetap singkat
		ids := make([]map[string]interface{}, 0)
		for _, doc := range documents {
			id := doc["_id"]
			if oid, ok := id.(primitive.ObjectID); ok {
				item := map[string]interface{}{
					"id":  oid.Hex(),
					"raw": oid,
				}
				ids = append(ids, item)
			}
		}

		results[collName] = ids
	}

	return c.JSON(fiber.Map{
		"message": "Debug IDs",
		"data":    results,
	})
}

// DebugCheckJadwal godoc
// @Summary Check if jadwal exists with specific ID
// @Tags Debug
// @Produce json
// @Param id path string true "Jadwal ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/debug/jadwal/{id} [get]
func DebugCheckJadwal(c *fiber.Ctx) error {
	jadwalCollection := config.DB.Collection("jadwals")
	id := c.Params("id")

	// Hapus spasi jika ada
	id = strings.TrimSpace(id)

	// Coba konversi string ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid ID format",
			"id":      id,
			"message": err.Error(),
		})
	}

	// Cek apakah jadwal ada dengan _id yang diberikan
	var jadwal models.Jadwal
	err = jadwalCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&jadwal)
	if err != nil {
		// Kembalikan hasil pencarian raw query untuk debugging
		var allJadwals []bson.M
		cursor, _ := jadwalCollection.Find(context.TODO(), bson.M{})
		_ = cursor.All(context.TODO(), &allJadwals)

		// Hanya ambil beberapa jadwal untuk display
		jadwalIDs := []string{}
		for _, j := range allJadwals {
			if oid, ok := j["_id"].(primitive.ObjectID); ok {
				jadwalIDs = append(jadwalIDs, oid.Hex())
			}
		}

		return c.Status(404).JSON(fiber.Map{
			"error":                "Jadwal not found",
			"searched_id":          id,
			"searched_object_id":   objectID.Hex(),
			"available_jadwal_ids": jadwalIDs,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Jadwal found",
		"id":      id,
		"jadwal":  jadwal,
	})
}
