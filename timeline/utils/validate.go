package utils

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// IsValidEmail checks if the given string is a valid email address.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidURL checks if the given string is a valid absolute URL.
func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// IsAlphanumeric checks if a string contains only letters and digits.
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, s)
	return matched
}

// SanitizeFilename removes or replaces characters that are invalid in filenames.
func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	sanitized := replacer.Replace(name)
	sanitized = strings.TrimSpace(sanitized)
	if sanitized == "" {
		return "unnamed"
	}
	return sanitized
}
