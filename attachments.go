package wxr

import (
	"fmt"
	"strings"
)

// AttachmentIndex holds attachment URL mappings for resolving featured images.
type AttachmentIndex struct {
	URLsByID     map[int]string
	URLsByParent map[int][]string
}

// buildAttachmentIndex builds an index of attachments from WXR items.
// It creates mappings from attachment IDs to URLs and from parent post IDs to attachment URLs.
func buildAttachmentIndex(ch channel) *AttachmentIndex {
	index := &AttachmentIndex{
		URLsByID:     make(map[int]string),
		URLsByParent: make(map[int][]string),
	}

	baseUploadsURL := determineUploadsBaseURL(ch)

	for _, item := range ch.Items {
		if item.PostType != "attachment" {
			continue
		}

		url := resolveAttachmentURL(item, baseUploadsURL)
		if url != "" {
			index.URLsByID[item.PostID] = url
			if item.PostParent > 0 {
				index.URLsByParent[item.PostParent] = append(index.URLsByParent[item.PostParent], url)
			}
		}
	}

	return index
}

// resolveAttachmentURL resolves the URL for an attachment item.
// It tries multiple sources in order: AttachmentURL, Link, GUID, or constructs from _wp_attached_file meta.
func resolveAttachmentURL(item item, baseUploadsURL string) string {
	url := strings.TrimSpace(item.AttachmentURL)
	if url != "" {
		return url
	}

	url = strings.TrimSpace(item.Link)
	if url != "" {
		return url
	}

	url = strings.TrimSpace(item.GUID)
	if url != "" {
		return url
	}

	// Try to construct URL from _wp_attached_file meta
	if rel := getMetaValue(item.PostMeta, "_wp_attached_file"); rel != "" && baseUploadsURL != "" {
		rel = strings.TrimPrefix(rel, "/")
		return fmt.Sprintf("%s/wp-content/uploads/%s", baseUploadsURL, rel)
	}

	return ""
}

// determineUploadsBaseURL determines the base URL for WordPress uploads.
// It tries BaseSiteURL, BaseBlogURL, and Link in order.
func determineUploadsBaseURL(ch channel) string {
	candidates := []string{ch.BaseSiteURL, ch.BaseBlogURL, ch.Link}
	for _, candidate := range candidates {
		clean := strings.TrimSpace(candidate)
		if clean != "" {
			return strings.TrimRight(clean, "/")
		}
	}
	return ""
}
