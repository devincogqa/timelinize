package utils

// Clamp returns val clamped to the range [min, max].
func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
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
	min, max := nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return min, max
}
