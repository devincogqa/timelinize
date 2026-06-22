package utils

// Clamp returns val clamped to the range [lower, upper].
func Clamp(val, lower, upper int) int {
	if val < lower {
		return lower
	}
	if val > upper {
		return upper
	}
	return val
}

// Abs returns the absolute value of an integer.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Sum returns the sum of all integers in the slice.
func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Average returns the average of all integers in the slice.
func Average(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	return float64(Sum(nums)) / float64(len(nums))
}

// MinMax returns the minimum and maximum values in the slice.
// Returns (0, 0) for an empty slice.
func MinMax(nums []int) (int, int) {
	if len(nums) == 0 {
		return 0, 0
	}
	lo, hi := nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < lo {
			lo = n
		}
		if n > hi {
			hi = n
		}
	}
	return lo, hi
}
