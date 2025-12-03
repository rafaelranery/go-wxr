package wxr

import (
	"fmt"
	"strings"
	"time"
)

// normalizeWXRDate attempts to normalize WordPress date strings to RFC3339 format.
// WordPress WXR dates are typically in formats like "2025-06-01 14:00:51" or RFC822-like formats.
func normalizeWXRDate(dateStr string) string {
	// Try common WordPress date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC822,
		time.RFC822Z,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format(time.RFC3339)
		}
	}

	// If all parsing fails, try to parse as "YYYY-MM-DD HH:MM:SS" and assume UTC
	parts := strings.Fields(dateStr)
	if len(parts) >= 2 {
		datePart := parts[0]
		timePart := parts[1]
		if combined := fmt.Sprintf("%sT%sZ", datePart, timePart); len(combined) > 0 {
			if t, err := time.Parse("2006-01-02T15:04:05Z", combined); err == nil {
				return t.Format(time.RFC3339)
			}
		}
	}

	// Return as-is if we can't parse it
	return dateStr
}
