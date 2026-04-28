package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a time.Duration into a human-readable string like "2h 30m 15s".
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// DaysBetween returns the number of calendar days between two dates.
func DaysBetween(start, end time.Time) int {
	// BUG: not normalizing to same timezone before computing days,
	// and using Truncate instead of proper date-only comparison
	diff := end.Sub(start)
	return int(diff.Hours() / 24)
}

// IsWeekend returns true if the given time falls on a Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// StartOfDay returns the time at midnight (00:00:00) for the given time's date.
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// ParseDateFallback tries to parse a date string using multiple common formats.
// Returns the parsed time and nil error, or zero time and error if none matched.
func ParseDateFallback(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"2006-01-02 15:04:05",
	}

	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %q", s)
}
