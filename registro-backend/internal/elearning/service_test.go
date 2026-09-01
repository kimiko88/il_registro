package elearning

import (
	"context"
	"testing"

	"registro-backend/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestElearningService_GetProvidersStatus(t *testing.T) {
	cfg := config.ElearningConfig{
		GoogleClientID:        "google-client-123",
		GoogleClientSecret:    "google-secret-456",
		GoogleRedirectURI:     "https://app.school.it/auth/google",
		MicrosoftClientID:     "ms-client-789",
		MicrosoftClientSecret: "ms-secret-000",
		MicrosoftRedirectURI:  "https://app.school.it/auth/microsoft",
		MicrosoftTenantID:     "ms-tenant-111",
	}

	svc := NewService(cfg)
	ctx := context.Background()

	status, err := svc.GetProvidersStatus(ctx, "user-uuid-1")
	assert.NoError(t, err)
	assert.True(t, status.Google.Configured)
	assert.False(t, status.Google.Connected)
	assert.True(t, status.Microsoft.Configured)
	assert.False(t, status.Microsoft.Connected)
}

func TestElearningService_ConnectProvider(t *testing.T) {
	svc := NewService(config.ElearningConfig{
		GoogleClientID:     "client-id",
		GoogleClientSecret: "secret",
	})
	ctx := context.Background()

	// Unsupported provider
	err := svc.ConnectProvider(ctx, "user-1", "unsupported_provider", "auth-code-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non supportato")

	// Missing auth code
	err = svc.ConnectProvider(ctx, "user-1", "google", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "codice di autorizzazione mancante")

	// Success google
	err = svc.ConnectProvider(ctx, "user-1", "google", "valid-code")
	assert.NoError(t, err)

	status, err := svc.GetProvidersStatus(ctx, "user-1")
	assert.NoError(t, err)
	assert.True(t, status.Google.Connected)
	assert.False(t, status.Microsoft.Connected)

	// Success microsoft
	err = svc.ConnectProvider(ctx, "user-1", "microsoft", "valid-code-ms")
	assert.NoError(t, err)

	status, err = svc.GetProvidersStatus(ctx, "user-1")
	assert.NoError(t, err)
	assert.True(t, status.Google.Connected)
	assert.True(t, status.Microsoft.Connected)
}

func TestElearningService_SyncEndpoints(t *testing.T) {
	svc := NewService(config.ElearningConfig{})
	ctx := context.Background()

	resCourses, err := svc.SyncCourses(ctx, "user-1", "google")
	assert.NoError(t, err)
	assert.Contains(t, resCourses.Message, "google")
	assert.Greater(t, resCourses.Count, 0)

	resAssignments, err := svc.SyncAssignments(ctx, "user-1", "google", "class-1")
	assert.NoError(t, err)
	assert.Greater(t, resAssignments.Count, 0)

	resGrades, err := svc.SyncGrades(ctx, "user-1", "google", "class-1")
	assert.NoError(t, err)
	assert.Greater(t, resGrades.Count, 0)
}
