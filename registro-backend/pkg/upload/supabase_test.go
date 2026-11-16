package upload

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupabaseUploaderRequiresVerifiedPrivateBucket(t *testing.T) {
	uploader := NewSupabaseUploader("https://storage.example.test", "sb_secret_test", "documents")
	_, err := uploader.UploadFile(context.Background(), "file.pdf", "application/pdf", strings.NewReader("pdf"))
	require.ErrorContains(t, err, "non verificato come privato")

	_, _, _, err = uploader.DownloadFile(context.Background(), "supabase://documents/file.pdf")
	require.ErrorContains(t, err, "non verificato come privato")
}

func TestSupabaseUploaderVerifiesPrivateBucketAndUsesOpaqueStorageKeys(t *testing.T) {
	var uploadCalled, downloadCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "sb_secret_test", r.Header.Get("apikey"))
		require.Empty(t, r.Header.Get("Authorization"), "new Supabase secret keys must not be sent as Bearer JWTs")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/storage/v1/bucket/documents":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":"documents","public":false}`)
		case r.Method == http.MethodPost && r.URL.Path == "/storage/v1/object/documents/communications/school-1/user-1/test.pdf":
			uploadCalled = true
			_, _ = io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPost && r.URL.Path == "/storage/v1/object/sign/documents/communications/school-1/user-1/test.pdf":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"signedURL":"/object/sign/documents/communications/school-1/user-1/test.pdf?token=signed-test-token"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/storage/v1/object/sign/documents/communications/school-1/user-1/test.pdf":
			downloadCalled = true
			require.Equal(t, "token=signed-test-token", r.URL.RawQuery)
			require.Empty(t, r.Header.Get("apikey"))
			require.Empty(t, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Length", "9")
			_, _ = io.WriteString(w, "%PDF-test")
		default:
			t.Fatalf("unexpected storage request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	uploader := NewSupabaseUploader(server.URL, "sb_secret_test", "documents")
	uploader.HTTPClient = server.Client()
	require.NoError(t, uploader.VerifyPrivateBucket(context.Background()))

	key, err := uploader.UploadFile(context.Background(), "communications/school-1/user-1/test.pdf", "application/pdf", strings.NewReader("%PDF-test"))
	require.NoError(t, err)
	require.Equal(t, "supabase://documents/communications/school-1/user-1/test.pdf", key)
	require.NotContains(t, key, "/object/public/")

	reader, contentType, size, err := uploader.DownloadFile(context.Background(), key)
	require.NoError(t, err)
	defer func() { _ = reader.Close() }()
	body, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "application/pdf", contentType)
	require.EqualValues(t, 9, size)
	require.Equal(t, "%PDF-test", string(body))
	require.True(t, uploadCalled)
	require.True(t, downloadCalled)
}

func TestSupabaseUploaderRejectsPublicBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"documents","public":true}`)
	}))
	defer server.Close()

	uploader := NewSupabaseUploader(server.URL, "sb_secret_test", "documents")
	uploader.HTTPClient = server.Client()
	err := uploader.VerifyPrivateBucket(context.Background())
	require.ErrorContains(t, err, "è pubblico")
	require.False(t, uploader.PrivateBucketVerified)
}

func TestSupabaseUploaderRejectsMalformedStorageKeys(t *testing.T) {
	uploader := NewSupabaseUploader("https://storage.example.test", "sb_secret_test", "documents")
	uploader.PrivateBucketVerified = true

	for _, key := range []string{
		"https://storage.example.test/storage/v1/object/public/documents/file.pdf",
		"supabase://other-bucket/file.pdf",
		"supabase://documents/../secrets.txt",
		"supabase://documents/%2e%2e",
		"supabase://documents/file.pdf?download=1",
	} {
		_, _, _, err := uploader.DownloadFile(context.Background(), key)
		require.Error(t, err, "storage key: %s", key)
	}
}

func TestSupabaseUploaderRejectsClientKeys(t *testing.T) {
	require.Error(t, validateServerStorageKey("sb_publishable_test"))
	require.Error(t, validateServerStorageKey("sb_anon_test"))

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	anonPayload := base64.RawURLEncoding.EncodeToString([]byte(`{"role":"anon"}`))
	anonJWT := header + "." + anonPayload + ".signature"
	require.Error(t, validateServerStorageKey(anonJWT))

	servicePayload := base64.RawURLEncoding.EncodeToString([]byte(`{"role":"service_role"}`))
	serviceJWT := header + "." + servicePayload + ".signature"
	require.NoError(t, validateServerStorageKey(serviceJWT))
	require.NoError(t, validateServerStorageKey("sb_secret_server_only_key"))
}
