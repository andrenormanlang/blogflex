package helpers

import (
	"time"
)

// FormatTime formats the time in a readable format
func FormatTime(t time.Time) string {
	// Use a more user-friendly format
	return t.Format("January 2, 2006 at 3:04 PM")
}
