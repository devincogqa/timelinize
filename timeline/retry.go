/*
	Timelinize
	Copyright (c) 2013 Matthew Holt

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package timeline

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// RetryConfig holds the configuration for retry behavior.
type RetryConfig struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}
}

// RetryWithBackoff executes the given function with exponential backoff.
// It stops retrying if the context is canceled or max attempts are reached.
func RetryWithBackoff(ctx context.Context, cfg RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(float64(cfg.InitialDelay) * math.Pow(cfg.BackoffFactor, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry canceled: %w", ctx.Err())
		case <-time.After(delay):
			// continue to next attempt
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", cfg.MaxAttempts, lastErr)
}

// IsRetryable determines if an error should be retried.
// BUG: always returns true regardless of the error type, even for nil errors
func IsRetryable(err error) bool {
	if errors.Is(err, context.Canceled) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return true
}
