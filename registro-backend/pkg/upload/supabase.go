package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"registro-backend/pkg/circuitbreaker"
)

// StorageUploader defines an interface for uploading files to remote storage.
type StorageUploader interface {
	UploadFile(ctx context.Context, storagePath string, contentType string, reader io.Reader) (string, error)
}

// PrivateStorageDownloader retrieves an object using server-side credentials.
type PrivateStorageDownloader interface {
	DownloadFile(ctx context.Context, storageKey string) (io.ReadCloser, string, int64, error)
}

// StorageObjectDeleter cleans up an upload if the database association fails.
type StorageObjectDeleter interface {
	DeleteFile(ctx context.Context, storageKey string) error
}

// SupabaseUploader uploads files to Supabase Storage using its REST API.
type SupabaseUploader struct {
	URL                   string
	Key                   string
	Bucket                string
	HTTPClient            *http.Client
	PrivateBucketVerified bool
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

// VerifyPrivateBucket checks the actual Storage bucket configuration before enabling file access.
// Public buckets bypass object read authorization, so an unverified bucket fails closed.
func (s *SupabaseUploader) VerifyPrivateBucket(ctx context.Context) error {
	if err := s.validateConfig(); err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s/storage/v1/bucket/%s", s.URL, url.PathEscape(s.bucket()))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("impossibile creare la richiesta di verifica bucket: %w", err)
	}
	s.setAuthHeaders(req)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("verifica bucket Supabase non riuscita: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verifica bucket Supabase fallita (status %d)", resp.StatusCode)
	}
	var info struct {
		Public *bool `json:"public"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 16<<10)).Decode(&info); err != nil {
		return fmt.Errorf("risposta verifica bucket non valida: %w", err)
	}
	if info.Public == nil {
		return fmt.Errorf("risposta verifica bucket priva del flag public")
	}
	if *info.Public {
		return fmt.Errorf("il bucket Supabase %q è pubblico: configurarlo come privato", s.bucket())
	}
	s.PrivateBucketVerified = true
	return nil
}

// UploadFile uploads to a private bucket and returns an opaque storage key.
// The returned value is not a public URL and cannot be fetched without server authorization.
func (s *SupabaseUploader) UploadFile(ctx context.Context, storagePath string, contentType string, reader io.Reader) (string, error) {
	if err := s.validateConfig(); err != nil {
		return "", err
	}
	if !s.PrivateBucketVerified {
		return "", fmt.Errorf("bucket documenti non verificato come privato")
	}
	bucket := s.bucket()
	if err := validateObjectName(storagePath); err != nil {
		return "", err
	}
	uploadEndpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.URL, url.PathEscape(bucket), escapeObjectPath(storagePath))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadEndpoint, reader)
	if err != nil {
		return "", fmt.Errorf("impossibile creare la richiesta di upload: %w", err)
	}
	s.setAuthHeaders(req)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	_, err = circuitbreaker.Storage().Execute(func() (interface{}, error) {
		resp, err := s.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("impossibile eseguire la richiesta di upload su Supabase: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return nil, fmt.Errorf("upload su Supabase Storage fallito (status %d): %s", resp.StatusCode, string(bodyBytes))
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("supabase://%s/%s", bucket, storagePath), nil
}

// DownloadFile retrieves an object from a private bucket using backend credentials.
func (s *SupabaseUploader) DownloadFile(ctx context.Context, storageKey string) (io.ReadCloser, string, int64, error) {
	if err := s.validateConfig(); err != nil {
		return nil, "", 0, err
	}
	if !s.PrivateBucketVerified {
		return nil, "", 0, fmt.Errorf("bucket documenti non verificato come privato")
	}
	bucket, objectName, err := s.parseStorageKey(storageKey)
	if err != nil {
		return nil, "", 0, err
	}
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.URL, url.PathEscape(bucket), escapeObjectPath(objectName))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", 0, fmt.Errorf("impossibile creare la richiesta di download: %w", err)
	}
	s.setAuthHeaders(req)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, "", 0, fmt.Errorf("download da Supabase Storage non riuscito: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		_ = resp.Body.Close()
		return nil, "", 0, fmt.Errorf("download da Supabase Storage fallito (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return resp.Body, contentType, resp.ContentLength, nil
}

// DeleteFile removes an upload that could not be linked to an authorized document.
func (s *SupabaseUploader) DeleteFile(ctx context.Context, storageKey string) error {
	if err := s.validateConfig(); err != nil {
		return err
	}
	if !s.PrivateBucketVerified {
		return fmt.Errorf("bucket documenti non verificato come privato")
	}
	bucket, objectName, err := s.parseStorageKey(storageKey)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s", s.URL, url.PathEscape(bucket))
	payload := strings.NewReader(fmt.Sprintf(`{"prefixes":[%q]}`, objectName))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, payload)
	if err != nil {
		return err
	}
	s.setAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("eliminazione da Supabase Storage fallita (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *SupabaseUploader) validateConfig() error {
	if s.URL == "" || s.Key == "" {
		return fmt.Errorf("Supabase URL o chiave server non configurate")
	}
	parsed, err := url.Parse(s.URL)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("Supabase URL non valida")
	}
	if s.HTTPClient == nil {
		return fmt.Errorf("client HTTP Supabase non configurato")
	}
	return nil
}

func (s *SupabaseUploader) bucket() string {
	if s.Bucket == "" {
		return "documents"
	}
	return s.Bucket
}

func (s *SupabaseUploader) setAuthHeaders(req *http.Request) {
	req.Header.Set("apikey", s.Key)
	// New Supabase secret keys are sent as apikey only. Legacy service_role JWT
	// keys continue to use both headers for backward compatibility.
	if !strings.HasPrefix(s.Key, "sb_secret_") {
		req.Header.Set("Authorization", "Bearer "+s.Key)
	}
}

func validateObjectName(name string) error {
	if name == "" || path.IsAbs(name) || path.Clean(name) != name ||
		strings.ContainsAny(name, "\\?#") || strings.ContainsAny(name, "\r\n\x00") {
		return fmt.Errorf("storage object name non valido")
	}
	for _, segment := range strings.Split(name, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return fmt.Errorf("storage object name non valido")
		}
	}
	return nil
}

func escapeObjectPath(name string) string {
	segments := strings.Split(name, "/")
	for i := range segments {
		segments[i] = url.PathEscape(segments[i])
	}
	return strings.Join(segments, "/")
}
func (s *SupabaseUploader) parseStorageKey(storageKey string) (string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(storageKey))
	if err != nil || parsed.Scheme != "supabase" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", "", fmt.Errorf("storage key non valida")
	}
	bucket := s.bucket()
	objectName := strings.TrimPrefix(parsed.Path, "/")
	if parsed.Host != bucket || objectName == "" {
		return "", "", fmt.Errorf("storage key non valida")
	}
	if err := validateObjectName(objectName); err != nil {
		return "", "", err
	}
	return bucket, objectName, nil
}
