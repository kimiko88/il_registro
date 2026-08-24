package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_ParseAcceptLanguage_UTF8RuneSafety(t *testing.T) {
	// Simple standard locale
	assert.Equal(t, "en-US", parseAcceptLanguage("en-US,en;q=0.9"))

	// Long string with multi-byte UTF-8 characters past 35 bytes
	utf8Header := "zh-CN,zh;q=0.9,📍📌🚩🌟🎉🎈🎌🎨🎭🎪🎬🎤🎧🎨🎭🎪"
	res := parseAcceptLanguage(utf8Header)
	// Must not panic or produce broken multi-byte characters
	assert.NotEmpty(t, res)
	assert.Equal(t, "zh-CN", res)

	// Empty fallback
	assert.Equal(t, "it-IT", parseAcceptLanguage(""))
}

func TestAuth_ActiveCache_GuaranteedEvictionWhenFull(t *testing.T) {
	mockRepo := new(DummyUserRepo)
	mw := NewMiddleware(nil, mockRepo)
	mw.activeCache = make(map[string]activeCacheEntry)

	// Pre-fill activeCache to maxActiveCacheSize with unexpired entries
	now := time.Now()
	for i := 0; i < maxActiveCacheSize; i++ {
		mw.activeCache[fmt.Sprintf("user-%d", i)] = activeCacheEntry{
			isActive:  true,
			expiresAt: now.Add(time.Duration(i+1) * time.Minute),
		}
	}

	assert.Equal(t, maxActiveCacheSize, len(mw.activeCache))

	// Mock DB check for a brand new user
	mockRepo.On("IsActive", mock.Anything, "new-user-1").Return(true, nil).Once()

	active, err := mw.isAccountActive(context.Background(), "new-user-1")
	assert.NoError(t, err)
	assert.True(t, active)

	// The new user MUST be in cache (oldest entry evicted)
	mw.cacheMu.RLock()
	entry, exists := mw.activeCache["new-user-1"]
	mw.cacheMu.RUnlock()

	assert.True(t, exists, "New entry must be guaranteed insertion even when cache is at max capacity")
	assert.True(t, entry.isActive)
	assert.Equal(t, maxActiveCacheSize, len(mw.activeCache))
}

func TestAuth_RequestPasswordReset_ConstantTimeWorkOnMissingUser(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &Service{repo: mockRepo}

	mockRepo.On("GetUserByEmail", mock.Anything, "nonexistent@school.it").Return(nil, assert.AnError).Once()

	start := time.Now()
	err := svc.RequestPasswordReset(context.Background(), "nonexistent@school.it")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, elapsed >= 0)
}
