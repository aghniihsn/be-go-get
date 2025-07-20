package database_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go-get-backend/models"
	"go-get-backend/pkg/password"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseTestSuite adalah test suite untuk testing database integration
type DatabaseTestSuite struct {
	suite.Suite
	db     *mongo.Database
	client *mongo.Client
	ctx    context.Context
}

// SetupSuite dijalankan sekali sebelum semua test
func (suite *DatabaseTestSuite) SetupSuite() {
	// Load environment variables dari .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		// Fallback to default untuk testing lokal jika .env tidak ditemukan
		os.Setenv("MONGOSTRING", "mongodb://localhost:27017")
	}

	// Get MongoDB connection string dari environment
	mongoString := os.Getenv("MONGOSTRING")
	if mongoString == "" {
		log.Fatal("MONGOSTRING environment variable not set")
	}

	// Connect ke database menggunakan connection string dari .env
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString))
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	suite.client = client
	suite.db = client.Database("cinema_booking_test") // Gunakan database terpisah untuk testing
	suite.ctx = context.Background()

	log.Printf("Connected to test database successfully using: %s", mongoString)
}

// TearDownSuite dijalankan sekali setelah semua test
func (suite *DatabaseTestSuite) TearDownSuite() {
	// Cleanup: Clear collections instead of dropping database (untuk MongoDB Atlas)
	collections := []string{"users", "films", "jadwals", "tikets", "pembayarans", "performance_test"}

	for _, collName := range collections {
		collection := suite.db.Collection(collName)
		_, err := collection.DeleteMany(suite.ctx, bson.M{})
		if err != nil {
			log.Printf("Failed to clear collection %s: %v", collName, err)
		}
	}

	// Disconnect
	err := suite.client.Disconnect(suite.ctx)
	if err != nil {
		log.Printf("Failed to disconnect from database: %v", err)
	}

	log.Println("Test database cleaned up successfully")
}

// SetupTest dijalankan sebelum setiap test
func (suite *DatabaseTestSuite) SetupTest() {
	// Clear all collections before each test
	collections := []string{"users", "films", "jadwals", "tikets", "pembayarans"}

	for _, collName := range collections {
		collection := suite.db.Collection(collName)
		_, err := collection.DeleteMany(suite.ctx, bson.M{})
		if err != nil {
			log.Printf("Failed to clear collection %s: %v", collName, err)
		}
	}
}

// TestUserCRUD menguji operasi CRUD untuk User dengan ObjectID
func (suite *DatabaseTestSuite) TestUserCRUD() {
	collection := suite.db.Collection("users")

	// Test Create User
	hashedPassword, err := password.HashPassword("testpassword123")
	suite.NoError(err)

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     "user",
	}

	// Insert user
	result, err := collection.InsertOne(suite.ctx, user)
	suite.NoError(err)
	suite.NotNil(result.InsertedID)

	insertedID := result.InsertedID.(primitive.ObjectID)
	suite.Equal(user.ID, insertedID)

	// Test Read User by ObjectID
	var foundUser models.User
	err = collection.FindOne(suite.ctx, bson.M{"_id": user.ID}).Decode(&foundUser)
	suite.NoError(err)
	suite.Equal(user.Username, foundUser.Username)
	suite.Equal(user.Email, foundUser.Email)
	suite.Equal(user.Role, foundUser.Role)

	// Verify password
	err = password.VerifyPassword("testpassword123", foundUser.Password)
	suite.NoError(err)

	// Test Update User
	updateData := bson.M{
		"$set": bson.M{
			"username": "updateduser",
			"email":    "updated@example.com",
		},
	}

	updateResult, err := collection.UpdateOne(suite.ctx, bson.M{"_id": user.ID}, updateData)
	suite.NoError(err)
	suite.Equal(int64(1), updateResult.ModifiedCount)

	// Verify update
	var updatedUser models.User
	err = collection.FindOne(suite.ctx, bson.M{"_id": user.ID}).Decode(&updatedUser)
	suite.NoError(err)
	suite.Equal("updateduser", updatedUser.Username)
	suite.Equal("updated@example.com", updatedUser.Email)

	// Test Delete User
	deleteResult, err := collection.DeleteOne(suite.ctx, bson.M{"_id": user.ID})
	suite.NoError(err)
	suite.Equal(int64(1), deleteResult.DeletedCount)

	// Verify deletion
	err = collection.FindOne(suite.ctx, bson.M{"_id": user.ID}).Decode(&foundUser)
	suite.Equal(mongo.ErrNoDocuments, err)
}

