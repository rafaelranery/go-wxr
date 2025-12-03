package wxr

import "strings"

// CategoryExtractor handles category and tag extraction from WXR items.
type CategoryExtractor struct{}

// ExtractCategories extracts category names from an item.
// Categories are identified by domain="category" or no domain attribute.
func (e *CategoryExtractor) ExtractCategories(item *item) []string {
	var categories []string
	for _, cat := range item.Categories {
		domain := strings.ToLower(strings.TrimSpace(cat.Domain))
		// Categories have domain="category" or no domain (defaults to category)
		if domain == "category" || domain == "" {
			if name := strings.TrimSpace(cat.Value); name != "" {
				categories = append(categories, name)
			}
		}
	}
	return categories
}

// ExtractTags extracts tag names from an item.
// Tags are identified by domain="post_tag".
func (e *CategoryExtractor) ExtractTags(item *item) []string {
	var tags []string
	for _, cat := range item.Categories {
		domain := strings.ToLower(strings.TrimSpace(cat.Domain))
		if domain == "post_tag" {
			if name := strings.TrimSpace(cat.Value); name != "" {
				tags = append(tags, name)
			}
		}
	}
	return tags
}

