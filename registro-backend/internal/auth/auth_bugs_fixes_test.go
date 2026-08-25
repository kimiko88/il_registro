package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_Middleware_ActiveCacheCap(t *testing.T) {
	mw := &Middleware{
		activeCache: make(map[string]activeCacheEntry),
	}

	// Populate 10,000 entries
	mw.cacheMu.Lock()
	for i := 0; i < maxActiveCacheSize; i++ {
		mw.activeCache[string(rune(i))] = activeCacheEntry{
			isActive:  true,
			expiresAt: time.Now().Add(-1 * time.Minute), // expired
		}
	}
	mw.cacheMu.Unlock()

	// Invalidate & Cleanup
	mw.CleanupExpiredEntries()

	mw.cacheMu.RLock()
	size := len(mw.activeCache)
	mw.cacheMu.RUnlock()

	assert.Equal(t, 0, size, "CleanupExpiredEntries should evict all expired entries")
}

func TestAuth_RecordFailedAttempt_Normalization(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil)

	mockRepo.On("RecordLoginAttempt", mock.Anything, mock.MatchedBy(func(a *LoginAttempt) bool {
		return a.Email == "user@test.com" && !a.Success
	})).Return(nil).Once()

	svc.RecordFailedAttempt(context.Background(), "  USER@TEST.COM ", "127.0.0.1")
	mockRepo.AssertExpectations(t)
}
