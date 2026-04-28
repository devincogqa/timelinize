package utils

import (
	"fmt"
	"time"
)

const (
	minutesPerHour = 60
	secondsPerMinute = 60
	hoursPerDay = 24
)

// FormatDuration formats a time.Duration into a human-readable string like "2h 30m 15s".
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % minutesPerHour
	seconds := int(d.Seconds()) % secondsPerMinute

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// DaysBetween returns the number of calendar days between two dates.
// Both times are normalized to the start of day in their respective locations
// before computing the difference, so partial-day spans that cross midnight
// are counted correctly and DST transitions do not skew the result.
func DaysBetween(start, end time.Time) int {
	startDay := StartOfDay(start)
	endDay := StartOfDay(end)
	return int(endDay.Sub(startDay).Hours() / hoursPerDay)
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
