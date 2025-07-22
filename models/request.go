package models

import "time"

// CreateBatchTicketsRequest represents the request to create multiple tickets at once
// swagger:model CreateBatchTicketsRequest
type CreateBatchTicketsRequest struct {
	// The ID of the schedule for which tickets are being booked
	// example: 60d21b4667d0d8992e610c85
	JadwalID string `json:"jadwal_id"`
	// The ID of the user booking the tickets
	// example: 60d21b4667d0d8992e610c84
	UserID string `json:"user_id"`
	// Array of seat codes to book
	// example: ["A1", "A2", "A3"]
	Kursis []string `json:"kursis"`
}

// PaymentRequest represents a payment request with additional information
// swagger:model PaymentRequest
type PaymentRequest struct {
	// The ID of the ticket being paid for
	// example: 60d21b4667d0d8992e610c86
	TiketID string `json:"tiket_id"`
	// The ID of the user making the payment
	// example: 60d21b4667d0d8992e610c84
	UserID string `json:"user_id"`
	// The amount being paid
	// example: 50000
	Jumlah float64 `json:"jumlah"`
	// The payment method used
	// example: transfer_bank
	MetodePembayaran string `json:"metode_pembayaran"`
	// URL or base64 of payment proof (optional)
	// example: https://example.com/payment-proof.jpg
	BuktiPembayaran string `json:"bukti_pembayaran,omitempty"`
}

// TicketSummary represents a ticket with additional information for display
// swagger:model TicketSummary
type TicketSummary struct {
	// The ticket details
	Tiket Tiket `json:"tiket"`
	// The title of the film
	// example: Avengers: Endgame
	FilmTitle string `json:"film_title"`
	// The URL to the film's poster image
	// example: https://example.com/posters/avengers.jpg
	FilmPoster string `json:"film_poster"`
	// The scheduled time for the film
	// example: 19:30
	JadwalWaktu string `json:"jadwal_waktu"`
	// The scheduled date for the film
	// example: 2025-07-25
	JadwalTanggal string `json:"jadwal_tanggal"`
	// The theater room for the film
	// example: Studio 3
	JadwalRuangan string `json:"jadwal_ruangan"`
	// The ticket price
	// example: 50000
	HargaTiket float64 `json:"harga_tiket"`
	// The payment status
	// example: completed
	StatusPembayaran string `json:"status_pembayaran"`
	// The payment date (if paid)
	TanggalPembayaran time.Time `json:"tanggal_pembayaran,omitempty"`
}

// MetodePembayaran represents available payment methods
// swagger:model MetodePembayaran
type MetodePembayaran struct {
	// The payment method ID
	// example: transfer_bank
	ID string `json:"id"`
	// The display name of the payment method
	// example: Transfer Bank
	Name string `json:"name"`
	// A description of the payment method
	// example: Transfer ke rekening bank yang tersedia
	Description string `json:"description"`
	// Optional URL to the payment method's logo
	// example: /images/transfer_bank.png
	Logo string `json:"logo,omitempty"`
}

// Available payment methods
var AvailablePaymentMethods = []MetodePembayaran{
	{
		ID:          "transfer_bank",
		Name:        "Transfer Bank",
		Description: "Transfer ke rekening bank yang tersedia",
		Logo:        "/images/transfer_bank.png",
	},
	{
		ID:          "e_wallet",
		Name:        "E-Wallet",
		Description: "Pembayaran menggunakan e-wallet seperti OVO, GoPay, dll",
		Logo:        "/images/e_wallet.png",
	},
	{
		ID:          "kartu_kredit",
		Name:        "Kartu Kredit",
		Description: "Pembayaran menggunakan kartu kredit",
		Logo:        "/images/credit_card.png",
	},
}
