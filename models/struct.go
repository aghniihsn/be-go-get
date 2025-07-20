package models

// Film represents a movie in the cinema
// @Description Film information
type Film struct {
	ID       string `json:"id" bson:"id" example:"film001"`
	Title    string `json:"title" bson:"title" example:"Avengers: Endgame"`
	Genre    string `json:"genre" bson:"genre" example:"Action"`
	Duration int    `json:"duration" bson:"duration" example:"181"`
}

// Jadwal represents a movie schedule
// @Description Movie schedule information
type Jadwal struct {
	ID      string  `json:"id" bson:"id" example:"jadwal001"`
	FilmID  string  `json:"film_id" bson:"film_id" example:"film001"`
	Tanggal string  `json:"tanggal" bson:"tanggal" example:"2024-01-15"`
	Waktu   string  `json:"waktu" bson:"waktu" example:"19:00"`
	Ruangan string  `json:"ruangan" bson:"ruangan" example:"Cinema 1"`
	Harga   float64 `json:"harga" bson:"harga" example:"50000"`
}

// Tiket represents a ticket booking
// @Description Ticket booking information
type Tiket struct {
	ID         string  `json:"id" bson:"id" example:"tiket001"`
	JadwalID   string  `json:"jadwal_id" bson:"jadwal_id" example:"jadwal001"`
	Nama       string  `json:"nama" bson:"nama" example:"John Doe"`
	Email      string  `json:"email" bson:"email" example:"john@example.com"`
	Jumlah     int     `json:"jumlah" bson:"jumlah" example:"2"`
	TotalHarga float64 `json:"total_harga" bson:"total_harga" example:"100000"`
	UserID     string  `json:"user_id" bson:"user_id" example:"user001"`
}

// Pembayaran represents a payment record
// @Description Payment information
type Pembayaran struct {
	ID      string  `json:"id" bson:"id" example:"pay001"`
	TiketID string  `json:"tiket_id" bson:"tiket_id" example:"tiket001"`
	Metode  string  `json:"metode" bson:"metode" example:"credit_card"`
	Status  string  `json:"status" bson:"status" example:"pending"`
	Total   float64 `json:"total" bson:"total" example:"100000"`
}

// User represents a user in the system
// @Description User information
type User struct {
	ID       string `json:"id" bson:"id" example:"user001"`
	Nama     string `json:"nama" bson:"nama" example:"John Doe"`
	Email    string `json:"email" bson:"email" example:"john@example.com"`
	Password string `json:"password,omitempty" bson:"password" example:"password123"`
	Role     string `json:"role" bson:"role" example:"user"`
}

// UserLogin represents login request
// @Description User login information
type UserLogin struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// UserRegister represents registration request
// @Description User registration information
type UserRegister struct {
	ID       string `json:"id" example:"user001"`
	Nama     string `json:"nama" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
	Role     string `json:"role" example:"user"`
}

// JWTPayload represents JWT token payload
// @Description JWT token payload
type JWTPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
