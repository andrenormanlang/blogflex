package helpers

import (
    "time"
)

// FormatTime formats the time in a readable format
func FormatTime(t time.Time) string {
    return t.Format("Jan 2, 2006 at 3:04pm")
}
