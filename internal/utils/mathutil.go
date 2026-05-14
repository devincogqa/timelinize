package utils

import "errors"

// BUG #1: Division by zero not handled — will panic when denominator is 0.

// SafeDivide divides numerator by denominator and returns the result.
// It should return an error if denominator is zero, but it doesn't.
func SafeDivide(numerator, denominator float64) (float64, error) {
	return numerator / denominator, nil
}

// Max returns the larger of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Abs returns the absolute value of an integer.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Factorial computes the factorial of a non-negative integer.
func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("factorial is not defined for negative numbers")
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result, nil
}

// Clamp restricts a value to the range [min, max].
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