// TestFilmCRUD menguji operasi CRUD untuk Film dengan ObjectID
func (suite *DatabaseTestSuite) TestFilmCRUD() {
	collection := suite.db.Collection("films")

	// Test Create Film
	film := models.Film{
		ID:          primitive.NewObjectID(),
		Title:       "Test Movie",
		Description: "A test movie description",
		Duration:    120,
		Genre:       "Action",
	}

	// Insert film
	result, err := collection.InsertOne(suite.ctx, film)
	suite.NoError(err)
	suite.NotNil(result.InsertedID)

	// Test Read Film
	var foundFilm models.Film
	err = collection.FindOne(suite.ctx, bson.M{"_id": film.ID}).Decode(&foundFilm)
	suite.NoError(err)
	suite.Equal(film.Title, foundFilm.Title)
	suite.Equal(film.Duration, foundFilm.Duration)
	suite.Equal(film.Genre, foundFilm.Genre)

	// Test Update Film
	updateData := bson.M{
		"$set": bson.M{
			"title":    "Updated Movie",
			"duration": 150,
		},
	}

	updateResult, err := collection.UpdateOne(suite.ctx, bson.M{"_id": film.ID}, updateData)
	suite.NoError(err)
	suite.Equal(int64(1), updateResult.ModifiedCount)

	// Test Delete Film
	deleteResult, err := collection.DeleteOne(suite.ctx, bson.M{"_id": film.ID})
	suite.NoError(err)
	suite.Equal(int64(1), deleteResult.DeletedCount)
}

