package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type SupabaseStorageService struct {
	Url       string
	ApiKey    string
	BucketMap map[FileCategory]string
}

func (s *SupabaseStorageService) Get(fileURL string) (io.Reader, error) {
	return nil, fmt.Errorf("Get not implemented for SupabaseStorageService")
}

func NewSupabaseStorageService() (*SupabaseStorageService, error) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	bucketMap := map[FileCategory]string{
		FilmPoster:     os.Getenv("SUPABASE_BUCKET_POSTERS"),
		ProfilePicture: os.Getenv("SUPABASE_BUCKET_PROFILES"),
		PaymentReceipt: os.Getenv("SUPABASE_BUCKET_RECEIPTS"),
	}
	if url == "" || key == "" {
		return nil, fmt.Errorf("Supabase URL or Service Role Key not set")
	}
	return &SupabaseStorageService{
		Url:       url,
		ApiKey:    key,
		BucketMap: bucketMap,
	}, nil
}

func (s *SupabaseStorageService) Upload(file *multipart.FileHeader, category FileCategory, filename string) (string, error) {
	bucket, ok := s.BucketMap[category]
	if !ok || bucket == "" {
		return "", fmt.Errorf("Invalid or missing bucket for category %s", category)
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.Url, bucket, filename)
	req, err := http.NewRequest("POST", endpoint, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))
	req.Header.Set("x-upsert", "true")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Supabase upload failed: %s", resp.Status)
	}
	// Public URL format: {SUPABASE_URL}/storage/v1/object/public/{bucket}/{filename}
	publicUrl := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.Url, bucket, filename)
	return publicUrl, nil
}

func (s *SupabaseStorageService) Delete(fileURL string) error {
	// Optional: implement delete if needed
	return nil
}
