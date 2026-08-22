package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimit_MemoryBackendAuthBurst(t *testing.T) {
	mb := newMemoryBackend()
	ip := "192.168.1.100"

	// First auth request should be allowed (burst = 1)
	allowed1 := mb.allowAuth(nil, ip)
	assert.True(t, allowed1, "First auth request should be allowed")

	// Immediate 2nd auth request must be rate limited (burst = 1, not 5)
	allowed2 := mb.allowAuth(nil, ip)
	assert.False(t, allowed2, "Second immediate auth request must be denied (burst b=1)")
}

func TestRateLimit_IPRateLimiterEviction(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 1)
	defer limiter.Close()

	l1 := limiter.GetLimiter("10.0.0.1")
	assert.NotNil(t, l1)

	limiter.mu.RLock()
	count := len(limiter.ips)
	limiter.mu.RUnlock()
	assert.Equal(t, 1, count)
}
