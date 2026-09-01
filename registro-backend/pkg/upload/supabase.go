package upload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StorageUploader defines an interface for uploading files to remote storage.
type StorageUploader interface {
	UploadFile(ctx context.Context, storagePath string, contentType string, reader io.Reader) (string, error)
}

// SupabaseUploader uploads files to Supabase Storage using its REST API.
type SupabaseUploader struct {
	URL        string
	Key        string
	Bucket     string
	HTTPClient *http.Client
}

// NewSupabaseUploader initializes a new SupabaseUploader instance.
func NewSupabaseUploader(url, key, bucket string) *SupabaseUploader {
	return &SupabaseUploader{
		URL:        strings.TrimRight(url, "/"),
		Key:        key,
		Bucket:     bucket,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// UploadFile uploads a file stream to Supabase Storage and returns its public URL.
func (s *SupabaseUploader) UploadFile(ctx context.Context, storagePath string, contentType string, reader io.Reader) (string, error) {
	if s.URL == "" || s.Key == "" {
		return "", fmt.Errorf("supabase URL o Key non configurate")
	}

	bucket := s.Bucket
	if bucket == "" {
		bucket = "documents"
	}

	uploadEndpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.URL, bucket, storagePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadEndpoint, reader)
	if err != nil {
		return "", fmt.Errorf("impossibile creare la richiesta di upload: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Key)
	req.Header.Set("apikey", s.Key)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("impossibile eseguire la richiesta di upload su Supabase: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload su Supabase Storage fallito (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.URL, bucket, storagePath)
	return publicURL, nil
}
