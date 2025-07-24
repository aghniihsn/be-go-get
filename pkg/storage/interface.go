package storage

import (
	"io"
	"mime/multipart"
)

// FileCategory menentukan jenis file yang disimpan
type FileCategory string

const (
	FilmPoster     FileCategory = "film_posters"
	ProfilePicture FileCategory = "profile_pictures"
	PaymentReceipt FileCategory = "payment_receipts"
)

// StorageService interface untuk operasi penyimpanan file
type StorageService interface {
	// Upload menyimpan file dan mengembalikan URL publik
	Upload(file *multipart.FileHeader, category FileCategory, filename string) (string, error)

	// Delete menghapus file berdasarkan URL atau path
	Delete(fileURL string) error

	// Get mendapatkan file sebagai io.Reader (optional untuk implementasi)
	Get(fileURL string) (io.Reader, error)
}
