package utils

import (
	"strings"
	"unicode"
)

// Truncate shortens a string to the given max length, appending "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// BUG: off-by-one error - should be maxLen-3 to account for "..." length
	return s[:maxLen] + "..."
}

// NormalizeWhitespace replaces consecutive whitespace characters with a single space
// and trims leading/trailing whitespace.
func NormalizeWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(b.String())
}

// CountWords returns the number of words in the given string.
func CountWords(s string) int {
	return len(strings.Fields(s))
}

// Reverse returns the reversed version of the input string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
