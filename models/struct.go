package models

type Film struct {
	ID       string `json:"id" bson:"id"`
	Title    string `json:"title" bson:"title"`
	Genre    string `json:"genre" bson:"genre"`
	Duration int    `json:"duration" bson:"duration"`
}

type Jadwal struct {
	ID      string  `json:"id" bson:"id"`
	FilmID  string  `json:"film_id" bson:"film_id"`
	Tanggal string  `json:"tanggal" bson:"tanggal"`
	Waktu   string  `json:"waktu" bson:"waktu"`
	Ruangan string  `json:"ruangan" bson:"ruangan"`
	Harga   float64 `json:"harga" bson:"harga"`
}

type Tiket struct {
	ID         string  `json:"id" bson:"id"`
	JadwalID   string  `json:"jadwal_id" bson:"jadwal_id"`
	Nama       string  `json:"nama" bson:"nama"`
	Email      string  `json:"email" bson:"email"`
	Jumlah     int     `json:"jumlah" bson:"jumlah"`
	TotalHarga float64 `json:"total_harga" bson:"total_harga"`
	UserID     string  `json:"user_id" bson:"user_id"`
}

type Pembayaran struct {
	ID      string  `json:"id" bson:"id"`
	TiketID string  `json:"tiket_id" bson:"tiket_id"`
	Metode  string  `json:"metode" bson:"metode"`
	Status  string  `json:"status" bson:"status"`
	Total   float64 `json:"total" bson:"total"`
}

type User struct {
	ID    string `json:"id" bson:"id"`
	Nama  string `json:"nama" bson:"nama"`
	Email string `json:"email" bson:"email"`
}
