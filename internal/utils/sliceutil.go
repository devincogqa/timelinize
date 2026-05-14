package utils

// Contains checks if the given slice contains the target element.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// Unique returns a new slice with duplicate elements removed.
// BUG #2: Off-by-one — the function skips the first element, so duplicates
// of the first element are never caught and unique results may be wrong.
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for i := 1; i < len(slice); i++ {
		if !seen[slice[i]] {
			seen[slice[i]] = true
			result = append(result, slice[i])
		}
	}
	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map applies a transformation function to each element and returns the results.
func Map[T any, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

// Chunk splits a slice into chunks of the given size.
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
