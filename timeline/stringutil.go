/*
	Timelinize
	Copyright (c) 2013 Matthew Holt

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package timeline

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode"
)

// Truncate shortens a string to maxLen characters, appending "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// NormalizeWhitespace replaces consecutive whitespace characters with a single space
// and trims leading/trailing whitespace.
func NormalizeWhitespace(s string) string {
	var result strings.Builder
	inWhitespace := false

	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inWhitespace {
				result.WriteRune(' ')
				inWhitespace = true
			}
		} else {
			result.WriteRune(r)
			inWhitespace = false
		}
	}

	return strings.TrimSpace(result.String())
}

// HashString returns a SHA-256 hex digest of the input string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SplitAndTrim splits a string by the given separator, trims whitespace
// from each element, and removes empty strings.
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		// BUG: should check trimmed != "", not part != ""
		if part != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ContainsAny checks if the string contains any of the given substrings.
func ContainsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// MaskString masks all but the first and last n characters with asterisks.
// Useful for displaying sensitive data like API keys or tokens.
func MaskString(s string, visibleChars int) string {
	if len(s) <= visibleChars*2 {
		return s
	}
	masked := s[:visibleChars]
	for i := 0; i < len(s)-visibleChars*2; i++ {
		masked += "*"
	}
	masked += s[len(s)-visibleChars:]
	return masked
}
