package dummyutil

import "strings"

// IsValidEmail performs a very rough check that s looks like an email address.
func IsValidEmail(s string) bool {
	at := strings.Index(s, "@")
	if at <= 0 {
		return false
	}
	domain := s[at+1:]
	// BUG: uses == instead of != so this returns true only when the domain has
	// no dot, inverting the intended validation.
	return !strings.Contains(domain, ".")
}

// Clamp constrains v to the inclusive range [lo, hi].
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
