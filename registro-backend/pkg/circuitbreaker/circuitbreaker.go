package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open: external service temporarily unavailable")
)

// CircuitBreaker wraps sony/gobreaker with convenient execution helpers
type CircuitBreaker struct {
	cb   *gobreaker.CircuitBreaker
	name string
}

type Config struct {
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
	Threshold   uint32
}

var (
	defaultStorageBreaker *CircuitBreaker
	defaultSidiBreaker    *CircuitBreaker
	defaultWebhookBreaker *CircuitBreaker
	once                  sync.Once
)

// New creates a new CircuitBreaker with custom parameters
func New(name string, cfg Config) *CircuitBreaker {
	if cfg.MaxRequests == 0 {
		cfg.MaxRequests = 3
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}
	if cfg.Threshold == 0 {
		cfg.Threshold = 5
	}

	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.Threshold
		},
	}

	return &CircuitBreaker{
		cb:   gobreaker.NewCircuitBreaker(st),
		name: name,
	}
}

// Storage returns the global circuit breaker for external storage (Supabase / S3)
func Storage() *CircuitBreaker {
	initDefaults()
	return defaultStorageBreaker
}

// Sidi returns the global circuit breaker for MIM SIDI flussi
func Sidi() *CircuitBreaker {
	initDefaults()
	return defaultSidiBreaker
}

// Webhook returns the global circuit breaker for external webhooks / push notifications
func Webhook() *CircuitBreaker {
	initDefaults()
	return defaultWebhookBreaker
}

func initDefaults() {
	once.Do(func() {
		defaultStorageBreaker = New("supabase-storage", Config{
			MaxRequests: 2,
			Timeout:     20 * time.Second,
			Threshold:   4,
		})
		defaultSidiBreaker = New("sidi-mim", Config{
			MaxRequests: 1,
			Timeout:     30 * time.Second,
			Threshold:   3,
		})
		defaultWebhookBreaker = New("webhooks", Config{
			MaxRequests: 3,
			Timeout:     15 * time.Second,
			Threshold:   5,
		})
	})
}

// Execute runs the given operation protected by the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	res, err := cb.cb.Execute(func() (interface{}, error) {
		return fn()
	})
	if errors.Is(err, gobreaker.ErrOpenState) {
		return nil, ErrCircuitOpen
	}
	return res, err
}

// ExecuteWithFallback runs operation with a fallback function if the circuit is open or errors occur
func (cb *CircuitBreaker) ExecuteWithFallback(
	fn func() (interface{}, error),
	fallback func(err error) (interface{}, error),
) (interface{}, error) {
	res, err := cb.Execute(fn)
	if err != nil {
		if fallback != nil {
			return fallback(err)
		}
		return nil, err
	}
	return res, nil
}

// State returns the current circuit breaker state
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.cb.State()
}

// Counts returns the current execution counts
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.cb.Counts()
}

// Name returns the breaker identifier
func (cb *CircuitBreaker) Name() string {
	return cb.name
}
