package elearning

import (
	"context"
	"testing"

	"registro-backend/internal/config"
)

func FuzzConnectProvider(f *testing.F) {
	seeds := [][3]string{
		{"teacher-1", "google", "valid_auth_code_123"},
		{"teacher-2", "microsoft", "valid_ms_token_abc"},
		{"teacher-3", "zoom", "some_token"},
		{"", "google", "code"},
		{"teacher-4", "", "code"},
		{"teacher-5", "google", ""},
		{"teacher-6", "MICROSOFT", "code"},
	}

	for _, s := range seeds {
		f.Add(s[0], s[1], s[2])
	}

	svc := NewService(config.ElearningConfig{
		GoogleClientID:     "g-client",
		GoogleClientSecret: "g-secret",
	})
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, userID, provider, authCode string) {
		err := svc.ConnectProvider(ctx, userID, provider, authCode)
		if provider != "google" && provider != "microsoft" {
			if err == nil {
				t.Fatalf("expected error for unsupported provider %q", provider)
			}
			return
		}

		if authCode == "" {
			if err == nil {
				t.Fatalf("expected error for empty authCode")
			}
			return
		}

		// When valid provider and non-empty authCode
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		status, err := svc.GetProvidersStatus(ctx, userID)
		if err != nil {
			t.Fatalf("failed to get provider status: %v", err)
		}

		if provider == "google" && !status.Google.Connected {
			t.Fatalf("expected google to be connected for user %q", userID)
		}
		if provider == "microsoft" && !status.Microsoft.Connected {
			t.Fatalf("expected microsoft to be connected for user %q", userID)
		}
	})
}
