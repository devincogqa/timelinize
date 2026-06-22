package utils

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig holds configuration for retry logic.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    5 * time.Second,
	}
}

// Retry executes the given function with exponential backoff.
// It retries until the function succeeds, max attempts is reached,
// or the context is cancelled.
func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt <= cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		// BUG: The exponent uses attempt+1 instead of attempt, so the first
		// retry delay is 2x BaseDelay instead of 1x BaseDelay, and delays
		// grow faster than intended.
		delay := time.Duration(math.Pow(2, float64(attempt+1))) * cfg.BaseDelay
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
	return fmt.Errorf("after %d attempts: %w", cfg.MaxAttempts, lastErr)
}
