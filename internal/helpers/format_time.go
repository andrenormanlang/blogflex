package helpers

import (
    "time"
)

// FormatTime formats the time in a readable format
func FormatTime(timeStr string) string {
    t, err := time.Parse(time.RFC850, timeStr)
    if err != nil {
        return timeStr // return the original string if parsing fails
    }
    return t.Format("Jan 2, 2006 at 3:04pm")
}
