package dummyutil

import "errors"

// ErrDivideByZero is returned when a division by zero is attempted.
var ErrDivideByZero = errors.New("division by zero")

// Divide returns a divided by b.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

// Factorial returns n! for non-negative n.
func Factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// Percent returns what percentage part is of whole.
func Percent(part, whole float64) float64 {
	return part / whole * 100
}
