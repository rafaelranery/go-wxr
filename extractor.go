package wxr

import (
	"strconv"
	"strings"
)

// AuthorExtractor handles author name resolution from WXR items.
type AuthorExtractor struct{}

// Extract extracts the author name from an item.
// It prefers explicit meta fields (redator, autor, author_name), then falls back to dc:creator.
func (e *AuthorExtractor) Extract(item *item) string {
	author := getMetaValue(item.PostMeta, "redator", "autor", "author_name")
	if author == "" {
		author = item.DCCreator
	}
	return strings.TrimSpace(author)
}

// ExcerptExtractor handles excerpt extraction from WXR items.
type ExcerptExtractor struct{}

// Extract extracts the excerpt from an item.
// It prefers the excerpt:encoded field, then falls back to the "subtitulo" meta field.
func (e *ExcerptExtractor) Extract(item *item) string {
	excerpt := item.ExcerptEncoded
	if strings.TrimSpace(excerpt) == "" {
		excerpt = getMetaValue(item.PostMeta, "subtitulo")
	}
	return excerpt
}

// FeaturedImageExtractor handles featured image URL resolution from WXR items.
type FeaturedImageExtractor struct {
	attachmentIndex *AttachmentIndex
}

// Extract extracts the featured image URL from an item.
// It tries multiple sources in order:
// 1. Custom meta fields (banner_da_materia, banner_old, link_do_banner)
// 2. WordPress thumbnail ID (_thumbnail_id) pointing to an attachment
// 3. First attachment associated with the post
func (e *FeaturedImageExtractor) Extract(item *item) string {
	// Try custom meta fields first
	featuredImage := getMetaValue(item.PostMeta, "banner_da_materia", "banner_old", "link_do_banner")
	if featuredImage != "" {
		return featuredImage
	}

	// Try thumbnail ID
	if thumbIDStr := getMetaValue(item.PostMeta, "_thumbnail_id"); thumbIDStr != "" {
		if thumbID, err := strconv.Atoi(thumbIDStr); err == nil {
			if url, ok := e.attachmentIndex.URLsByID[thumbID]; ok {
				return url
			}
		}
	}

	// Try first attachment associated with post
	if urls := e.attachmentIndex.URLsByParent[item.PostID]; len(urls) > 0 {
		return urls[0]
	}

	return ""
}

// DateExtractor handles date extraction and normalization from WXR items.
type DateExtractor struct{}

// Extract extracts and normalizes the date from an item.
// It prefers wp:post_date_gmt, then wp:post_date, then pubDate.
// The date is normalized to RFC3339 format.
func (e *DateExtractor) Extract(item *item) string {
	dateStr := item.PostDateGMT
	if dateStr == "" {
		dateStr = item.PostDate
	}
	if dateStr == "" {
		dateStr = item.PubDate
	}

	// Normalize date format to RFC3339 if possible
	if dateStr != "" {
		dateStr = normalizeWXRDate(dateStr)
	}

	return dateStr
}
