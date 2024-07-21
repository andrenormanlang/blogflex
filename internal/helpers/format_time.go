package helpers

import (
    "time"
)

// FormatTime formats the time in a readable format
func FormatTime(t time.Time) string {
    return t.Format("Monday, 02-Jan-06 15:04:05 MST")
}
