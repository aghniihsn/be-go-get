package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTiketRequest struct {
	JadwalID string `json:"jadwal_id"`
	UserID   string `json:"user_id"`
	Kursi    string `json:"kursi"`
	Status   string `json:"status"`
}

var StudioSeats = []string{
	"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10",
	"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B10",
	"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "C10",
	"D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10",
	"E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9", "E10",
}

type User struct {
	ID                primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Username          string             `json:"username" bson:"username"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"password" bson:"password"`
	Role              string             `json:"role" bson:"role"`
	Firstname         string             `json:"firstname" bson:"firstname"`
	Lastname          string             `json:"lastname" bson:"lastname"`
	Gender            string             `json:"gender" bson:"gender"`
	PhoneNumber       string             `json:"phone_number" bson:"phone_number"`
	ProfilePictureURL string             `json:"profile_picture_url" bson:"profile_picture_url"`
	Address           string             `json:"address" bson:"address"`
	CreatedAt         time.Time          `json:"-" bson:"created_at"`
	UpdatedAt         time.Time          `json:"-" bson:"updated_at"`
}

type Film struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Genre       []string           `json:"genre" bson:"genre"`
	Duration    int                `json:"duration" bson:"duration"`
	Rating      string             `json:"rating" bson:"rating"`
	Description string             `json:"description" bson:"description"`
	PosterURL   string             `json:"poster_url" bson:"poster_url"`
	CreatedAt   time.Time          `json:"-" bson:"created_at"`
	UpdatedAt   time.Time          `json:"-" bson:"updated_at"`
}

const (
	RatingSemua  = "Semua Umur"
	RatingAnak   = "Anak-anak"
	RatingRemaja = "Remaja"
	RatingDewasa = "Dewasa"
)

type Jadwal struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FilmID    primitive.ObjectID `json:"film_id" bson:"film_id"`
	FilmTitle string             `json:"film_title" bson:"film_title"`
	Tanggal   string             `json:"tanggal" bson:"tanggal"`
	Waktu     string             `json:"waktu" bson:"waktu"`
	Ruangan   string             `json:"ruangan" bson:"ruangan"`
	Harga     float64            `json:"harga" bson:"harga"`
	CreatedAt time.Time          `json:"-" bson:"created_at"`
	UpdatedAt time.Time          `json:"-" bson:"updated_at"`
}

var ValidStudios = []string{"Studio 1", "Studio 2", "Studio 3", "Studio 4", "Studio 5"}

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

// Enum status tiket
const (
	TiketStatusConfirmed = "confirmed"
	TiketStatusCancelled = "cancelled"
	TiketStatusUsed      = "used"
)

type Pembayaran struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TiketID           primitive.ObjectID `json:"tiket_id" bson:"tiket_id"`
	Jumlah            float64            `json:"jumlah" bson:"jumlah"`
	MetodePembayaran  string             `json:"metode_pembayaran" bson:"metode_pembayaran"`
	Status            string             `json:"status" bson:"status"`
	BuktiPembayaran   string             `json:"bukti_pembayaran,omitempty" bson:"bukti_pembayaran,omitempty"`
	TanggalPembayaran time.Time          `json:"tanggal_pembayaran" bson:"tanggal_pembayaran"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

const (
	PembayaranStatusPending   = "pending"
	PembayaranStatusCompleted = "completed"
	PembayaranStatusFailed    = "failed"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Address           string `json:"address"`
}

type JWTPayload struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Role     string             `json:"role"`
}
