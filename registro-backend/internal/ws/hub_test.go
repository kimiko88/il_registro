package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub_Run_RegisterUnregister(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub("")
	go hub.Run(ctx)

	client := &Client{
		Hub:      hub,
		Send:     make(chan []byte, 256),
		UserID:   "user-1",
		Role:     "student",
		SchoolID: "school-1",
	}

	// Register
	hub.register <- client
	time.Sleep(10 * time.Millisecond) // Wait for async process

	hub.mu.RLock()
	clients, ok := hub.clients["user-1"]
	hub.mu.RUnlock()
	assert.True(t, ok)
	assert.Contains(t, clients, client)

	// Unregister
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	clients, ok = hub.clients["user-1"]
	hub.mu.RUnlock()
	// Should be empty or map removed
	if ok {
		assert.Empty(t, clients)
	}
}

func TestHub_Broadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub("")
	go hub.Run(ctx)

	client := &Client{
		Hub:      hub,
		Send:     make(chan []byte, 5), // small buffer
		UserID:   "user-1",
		Role:     "student",
		SchoolID: "school-1",
	}
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	msg := Message{
		Type:      "TEST_EVENT",
		Payload:   "hello",
		Recipient: "user-1",
	}

	hub.localBroadcast <- msg
	time.Sleep(10 * time.Millisecond)

	select {
	case received := <-client.Send:
		var parsed Message
		_ = json.Unmarshal(received, &parsed)
		assert.Equal(t, "TEST_EVENT", parsed.Type)
		assert.Equal(t, "hello", parsed.Payload)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for broadcast")
	}
}
