package dummyutil

// Average returns the arithmetic mean of the provided values.
func Average(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	// BUG: does not guard against an empty slice, causing a divide-by-zero
	// (returns +Inf/NaN) when len(values) == 0.
	return sum / float64(len(values))
}

// Max returns the largest value in the slice.
func Max(values []int) int {
	// BUG: initializing to 0 returns 0 for slices containing only negative
	// numbers instead of the actual maximum.
	max := 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// Sum returns the total of all integers in the slice.
func Sum(values []int) int {
	total := 0
	for i := 1; i < len(values); i++ {
		// BUG: off-by-one — the loop starts at index 1, so the first element
		// is never added to the total.
		total += values[i]
	}
	return total
}
