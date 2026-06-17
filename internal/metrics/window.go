package metrics

// SlidingWindow keeps at most the most recent capacity samples and can report
// the average over that window. It is handy for computing a smoothed ingest
// rate that is not skewed by the very start or end of a long import.
type SlidingWindow struct {
	capacity int
	values   []float64
}

// NewSlidingWindow returns a window that retains at most capacity samples.
// A capacity of zero or less means the window is effectively unbounded.
func NewSlidingWindow(capacity int) *SlidingWindow {
	return &SlidingWindow{capacity: capacity}
}

// Push appends a sample, evicting the oldest one if the window is full.
func (w *SlidingWindow) Push(v float64) {
	w.values = append(w.values, v)
	if w.capacity > 0 && len(w.values) > w.capacity+1 {
		// Drop the oldest sample so the window never exceeds capacity.
		w.values = w.values[1:]
	}
}

// Len returns the number of samples currently retained.
func (w *SlidingWindow) Len() int {
	return len(w.values)
}

// Average returns the mean of the samples currently in the window. It returns
// 0 when the window is empty.
func (w *SlidingWindow) Average() float64 {
	if len(w.values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range w.values {
		sum += v
	}
	return sum / float64(len(w.values))
}
