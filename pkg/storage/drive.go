package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveService implementasi StorageService dengan Google Drive
type DriveService struct {
	service       *drive.Service
	folderIDs     map[FileCategory]string // Menyimpan folder ID untuk setiap kategori
	publicBaseURL string                  // Base URL untuk akses file publik
}

// NewDriveService membuat instance baru DriveService
func NewDriveService(credentialsFile string) (*DriveService, error) {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// We're using service account credentials
	jwtConfig, err := google.JWTConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse service account file to config: %v", err)
	}

	// Get client from JWT config
	client := jwtConfig.Client(context.Background())

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	// Initialize folder IDs map - load from environment variables
	folderIDs := map[FileCategory]string{
		FilmPoster:     os.Getenv("DRIVE_FOLDER_POSTERS"),
		ProfilePicture: os.Getenv("DRIVE_FOLDER_PROFILES"),
		PaymentReceipt: os.Getenv("DRIVE_FOLDER_RECEIPTS"),
	}

	// Check if folder IDs are set, if not try to create folders
	for category, id := range folderIDs {
		if id == "" {
			folderID, err := createDriveFolder(srv, string(category))
			if err != nil {
				return nil, fmt.Errorf("failed to create folder for %s: %v", category, err)
			}
			folderIDs[category] = folderID
		}
	}

	return &DriveService{
		service:       srv,
		folderIDs:     folderIDs,
		publicBaseURL: "https://drive.google.com/uc?export=view&id=",
	}, nil
}

// Service account authentication eliminates the need for OAuth2 token handling.
// JWT tokens are generated automatically from the private key in the service account credentials.

// createDriveFolder creates a new folder in Google Drive and returns its ID
func createDriveFolder(service *drive.Service, folderName string) (string, error) {
	// Define folder metadata
	folderMetadata := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	// Create the folder
	folder, err := service.Files.Create(folderMetadata).Do()
	if err != nil {
		return "", fmt.Errorf("could not create folder: %v", err)
	}

	// Make the folder publicly readable
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err = service.Permissions.Create(folder.Id, permission).Do()
	if err != nil {
		return "", fmt.Errorf("could not share folder: %v", err)
	}

	return folder.Id, nil
}

// Upload mengimplementasikan StorageService.Upload
func (d *DriveService) Upload(file *multipart.FileHeader, category FileCategory, filename string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	folderID, ok := d.folderIDs[category]
	if !ok {
		return "", fmt.Errorf("category folder not configured: %s", category)
	}

	// Tambahkan timestamp ke nama file untuk mencegah duplikasi
	timestamp := time.Now().UnixNano()
	// Bersihkan nama file dari karakter yang tidak diinginkan
	cleanFilename := strings.ReplaceAll(filename, " ", "_")
	finalFilename := fmt.Sprintf("%d_%s", timestamp, cleanFilename)

	f := &drive.File{
		Name:     finalFilename,
		Parents:  []string{folderID},
		MimeType: file.Header.Get("Content-Type"),
	}

	// Buat file di Drive
	fileContent, err := ioutil.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("could not read file content: %v", err)
	}

	createdFile, err := d.service.Files.Create(f).Media(strings.NewReader(string(fileContent))).Do()
	if err != nil {
		return "", fmt.Errorf("could not create file: %v", err)
	}

	// Update izin file untuk dapat diakses publik
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err = d.service.Permissions.Create(createdFile.Id, permission).Do()
	if err != nil {
		return "", fmt.Errorf("could not share file: %v", err)
	}

	// Kembalikan URL publik
	return d.publicBaseURL + createdFile.Id, nil
}

// Delete mengimplementasikan StorageService.Delete
func (d *DriveService) Delete(fileURL string) error {
	// Ekstrak file ID dari URL
	fileID := extractFileIDFromURL(fileURL)
	if fileID == "" {
		return fmt.Errorf("invalid file URL format")
	}

	err := d.service.Files.Delete(fileID).Do()
	if err != nil {
		return fmt.Errorf("could not delete file: %v", err)
	}

	return nil
}

// Get mengimplementasikan StorageService.Get
func (d *DriveService) Get(fileURL string) (io.Reader, error) {
	fileID := extractFileIDFromURL(fileURL)
	if fileID == "" {
		return nil, fmt.Errorf("invalid file URL format")
	}

	// Google Drive files can be accessed directly via the URL
	// We're just verifying the file exists
	_, err := d.service.Files.Get(fileID).Do()
	if err != nil {
		return nil, fmt.Errorf("could not find file: %v", err)
	}

	// Return the file content via HTTP request
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("could not fetch file: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get file, status: %s", resp.Status)
	}

	return resp.Body, nil
}

// Fungsi utilitas untuk mengekstrak ID file dari URL
func extractFileIDFromURL(url string) string {
	const prefix = "https://drive.google.com/uc?export=view&id="
	if strings.HasPrefix(url, prefix) {
		return url[len(prefix):]
	}
	return ""
}
