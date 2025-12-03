package wxr

// Post represents a WordPress post parsed from a WXR export file.
// All fields are normalized and ready for use in applications.
type Post struct {
	// ID is the WordPress post ID.
	ID int

	// TitleRendered is the post title (HTML may be present).
	TitleRendered string

	// ContentRendered is the full post content (HTML).
	ContentRendered string

	// Slug is the URL-friendly post slug.
	Slug string

	// Link is the canonical permalink URL for the post.
	Link string

	// Excerpt is the post excerpt or summary.
	Excerpt string

	// Author is the post author name.
	Author string

	// Categories is a list of category names or IDs associated with the post.
	Categories []string

	// Date is the post publication date in RFC3339 format.
	Date string

	// FeaturedImage is the URL of the featured image for the post.
	FeaturedImage string
}
