package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorageService implementasi StorageService dengan penyimpanan lokal
type LocalStorageService struct {
	baseDir      string                  // Direktori dasar untuk penyimpanan
	baseURL      string                  // URL dasar untuk akses file
	categoryDirs map[FileCategory]string // Map untuk direktori kategori
}

// NewLocalStorageService membuat instance baru LocalStorageService
func NewLocalStorageService(baseDir, baseURL string) (*LocalStorageService, error) {
	if baseDir == "" {
		baseDir = "./uploads"
	}

	if baseURL == "" {
		baseURL = "http://localhost:3000/uploads"
	}

	// Buat map direktori kategori
	categoryDirs := map[FileCategory]string{
		FilmPoster:     filepath.Join(baseDir, "film_posters"),
		ProfilePicture: filepath.Join(baseDir, "profile_pictures"),
		PaymentReceipt: filepath.Join(baseDir, "payment_receipts"),
	}

	// Pastikan direktori ada
	for _, dir := range categoryDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("could not create directory %s: %v", dir, err)
		}
	}

	return &LocalStorageService{
		baseDir:      baseDir,
		baseURL:      baseURL,
		categoryDirs: categoryDirs,
	}, nil
}

// Upload mengimplementasikan StorageService.Upload
func (l *LocalStorageService) Upload(file *multipart.FileHeader, category FileCategory, filename string) (string, error) {
	categoryDir, exists := l.categoryDirs[category]
	if !exists {
		return "", fmt.Errorf("unknown category: %s", category)
	}

	// Buka file sumber
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("could not open uploaded file: %v", err)
	}
	defer src.Close()

	// Tambahkan timestamp ke nama file untuk mencegah duplikasi
	timestamp := time.Now().UnixNano()
	// Bersihkan nama file dari karakter yang tidak diinginkan
	cleanFilename := strings.ReplaceAll(filename, " ", "_")
	finalFilename := fmt.Sprintf("%d_%s", timestamp, cleanFilename)

	// Buat path file tujuan
	dstPath := filepath.Join(categoryDir, finalFilename)

	// Buat file tujuan
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("could not create destination file: %v", err)
	}
	defer dst.Close()

	// Salin konten dari file yang diupload ke file tujuan
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("could not copy file content: %v", err)
	}

	// Konstruksi URL untuk mengakses file
	fileURL := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight(l.baseURL, "/"),
		string(category),
		finalFilename,
	)

	return fileURL, nil
}

// Delete mengimplementasikan StorageService.Delete
func (l *LocalStorageService) Delete(fileURL string) error {
	// Ekstrak path file dari URL
	filePath, err := l.extractFilePathFromURL(fileURL)
	if err != nil {
		return err
	}

	// Hapus file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete file: %v", err)
	}

	return nil
}

// Get mengimplementasikan StorageService.Get
func (l *LocalStorageService) Get(fileURL string) (io.Reader, error) {
	// Ekstrak path file dari URL
	filePath, err := l.extractFilePathFromURL(fileURL)
	if err != nil {
		return nil, err
	}

	// Buka file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}

	return file, nil
}

// extractFilePathFromURL mengekstrak path file lokal dari URL
func (l *LocalStorageService) extractFilePathFromURL(fileURL string) (string, error) {
	if !strings.HasPrefix(fileURL, l.baseURL) {
		return "", fmt.Errorf("URL does not match base URL")
	}

	// Extract relative path from URL
	relPath := strings.TrimPrefix(fileURL, l.baseURL)
	relPath = strings.TrimPrefix(relPath, "/")

	// Map to actual file path
	parts := strings.SplitN(relPath, "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid URL format")
	}

	category := parts[0]
	filename := parts[1]

	var categoryPath string
	switch category {
	case string(FilmPoster):
		categoryPath = l.categoryDirs[FilmPoster]
	case string(ProfilePicture):
		categoryPath = l.categoryDirs[ProfilePicture]
	case string(PaymentReceipt):
		categoryPath = l.categoryDirs[PaymentReceipt]
	default:
		return "", fmt.Errorf("unknown category in URL: %s", category)
	}

	filePath := filepath.Join(categoryPath, filename)
	return filePath, nil
}
