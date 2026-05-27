package upload

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"backend-mosque-tahfidz-management/internal/config"
)

// CloudinaryResponse represents the JSON response from Cloudinary Upload API
type CloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Format    string `json:"format"`
	Bytes     int    `json:"bytes"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type CloudinaryService struct {
	cloudName string
	apiKey    string
	apiSecret string
}

func NewCloudinaryService(cfg *config.Config) *CloudinaryService {
	return &CloudinaryService{
		cloudName: cfg.CloudinaryCloudName,
		apiKey:    cfg.CloudinaryAPIKey,
		apiSecret: cfg.CloudinaryAPISecret,
	}
}

// Upload sends a file to Cloudinary's Upload API and returns the public secure URL.
// It uses Cloudinary's authenticated upload with a signature — no SDK required.
func (s *CloudinaryService) Upload(file io.Reader, filename string) (*CloudinaryResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	folder := "tahfidz/students"

	// Generate signature: sha1("folder=<folder>&timestamp=<ts><api_secret>")
	sigPayload := fmt.Sprintf("folder=%s&timestamp=%s%s", folder, timestamp, s.apiSecret)
	h := sha1.New()
	h.Write([]byte(sigPayload))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add auth fields
	writer.WriteField("api_key", s.apiKey)
	writer.WriteField("timestamp", timestamp)
	writer.WriteField("signature", signature)
	writer.WriteField("folder", folder)

	writer.Close()

	// POST to Cloudinary
	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.cloudName)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudinary request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudinary upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result CloudinaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode cloudinary response: %w", err)
	}

	return &result, nil
}
