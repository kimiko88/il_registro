package ws

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWS_AllowedOrigins_DynamicEvaluation(t *testing.T) {
	// Initially set allowed origins
	os.Setenv("ALLOWED_ORIGINS", "https://app1.example.com")
	origins1 := allowedOrigins()
	assert.True(t, origins1["https://app1.example.com"])
	assert.False(t, origins1["https://app2.example.com"])

	// Change env var dynamically (must take effect without sync.Once lock)
	os.Setenv("ALLOWED_ORIGINS", "https://app2.example.com")
	origins2 := allowedOrigins()
	assert.False(t, origins2["https://app1.example.com"])
	assert.True(t, origins2["https://app2.example.com"])

	os.Unsetenv("ALLOWED_ORIGINS")
}

func TestWS_Hub_DroppedMessageMetrics(t *testing.T) {
	hub := NewHub("")
	// Fill the localBroadcast channel completely
	for i := 0; i < 512; i++ {
		hub.localBroadcast <- Message{Type: "TEST"}
	}

	// Attempt sending one more message (will trigger drop logic after timeout)
	go func() {
		hub.sendBroadcast(Message{Type: "OVERFLOW"})
	}()

	time.Sleep(200 * time.Millisecond)
	assert.GreaterOrEqual(t, hub.GetDroppedMessageCount(), uint64(1))
}
