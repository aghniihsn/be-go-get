package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Models structs
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Film struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Genre       string             `json:"genre" bson:"genre"`
	Duration    int                `json:"duration" bson:"duration"`
	Rating      string             `json:"rating" bson:"rating"`
	Description string             `json:"description" bson:"description"`
	PosterURL   string             `json:"poster_url" bson:"poster_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Jadwal struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FilmID    primitive.ObjectID `json:"film_id" bson:"film_id"`
	Tanggal   string             `json:"tanggal" bson:"tanggal"`
	Waktu     string             `json:"waktu" bson:"waktu"`
	Ruangan   string             `json:"ruangan" bson:"ruangan"`
	Harga     float64            `json:"harga" bson:"harga"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Tiket struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	JadwalID         primitive.ObjectID `json:"jadwal_id" bson:"jadwal_id"`
	Kursi            string             `json:"kursi" bson:"kursi"`
	Status           string             `json:"status" bson:"status"`
	TanggalPembelian time.Time          `json:"tanggal_pembelian" bson:"tanggal_pembelian"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type Pembayaran struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TiketID           primitive.ObjectID `json:"tiket_id" bson:"tiket_id"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	Jumlah            float64            `json:"jumlah" bson:"jumlah"`
	MetodePembayaran  string             `json:"metode_pembayaran" bson:"metode_pembayaran"`
	Status            string             `json:"status" bson:"status"`
	TanggalPembayaran time.Time          `json:"tanggal_pembayaran" bson:"tanggal_pembayaran"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

var DB *mongo.Database

func main() {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Connect to database
	connectDB()

	fmt.Println("🌱 Starting Database Seeder...")

	// Run seeders in order (dependencies)
	seedUsers()
	seedFilms()
	seedJadwals()
	seedTikets()
	seedPembayarans()

	fmt.Println("✅ Database seeding completed successfully!")
}

// connectDB establishes database connection
func connectDB() {
	mongoURI := os.Getenv("MONGOSTRING")
	if mongoURI == "" {
		log.Fatal("MONGOSTRING environment variable is not set")
	}

	// Add database name to URI if not present
	if mongoURI[len(mongoURI)-1] != '/' {
		mongoURI += "/"
	}

	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "cinema_booking"
	}

	mongoURI += dbName

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	DB = client.Database(dbName)
	fmt.Println("✅ Connected to MongoDB!")
}

// seedUsers creates sample users including admin
func seedUsers() {
	collection := DB.Collection("users")
	ctx := context.Background()

	fmt.Println("🔹 Seeding Users...")

	// Clear existing users
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear users collection: %v", err)
	}

	// Sample users data
	users := []User{
		{
			ID:        primitive.NewObjectID(),
			Username:  "admin",
			Email:     "admin@cinema.com",
			Password:  hashPassword("admin123"),
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "johndoe",
			Email:     "john@example.com",
			Password:  hashPassword("password123"),
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "janedoe",
			Email:     "jane@example.com",
			Password:  hashPassword("password123"),
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "bobsmith",
			Email:     "bob@example.com",
			Password:  hashPassword("password123"),
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Username:  "alicechen",
			Email:     "alice@example.com",
			Password:  hashPassword("password123"),
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert users
	var docs []interface{}
	for _, user := range users {
		docs = append(docs, user)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	fmt.Printf("✅ Seeded %d users\n", len(result.InsertedIDs))

	// Store user IDs for later use
	for i, insertedID := range result.InsertedIDs {
		users[i].ID = insertedID.(primitive.ObjectID)
	}
}

// seedFilms creates sample films
func seedFilms() {
	collection := DB.Collection("films")
	ctx := context.Background()

	fmt.Println("🔹 Seeding Films...")

	// Clear existing films
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear films collection: %v", err)
	}

	// Sample films data
	films := []Film{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Avengers: Endgame",
			Genre:       "Action/Adventure",
			Duration:    181,
			Rating:      "PG-13",
			Description: "The epic conclusion to the Infinity Saga that became a critically acclaimed worldwide phenomenon.",
			PosterURL:   "https://example.com/posters/avengers-endgame.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Spider-Man: No Way Home",
			Genre:       "Action/Adventure",
			Duration:    148,
			Rating:      "PG-13",
			Description: "Spider-Man's identity is revealed and he asks Doctor Strange for help.",
			PosterURL:   "https://example.com/posters/spiderman-nwh.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Top Gun: Maverick",
			Genre:       "Action/Drama",
			Duration:    131,
			Rating:      "PG-13",
			Description: "After more than thirty years of service as one of the Navy's top aviators, Pete Mitchell is where he belongs.",
			PosterURL:   "https://example.com/posters/top-gun-maverick.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "The Batman",
			Genre:       "Action/Crime",
			Duration:    176,
			Rating:      "PG-13",
			Description: "Batman ventures into Gotham City's underworld when a sadistic killer leaves behind a trail of cryptic clues.",
			PosterURL:   "https://example.com/posters/the-batman.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Dune",
			Genre:       "Sci-Fi/Adventure",
			Duration:    155,
			Rating:      "PG-13",
			Description: "A noble family becomes embroiled in a war for control over the galaxy's most valuable asset.",
			PosterURL:   "https://example.com/posters/dune.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Parasite",
			Genre:       "Thriller/Drama",
			Duration:    132,
			Rating:      "R",
			Description: "A poor family schemes to become employed by a wealthy family by infiltrating their household.",
			PosterURL:   "https://example.com/posters/parasite.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert films
	var docs []interface{}
	for _, film := range films {
		docs = append(docs, film)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to seed films: %v", err)
	}

	fmt.Printf("✅ Seeded %d films\n", len(result.InsertedIDs))

	// Store film IDs for later use
	for i, insertedID := range result.InsertedIDs {
		films[i].ID = insertedID.(primitive.ObjectID)
	}
}

// seedJadwals creates sample schedules for films
func seedJadwals() {
	collection := DB.Collection("jadwals")
	filmCollection := DB.Collection("films")
	ctx := context.Background()

	fmt.Println("🔹 Seeding Jadwals...")

	// Clear existing jadwals
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear jadwals collection: %v", err)
	}

	// Get film IDs
	cursor, err := filmCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to get films for jadwal seeding: %v", err)
	}

	var films []Film
	if err = cursor.All(ctx, &films); err != nil {
		log.Fatalf("Failed to decode films: %v", err)
	}

	if len(films) == 0 {
		log.Fatal("No films found for jadwal seeding. Please seed films first.")
	}

	// Sample schedules
	var jadwals []Jadwal

	// Generate schedules for each film
	ruangans := []string{"Cinema 1", "Cinema 2", "Cinema 3", "IMAX", "VIP"}
	times := []string{"10:00", "13:00", "16:00", "19:00", "22:00"}
	dates := []string{"2025-07-22", "2025-07-23", "2025-07-24", "2025-07-25", "2025-07-26"}

	for _, film := range films {
		for i, date := range dates {
			if i >= 3 { // Limit to 3 schedules per film
				break
			}

			jadwal := Jadwal{
				ID:        primitive.NewObjectID(),
				FilmID:    film.ID,
				Tanggal:   date,
				Waktu:     times[i%len(times)],
				Ruangan:   ruangans[i%len(ruangans)],
				Harga:     float64(35000 + (i * 15000)), // Vary prices
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			jadwals = append(jadwals, jadwal)
		}
	}

	// Insert jadwals
	var docs []interface{}
	for _, jadwal := range jadwals {
		docs = append(docs, jadwal)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to seed jadwals: %v", err)
	}

	fmt.Printf("✅ Seeded %d jadwals\n", len(result.InsertedIDs))
}

// seedTikets creates sample tickets
func seedTikets() {
	collection := DB.Collection("tikets")
	userCollection := DB.Collection("users")
	jadwalCollection := DB.Collection("jadwals")
	ctx := context.Background()

	fmt.Println("🔹 Seeding Tikets...")

	// Clear existing tikets
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear tikets collection: %v", err)
	}

	// Get users (exclude admin)
	cursor, err := userCollection.Find(ctx, bson.M{"role": "user"})
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatalf("Failed to decode users: %v", err)
	}

	// Get jadwals
	cursor, err = jadwalCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to get jadwals: %v", err)
	}

	var jadwals []Jadwal
	if err = cursor.All(ctx, &jadwals); err != nil {
		log.Fatalf("Failed to decode jadwals: %v", err)
	}

	if len(users) == 0 || len(jadwals) == 0 {
		log.Fatal("No users or jadwals found for tiket seeding")
	}

	// Sample tickets
	var tikets []Tiket
	seats := []string{"A1", "A2", "B1", "B2", "C1", "C2", "D1", "D2"}
	statuses := []string{"confirmed", "used", "cancelled"}

	// Create tickets for different users and schedules
	for i, user := range users {
		for j := 0; j < 2; j++ { // 2 tickets per user
			if (i*2 + j) >= len(jadwals) {
				break
			}

			jadwal := jadwals[i*2+j]
			tiket := Tiket{
				ID:               primitive.NewObjectID(),
				UserID:           user.ID,
				JadwalID:         jadwal.ID,
				Kursi:            seats[(i*2+j)%len(seats)],
				Status:           statuses[j%len(statuses)],
				TanggalPembelian: time.Now().AddDate(0, 0, -j-1), // Vary purchase dates
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			tikets = append(tikets, tiket)
		}
	}

	// Insert tikets
	var docs []interface{}
	for _, tiket := range tikets {
		docs = append(docs, tiket)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to seed tikets: %v", err)
	}

	fmt.Printf("✅ Seeded %d tikets\n", len(result.InsertedIDs))
}

// seedPembayarans creates sample payments
func seedPembayarans() {
	collection := DB.Collection("pembayarans")
	tiketCollection := DB.Collection("tikets")
	ctx := context.Background()

	fmt.Println("🔹 Seeding Pembayarans...")

	// Clear existing pembayarans
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear pembayarans collection: %v", err)
	}

	// Get tikets that are confirmed or used (not cancelled)
	cursor, err := tiketCollection.Find(ctx, bson.M{"status": bson.M{"$in": []string{"confirmed", "used"}}})
	if err != nil {
		log.Fatalf("Failed to get tikets: %v", err)
	}

	var tikets []Tiket
	if err = cursor.All(ctx, &tikets); err != nil {
		log.Fatalf("Failed to decode tikets: %v", err)
	}

	if len(tikets) == 0 {
		log.Fatal("No valid tikets found for pembayaran seeding")
	}

	// Get jadwal info for pricing
	jadwalCollection := DB.Collection("jadwals")

	var pembayarans []Pembayaran
	metodes := []string{"credit_card", "debit_card", "bank_transfer", "e_wallet", "cash"}
	statuses := []string{"completed", "pending", "failed"}

	for i, tiket := range tikets {
		// Get jadwal for price
		var jadwal Jadwal
		err := jadwalCollection.FindOne(ctx, bson.M{"_id": tiket.JadwalID}).Decode(&jadwal)
		if err != nil {
			log.Printf("Warning: Could not find jadwal for tiket %s: %v", tiket.ID.Hex(), err)
			continue
		}

		pembayaran := Pembayaran{
			ID:                primitive.NewObjectID(),
			TiketID:           tiket.ID,
			UserID:            tiket.UserID,
			Jumlah:            jadwal.Harga, // Use jadwal price
			MetodePembayaran:  metodes[i%len(metodes)],
			Status:            statuses[i%len(statuses)],
			TanggalPembayaran: tiket.TanggalPembelian.Add(time.Minute * 5), // Payment 5 minutes after ticket purchase
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		pembayarans = append(pembayarans, pembayaran)
	}

	// Insert pembayarans
	var docs []interface{}
	for _, pembayaran := range pembayarans {
		docs = append(docs, pembayaran)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to seed pembayarans: %v", err)
	}

	fmt.Printf("✅ Seeded %d pembayarans\n", len(result.InsertedIDs))
}

// hashPassword utility function
func hashPassword(plainPassword string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(bytes)
}
