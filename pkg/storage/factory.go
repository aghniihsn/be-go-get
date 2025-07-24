package storage

import (
	"fmt"
	"os"
)

// GetStorageService returns the configured storage service
func GetStorageService() (StorageService, error) {
	// Check if we should use local storage
	storageProvider := os.Getenv("STORAGE_PROVIDER")
	if storageProvider == "local" {
		baseDir := os.Getenv("LOCAL_STORAGE_DIR")
		if baseDir == "" {
			baseDir = "./uploads"
		}
		baseURL := os.Getenv("LOCAL_STORAGE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:3000/uploads"
		}
		return NewLocalStorageService(baseDir, baseURL)
	}

	if storageProvider == "supabase" {
		return NewSupabaseStorageService()
	}

	return nil, fmt.Errorf("No valid storage provider configured")
}
