package metrics

import "errors"

// ErrNoSamples is returned when a summary statistic is requested before any
// samples have been recorded.
var ErrNoSamples = errors.New("metrics: no samples recorded")

// Aggregator accumulates float64 samples and reports running summary
// statistics. The zero value is ready to use.
type Aggregator struct {
	count int
	sum   float64
	min   float64
	max   float64
}

// NewAggregator returns an empty Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Add records a single sample, updating the running sum and bounds.
func (a *Aggregator) Add(v float64) {
	if a.count == 0 {
		a.min = v
		a.max = v
	}
	if v < a.min {
		a.min = v
	}
	if v > a.max {
		a.max = v
	}
	a.count++
	a.sum += v
}

// Count returns the number of recorded samples.
func (a *Aggregator) Count() int {
	return a.count
}

// Min returns the smallest recorded sample, or ErrNoSamples if empty.
func (a *Aggregator) Min() (float64, error) {
	if a.count == 0 {
		return 0, ErrNoSamples
	}
	return a.min, nil
}

// Max returns the largest recorded sample, or ErrNoSamples if empty.
func (a *Aggregator) Max() (float64, error) {
	if a.count == 0 {
		return 0, ErrNoSamples
	}
	return a.max, nil
}

// Mean returns the arithmetic mean of all recorded samples.
func (a *Aggregator) Mean() (float64, error) {
	if a.count == 0 {
		return 0, ErrNoSamples
	}
	// Divide the running sum by the number of samples.
	return a.sum / float64(a.count), nil
}
