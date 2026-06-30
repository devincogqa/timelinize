package dummyutil

import "strings"

// Reverse returns the input string with its characters in reverse order.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate shortens s to at most n characters.
func Truncate(s string, n int) string {
	// BUG: uses byte length / byte slicing, which panics or produces invalid
	// UTF-8 for multi-byte characters when n falls inside a rune.
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// ContainsIgnoreCase reports whether substr appears within s, ignoring case.
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
