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
		return "just now"
	case diff < time.Hour:
		duration = int(diff.Minutes())
		unit = "minute"
		if duration > 1 {
			unit += "s"
		}
	case diff < 24*time.Hour:
		duration := int(diff.Hours())
		unit = "hour"
		if duration > 1 {
			unit += "s"
		}
	case diff < 10*24*time.Hour:
		days := int(diff.Hours() / 24)
		unit = "day"
		if days > 1 {
			unit += "s"
		}
	default:
		return loginAt.Format("2011-01-02 13:04:25")
	}

	return fmt.Sprintf("%d %s ago", duration, unit)
}
