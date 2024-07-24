package helpers

import (
    "strings"
)

func TruncateWords(content string, limit int) string {
    words := strings.Fields(content)
    if len(words) > limit {
        return strings.Join(words[:limit], " ") + "..."
    }
    return strings.Join(words, " ")
}
