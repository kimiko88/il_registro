package ws

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWS_AllowedOrigins_Memoized(t *testing.T) {
	_ = os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
	defer func() { _ = os.Unsetenv("ALLOWED_ORIGINS") }()

	origins := getAllowedOrigins()
	assert.True(t, origins["http://localhost:3000"])
	assert.True(t, origins["http://example.com"])
	assert.False(t, origins["http://malicious.com"])
}

func TestWS_Hub_GlobalBroadcastFlag(t *testing.T) {
	hub := NewHub("")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)

	client := &Client{
		Hub:      hub,
		Send:     make(chan []byte, 10),
		UserID:   "superadmin-1",
		Role:     "superadmin",
		SchoolID: "school-1",
	}
	hub.register <- client

	// Short sleep to allow registration loop to execute
	time.Sleep(20 * time.Millisecond)

	msg := Message{
		Type:            "SYSTEM_MAINTENANCE",
		Payload:         "server update in 5 minutes",
		GlobalBroadcast: true,
	}

	hub.deliverLocally(msg)

	select {
	case received := <-client.Send:
		assert.Contains(t, string(received), "SYSTEM_MAINTENANCE")
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Expected client to receive global broadcast message")
	}
}

func TestWS_Hub_SchoolClientsIndex(t *testing.T) {
	hub := NewHub("")
	client1 := &Client{UserID: "u1", SchoolID: "school-A", Role: "teacher", Send: make(chan []byte, 10)}
	client2 := &Client{UserID: "u2", SchoolID: "school-B", Role: "teacher", Send: make(chan []byte, 10)}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	hub.register <- client1
	hub.register <- client2
	time.Sleep(20 * time.Millisecond)

	hub.mu.RLock()
	assert.Len(t, hub.schoolClients["school-A"], 1)
	assert.Len(t, hub.schoolClients["school-B"], 1)
	hub.mu.RUnlock()

	hub.unregister <- client1
	time.Sleep(20 * time.Millisecond)

	hub.mu.RLock()
	assert.Nil(t, hub.schoolClients["school-A"])
	assert.Len(t, hub.schoolClients["school-B"], 1)
	hub.mu.RUnlock()
}
