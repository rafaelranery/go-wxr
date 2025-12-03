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

	// Categories is a list of category names associated with the post.
	Categories []string

	// Tags is a list of tag names associated with the post.
	Tags []string

	// Date is the post publication date in RFC3339 format.
	Date string

	// ModifiedDate is the post last modification date in RFC3339 format.
	ModifiedDate string

	// FeaturedImage is the URL of the featured image for the post.
	FeaturedImage string

	// GUID is the globally unique identifier for the post.
	GUID string

	// ParentID is the ID of the parent post (for hierarchical post types like pages).
	ParentID int

	// Meta contains all post meta fields as key-value pairs.
	Meta map[string]string
}
