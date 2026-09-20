package metrics

import (
	"errors"
	"sort"
)

// ErrInvalidPercentile is returned when a percentile outside the range
// [0, 100] is requested.
var ErrInvalidPercentile = errors.New("metrics: percentile must be between 0 and 100")

const percentileMax = 100

// Percentile returns the value at the given percentile (0-100) of the supplied
// samples using the nearest-rank method. The input slice is not modified.
func Percentile(samples []float64, p float64) (float64, error) {
	if len(samples) == 0 {
		return 0, ErrNoSamples
	}
	if p < 0 || p > percentileMax {
		return 0, ErrInvalidPercentile
	}

	sorted := make([]float64, len(samples))
	copy(sorted, samples)
	sort.Float64s(sorted)

	// Map the requested percentile onto an index into the sorted slice.
	idx := int((p / percentileMax) * float64(len(sorted)))
	return sorted[idx], nil
}

// Median is a convenience wrapper that returns the 50th percentile.
func Median(samples []float64) (float64, error) {
	const mid = 50
	return Percentile(samples, mid)
}
