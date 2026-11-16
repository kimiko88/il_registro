package upload

import (
	"bytes"
	"context"
	"encoding/base64"
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

	// Generate a short-lived URL with server-side credentials. The signed URL
	// never leaves this process; the application authorizes the document first.
	signEndpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", s.URL, url.PathEscape(bucket), escapeObjectPath(objectName))
	payload, err := json.Marshal(map[string]int{"expiresIn": 60})
	if err != nil {
		return nil, "", 0, fmt.Errorf("impossibile preparare la richiesta di firma: %w", err)
	}
	signReq, err := http.NewRequestWithContext(ctx, http.MethodPost, signEndpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, "", 0, fmt.Errorf("impossibile creare la richiesta di firma: %w", err)
	}
	s.setAuthHeaders(signReq)
	signReq.Header.Set("Content-Type", "application/json")
	signResp, err := s.HTTPClient.Do(signReq)
	if err != nil {
		return nil, "", 0, fmt.Errorf("generazione URL firmata non riuscita: %w", err)
	}
	if signResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(signResp.Body, 4096))
		_ = signResp.Body.Close()
		return nil, "", 0, fmt.Errorf("generazione URL firmata fallita (status %d): %s", signResp.StatusCode, string(bodyBytes))
	}
	var signed struct {
		SignedURL string `json:"signedURL"`
	}
	decodeErr := json.NewDecoder(io.LimitReader(signResp.Body, 16<<10)).Decode(&signed)
	_ = signResp.Body.Close()
	if decodeErr != nil || signed.SignedURL == "" {
		return nil, "", 0, fmt.Errorf("risposta URL firmata non valida")
	}

	baseURL, err := url.Parse(s.URL)
	if err != nil {
		return nil, "", 0, fmt.Errorf("URL Supabase non valida")
	}
	signedURL, err := url.Parse(signed.SignedURL)
	if err != nil || signedURL.User != nil || signedURL.Fragment != "" || signedURL.RawQuery == "" {
		return nil, "", 0, fmt.Errorf("URL firmata non valida")
	}
	expectedObjectPath := "/object/sign/" + bucket + "/" + escapeObjectPath(objectName)
	expectedStoragePath := strings.TrimRight(baseURL.Path, "/") + "/storage/v1" + expectedObjectPath
	if signedURL.IsAbs() {
		if signedURL.Scheme != baseURL.Scheme || !strings.EqualFold(signedURL.Host, baseURL.Host) || signedURL.Path != expectedStoragePath {
			return nil, "", 0, fmt.Errorf("URL firmata punta a un host o percorso non autorizzato")
		}
	} else {
		if signedURL.Path != expectedObjectPath {
			return nil, "", 0, fmt.Errorf("percorso URL firmata non valido")
		}
		baseURL.Path = expectedStoragePath
		baseURL.RawQuery = signedURL.RawQuery
		signedURL = baseURL
	}

	// The signed token is the download capability; no Supabase API key or
	// service credential is forwarded on this second request.
	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, signedURL.String(), nil)
	if err != nil {
		return nil, "", 0, fmt.Errorf("impossibile creare la richiesta di download: %w", err)
	}
	resp, err := s.HTTPClient.Do(downloadReq)
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
	if s.URL == "" || strings.TrimSpace(s.Key) == "" {
		return fmt.Errorf("Supabase URL o chiave server non configurate")
	}
	if err := validateServerStorageKey(strings.TrimSpace(s.Key)); err != nil {
		return err
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

// validateServerStorageKey rejects browser/client keys. Storage credentials in
// this process must be a new sb_secret_* key or a legacy JWT with role=service_role.
func validateServerStorageKey(key string) error {
	if strings.HasPrefix(key, "sb_secret_") {
		return nil
	}
	if strings.HasPrefix(key, "sb_publishable_") || strings.HasPrefix(key, "sb_anon_") {
		return fmt.Errorf("le chiavi Supabase publishable/anon non sono ammesse per lo storage server-side")
	}
	parts := strings.Split(key, ".")
	if len(parts) != 3 {
		return fmt.Errorf("chiave Supabase server non valida: usare sb_secret_* o una legacy service_role")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("JWT Supabase server non valido")
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil || claims["role"] != "service_role" {
		return fmt.Errorf("la chiave Supabase deve avere role=service_role; le chiavi anon/client sono vietate")
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
