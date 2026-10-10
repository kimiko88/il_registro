package elearning

import (
	"context"
	"fmt"
	"testing"

	"registro-backend/internal/config"
)

func BenchmarkGetProvidersStatus(b *testing.B) {
	svc := NewService(config.ElearningConfig{
		GoogleClientID:        "g-id",
		GoogleClientSecret:    "g-sec",
		MicrosoftClientID:     "ms-id",
		MicrosoftClientSecret: "ms-sec",
	})
	ctx := context.Background()

	_ = svc.ConnectProvider(ctx, "bench-user", "google", "auth-code")

	for b.Loop() {
		_, _ = svc.GetProvidersStatus(ctx, "bench-user")
	}
}

func BenchmarkConnectProvider(b *testing.B) {
	svc := NewService(config.ElearningConfig{
		GoogleClientID:     "g-id",
		GoogleClientSecret: "g-sec",
	})
	ctx := context.Background()

	i := 0
	for b.Loop() {
		userID := fmt.Sprintf("user-%d", i%1000)
		i++
		_ = svc.ConnectProvider(ctx, userID, "google", "auth-code")
	}
}

func BenchmarkSyncAssignments(b *testing.B) {
	svc := NewService(config.ElearningConfig{})
	ctx := context.Background()

	for b.Loop() {
		_, _ = svc.SyncAssignments(ctx, "teacher-1", "google", "class-5a")
	}
}
