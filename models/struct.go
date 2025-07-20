package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents user data structure
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"role" bson:"role"` // "user" or "admin"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Film represents movie data structure
type Film struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Genre       string             `json:"genre" bson:"genre"`
	Duration    int                `json:"duration" bson:"duration"` // in minutes
	Rating      string             `json:"rating" bson:"rating"`
	Description string             `json:"description" bson:"description"`
	PosterURL   string             `json:"poster_url" bson:"poster_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Jadwal represents movie schedule data structure
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

// Tiket represents ticket data structure
type Tiket struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	JadwalID         primitive.ObjectID `json:"jadwal_id" bson:"jadwal_id"`
	Kursi            string             `json:"kursi" bson:"kursi"`
	Status           string             `json:"status" bson:"status"` // "confirmed", "cancelled", "used"
	TanggalPembelian time.Time          `json:"tanggal_pembelian" bson:"tanggal_pembelian"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// Pembayaran represents payment data structure
type Pembayaran struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TiketID           primitive.ObjectID `json:"tiket_id" bson:"tiket_id"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	Jumlah            float64            `json:"jumlah" bson:"jumlah"`
	MetodePembayaran  string             `json:"metode_pembayaran" bson:"metode_pembayaran"`
	Status            string             `json:"status" bson:"status"` // "pending", "completed", "failed"
	TanggalPembayaran time.Time          `json:"tanggal_pembayaran" bson:"tanggal_pembayaran"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserLogin represents login request
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegister represents registration request
type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// JWTPayload represents JWT token payload
type JWTPayload struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Role     string             `json:"role"`
}
