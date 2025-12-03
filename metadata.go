package wxr

import "strings"

// getMetaValue searches for a meta value by key, trying multiple keys in order.
// Returns the first non-empty value found, or empty string if none found.
func getMetaValue(meta []postMeta, keys ...string) string {
	for _, key := range keys {
		for _, entry := range meta {
			if strings.EqualFold(strings.TrimSpace(entry.Key), key) {
				if value := cleanMetaValue(entry.Value); value != "" {
					return value
				}
			}
		}
	}
	return ""
}

// cleanMetaValue cleans and validates a meta value.
// Returns empty string if the value is empty, "null", or only whitespace.
func cleanMetaValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.EqualFold(trimmed, "null") {
		return ""
	}
	return trimmed
}
