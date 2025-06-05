package tui

import (
	"fmt"
	"time"
)

func FormatLastLogin(loginAt time.Time) string {
	diff := time.Since(loginAt)

	var duration int
	var unit string

	switch {
	case diff < time.Minute:
		// If the last login was less than a minute ago
		return "just now"
	case diff < time.Hour:
		duration, unit = formatDuration(int(diff.Minutes()), "minute")
	case diff < 24*time.Hour:
		duration, unit = formatDuration(int(diff.Hours()), "hour")
	case diff < 10*24*time.Hour:
		duration, unit = formatDuration(int(diff.Hours()/24), "day")
	default:
		// Date format for long durations
		return loginAt.Format("at May 2, 2006")
	}

	// Format the output string for general durations
	return fmt.Sprintf("%d %s ago", duration, unit)
}

func formatDuration(duration int, unit string) (int, string) {
	if duration > 1 {
		unit += "s"
	}

	return duration, unit
}