// TestJadwalCRUD menguji operasi CRUD untuk Jadwal dengan relasi ObjectID
func (suite *DatabaseTestSuite) TestJadwalCRUD() {
	filmCollection := suite.db.Collection("films")
	jadwalCollection := suite.db.Collection("jadwals")

	// Setup: Create a film first
	film := models.Film{
		ID:          primitive.NewObjectID(),
		Title:       "Test Movie for Schedule",
		Description: "Test description",
		Duration:    120,
		Genre:       "Action",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := filmCollection.InsertOne(suite.ctx, film)
	suite.NoError(err)

	// Test Create Jadwal with Film ObjectID reference
	jadwal := models.Jadwal{
		ID:        primitive.NewObjectID(),
		FilmID:    film.ID,
		Tanggal:   "2025-12-25",
		Waktu:     "19:00",
		Ruangan:   "Studio 1",
		Harga:     50000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert jadwal
	result, err := jadwalCollection.InsertOne(suite.ctx, jadwal)
	suite.NoError(err)
	suite.NotNil(result.InsertedID)

	// Test Read Jadwal
	var foundJadwal models.Jadwal
	err = jadwalCollection.FindOne(suite.ctx, bson.M{"_id": jadwal.ID}).Decode(&foundJadwal)
	suite.NoError(err)
	suite.Equal(jadwal.FilmID, foundJadwal.FilmID)
	suite.Equal(jadwal.Ruangan, foundJadwal.Ruangan)
	suite.Equal(jadwal.Harga, foundJadwal.Harga)

	// Test Query by FilmID (ObjectID reference)
	var jadwalsByFilm []models.Jadwal
	cursor, err := jadwalCollection.Find(suite.ctx, bson.M{"film_id": film.ID})
	suite.NoError(err)

	err = cursor.All(suite.ctx, &jadwalsByFilm)
	suite.NoError(err)
	suite.Len(jadwalsByFilm, 1)
	suite.Equal(jadwal.ID, jadwalsByFilm[0].ID)

	// Cleanup
	_, err = jadwalCollection.DeleteOne(suite.ctx, bson.M{"_id": jadwal.ID})
	suite.NoError(err)

	_, err = filmCollection.DeleteOne(suite.ctx, bson.M{"_id": film.ID})
	suite.NoError(err)
}

// TestTiketBookingFlow menguji flow booking tiket dengan multiple ObjectID references
func (suite *DatabaseTestSuite) TestTiketBookingFlow() {
	userCollection := suite.db.Collection("users")
	filmCollection := suite.db.Collection("films")
	jadwalCollection := suite.db.Collection("jadwals")
	tiketCollection := suite.db.Collection("tikets")

	// Setup: Create user, film, and jadwal
	hashedPassword, err := password.HashPassword("password123")
	suite.NoError(err)

	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	film := models.Film{
		ID:          primitive.NewObjectID(),
		Title:       "Booking Test Movie",
		Description: "Test movie for booking",
		Duration:    120,
		Genre:       "Action",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jadwal := models.Jadwal{
		ID:        primitive.NewObjectID(),
		FilmID:    film.ID,
		Tanggal:   "2025-12-25",
		Waktu:     "19:00",
		Ruangan:   "Studio 1",
		Harga:     50000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert all entities
	_, err = userCollection.InsertOne(suite.ctx, user)
	suite.NoError(err)

	_, err = filmCollection.InsertOne(suite.ctx, film)
	suite.NoError(err)

	_, err = jadwalCollection.InsertOne(suite.ctx, jadwal)
	suite.NoError(err)

	// Create tiket with multiple ObjectID references
	tiket := models.Tiket{
		ID:               primitive.NewObjectID(),
		UserID:           user.ID,
		JadwalID:         jadwal.ID,
		Kursi:            "A1",
		Status:           "confirmed",
		TanggalPembelian: time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Insert tiket
	result, err := tiketCollection.InsertOne(suite.ctx, tiket)
	suite.NoError(err)
	suite.NotNil(result.InsertedID)

	// Test aggregation: Get tiket with user and jadwal details
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": tiket.ID},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user_details",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "jadwals",
				"localField":   "jadwal_id",
				"foreignField": "_id",
				"as":           "jadwal_details",
			},
		},
	}

	cursor, err := tiketCollection.Aggregate(suite.ctx, pipeline)
	suite.NoError(err)

	var results []bson.M
	err = cursor.All(suite.ctx, &results)
	suite.NoError(err)
	suite.Len(results, 1)

	result_data := results[0]
	suite.Equal(tiket.Kursi, result_data["kursi"])

	// Verify user details in aggregation
	userDetails := result_data["user_details"].(primitive.A)
	suite.Len(userDetails, 1)

	// Verify jadwal details in aggregation
	jadwalDetails := result_data["jadwal_details"].(primitive.A)
	suite.Len(jadwalDetails, 1)

	// Test query tikets by user
	var userTikets []models.Tiket
	cursor, err = tiketCollection.Find(suite.ctx, bson.M{"user_id": user.ID})
	suite.NoError(err)

	err = cursor.All(suite.ctx, &userTikets)
	suite.NoError(err)
	suite.Len(userTikets, 1)
	suite.Equal(tiket.ID, userTikets[0].ID)
}

// TestObjectIDValidation menguji validasi ObjectID di database level
func (suite *DatabaseTestSuite) TestObjectIDValidation() {
	collection := suite.db.Collection("users")

	// Test invalid ObjectID handling
	invalidID := "invalid-objectid-format"

	// This should fail because ObjectID format is invalid
	var user models.User
	err := collection.FindOne(suite.ctx, bson.M{"_id": invalidID}).Decode(&user)
	suite.Error(err)

	// Test valid ObjectID
	validID := primitive.NewObjectID()
	err = collection.FindOne(suite.ctx, bson.M{"_id": validID}).Decode(&user)
	suite.Equal(mongo.ErrNoDocuments, err) // Should be "not found" not "invalid format"
}

// TestConcurrentObjectIDOperations menguji operasi ObjectID secara concurrent
func (suite *DatabaseTestSuite) TestConcurrentObjectIDOperations() {
	collection := suite.db.Collection("films")

	// Test concurrent inserts dengan ObjectID
	numGoroutines := 10
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			film := models.Film{
				ID:          primitive.NewObjectID(),
				Title:       "Concurrent Movie " + string(rune(index)),
				Description: "Test concurrent operation",
				Duration:    120,
				Genre:       "Action",
			}

			_, err := collection.InsertOne(suite.ctx, film)
			if err != nil {
				errors <- err
				return
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// Success
		case err := <-errors:
			suite.Fail("Concurrent operation failed", err.Error())
		case <-time.After(5 * time.Second):
			suite.Fail("Timeout waiting for concurrent operations")
		}
	}

	// Verify all records were inserted
	count, err := collection.CountDocuments(suite.ctx, bson.M{})
	suite.NoError(err)
	suite.Equal(int64(numGoroutines), count)
}

// TestDatabaseIndexes menguji indexes untuk ObjectID fields
func (suite *DatabaseTestSuite) TestDatabaseIndexes() {
	// Test indexes pada collections
	collections := []string{"users", "films", "jadwals", "tikets", "pembayarans"}

	for _, collName := range collections {
		collection := suite.db.Collection(collName)

		// Get indexes
		cursor, err := collection.Indexes().List(suite.ctx)
		suite.NoError(err)

		var indexes []bson.M
		err = cursor.All(suite.ctx, &indexes)
		suite.NoError(err)

		// Should have at least _id index
		suite.GreaterOrEqual(len(indexes), 1)

		// Check _id index exists
		found := false
		for _, index := range indexes {
			if key, exists := index["key"]; exists {
				if keyMap, ok := key.(bson.M); ok {
					if _, hasID := keyMap["_id"]; hasID {
						found = true
						break
					}
				}
			}
		}
		suite.True(found, "Collection %s should have _id index", collName)
	}
}

// TestDatabasePerformance menguji performance operasi ObjectID
func (suite *DatabaseTestSuite) TestDatabasePerformance() {
	collection := suite.db.Collection("performance_test")

	// Insert many documents untuk performance test
	numDocs := 1000
	docs := make([]interface{}, numDocs)

	start := time.Now()
	for i := 0; i < numDocs; i++ {
		docs[i] = models.Film{
			ID:          primitive.NewObjectID(),
			Title:       "Performance Test Movie",
			Description: "Test performance",
			Duration:    120,
			Genre:       "Action",
		}
	}

	_, err := collection.InsertMany(suite.ctx, docs)
	suite.NoError(err)

	insertDuration := time.Since(start)
	suite.Less(insertDuration, 5*time.Second, "Bulk insert should complete within 5 seconds")

	// Test query performance
	start = time.Now()
	count, err := collection.CountDocuments(suite.ctx, bson.M{})
	suite.NoError(err)
	suite.Equal(int64(numDocs), count)

	queryDuration := time.Since(start)
	suite.Less(queryDuration, 1*time.Second, "Count query should complete within 1 second")

	// Cleanup
	err = collection.Drop(suite.ctx)
	suite.NoError(err)
}

// TestSuite runner function
func TestDatabaseSuite(t *testing.T) {
	// Skip if MongoDB is not available
	if testing.Short() {
		t.Skip("Skipping database integration tests in short mode")
	}

	suite.Run(t, new(DatabaseTestSuite))
}

// Benchmark functions for database operations
func BenchmarkUserInsert(b *testing.B) {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		// Fallback untuk testing lokal
		os.Setenv("MONGOSTRING", "mongodb://localhost:27017")
	}

	mongoString := os.Getenv("MONGOSTRING")
	if mongoString == "" {
		b.Fatal("MONGOSTRING environment variable not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString))
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}
	db := client.Database("cinema_booking_benchmark")

	collection := db.Collection("users")
	ctx := context.Background()

	// Clear collection
	collection.DeleteMany(ctx, bson.M{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashedPassword, _ := password.HashPassword("benchmarkpassword")

		user := models.User{
			ID:        primitive.NewObjectID(),
			Username:  "benchuser",
			Email:     "bench@example.com",
			Password:  hashedPassword,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			b.Fatalf("Insert failed: %v", err)
		}
	}
}

func BenchmarkUserFind(b *testing.B) {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		// Fallback untuk testing lokal
		os.Setenv("MONGOSTRING", "mongodb://localhost:27017")
	}

	mongoString := os.Getenv("MONGOSTRING")
	if mongoString == "" {
		b.Fatal("MONGOSTRING environment variable not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString))
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}
	db := client.Database("cinema_booking_benchmark")

	collection := db.Collection("users")
	ctx := context.Background()

	// Setup: Insert a user to find
	hashedPassword, _ := password.HashPassword("benchmarkpassword")
	userID := primitive.NewObjectID()

	user := models.User{
		ID:        userID,
		Username:  "benchuser",
		Email:     "bench@example.com",
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection.InsertOne(ctx, user)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var foundUser models.User
		err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&foundUser)
		if err != nil {
			b.Fatalf("Find failed: %v", err)
		}
	}
}
