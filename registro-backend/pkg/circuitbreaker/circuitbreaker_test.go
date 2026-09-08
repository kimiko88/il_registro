package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_Success(t *testing.T) {
	cb := New("test-success", Config{
		Threshold: 3,
		Timeout:   1 * time.Second,
	})

	res, err := cb.Execute(func() (interface{}, error) {
		return "hello", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "hello", res)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestCircuitBreaker_TripsAfterThreshold(t *testing.T) {
	customErr := errors.New("upstream failure")
	cb := New("test-tripping", Config{
		Threshold: 2,
		Timeout:   500 * time.Millisecond,
	})

	// First failure
	_, err := cb.Execute(func() (interface{}, error) {
		return nil, customErr
	})
	assert.ErrorIs(t, err, customErr)
	assert.Equal(t, gobreaker.StateClosed, cb.State())

	// Second failure -> Trips circuit breaker to OPEN
	_, err = cb.Execute(func() (interface{}, error) {
		return nil, customErr
	})
	assert.ErrorIs(t, err, customErr)
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// Third call -> Immediately returns ErrCircuitOpen without executing fn
	executed := false
	_, err = cb.Execute(func() (interface{}, error) {
		executed = true
		return "should not run", nil
	})
	assert.False(t, executed)
	assert.ErrorIs(t, err, ErrCircuitOpen)
}

func TestCircuitBreaker_ExecuteWithFallback(t *testing.T) {
	cb := New("test-fallback", Config{
		Threshold: 1,
		Timeout:   1 * time.Second,
	})

	// Trip circuit
	_, _ = cb.Execute(func() (interface{}, error) {
		return nil, errors.New("service down")
	})

	// Run with fallback
	res, err := cb.ExecuteWithFallback(
		func() (interface{}, error) {
			return "main", nil
		},
		func(err error) (interface{}, error) {
			return "cached_fallback", nil
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, "cached_fallback", res)
}

func TestCircuitBreaker_GlobalSingletons(t *testing.T) {
	s1 := Storage()
	s2 := Storage()
	assert.Equal(t, s1, s2)
	assert.Equal(t, "supabase-storage", s1.Name())

	sidi := Sidi()
	assert.Equal(t, "sidi-mim", sidi.Name())

	wh := Webhook()
	assert.Equal(t, "webhooks", wh.Name())
}
