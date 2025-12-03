package wxr

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

// Parser provides configurable parsing of WordPress WXR export files.
type Parser struct {
	logger            Logger
	filter            Filter
	authorExt         *AuthorExtractor
	excerptExt        *ExcerptExtractor
	dateExt           *DateExtractor
	modifiedDateExt   *ModifiedDateExtractor
	categoryExt       *CategoryExtractor
	metaExt           *MetaExtractor
	featuredImageExt  *FeaturedImageExtractor
}

// NewParser creates a new Parser with the default no-op logger.
func NewParser() *Parser {
	return &Parser{
		logger:          &noOpLogger{},
		filter:          NewDefaultFilter(),
		authorExt:       &AuthorExtractor{},
		excerptExt:      &ExcerptExtractor{},
		dateExt:         &DateExtractor{},
		modifiedDateExt: &ModifiedDateExtractor{},
		categoryExt:     &CategoryExtractor{},
		metaExt:         &MetaExtractor{},
	}
}

// NewParserWithLogger creates a new Parser with a custom logger.
func NewParserWithLogger(logger Logger) *Parser {
	return &Parser{
		logger:          logger,
		filter:          NewDefaultFilter(),
		authorExt:       &AuthorExtractor{},
		excerptExt:      &ExcerptExtractor{},
		dateExt:         &DateExtractor{},
		modifiedDateExt: &ModifiedDateExtractor{},
		categoryExt:     &CategoryExtractor{},
		metaExt:         &MetaExtractor{},
	}
}

// NewParserWithStdLogger creates a new Parser using the standard log.Logger.
// This is a convenience function for users who want to use the standard library logger.
func NewParserWithStdLogger(stdLogger *log.Logger) *Parser {
	if stdLogger == nil {
		return NewParser()
	}
	return &Parser{
		logger:          &stdLoggerAdapter{logger: stdLogger},
		filter:          NewDefaultFilter(),
		authorExt:       &AuthorExtractor{},
		excerptExt:      &ExcerptExtractor{},
		dateExt:         &DateExtractor{},
		modifiedDateExt: &ModifiedDateExtractor{},
		categoryExt:     &CategoryExtractor{},
		metaExt:         &MetaExtractor{},
	}
}

// SetLogger sets the logger for the parser.
func (p *Parser) SetLogger(logger Logger) {
	if logger == nil {
		p.logger = &noOpLogger{}
	} else {
		p.logger = logger
	}
}

// WithFilter sets a custom filter for the parser.
// Returns the parser for method chaining.
func (p *Parser) WithFilter(filter Filter) *Parser {
	if filter == nil {
		p.filter = NewDefaultFilter()
	} else {
		p.filter = filter
	}
	return p
}

// WithAuthorExtractor sets a custom author extractor for the parser.
// Returns the parser for method chaining.
func (p *Parser) WithAuthorExtractor(extractor *AuthorExtractor) *Parser {
	if extractor == nil {
		p.authorExt = &AuthorExtractor{}
	} else {
		p.authorExt = extractor
	}
	return p
}

// WithExcerptExtractor sets a custom excerpt extractor for the parser.
// Returns the parser for method chaining.
func (p *Parser) WithExcerptExtractor(extractor *ExcerptExtractor) *Parser {
	if extractor == nil {
		p.excerptExt = &ExcerptExtractor{}
	} else {
		p.excerptExt = extractor
	}
	return p
}

// WithDateExtractor sets a custom date extractor for the parser.
// Returns the parser for method chaining.
func (p *Parser) WithDateExtractor(extractor *DateExtractor) *Parser {
	if extractor == nil {
		p.dateExt = &DateExtractor{}
	} else {
		p.dateExt = extractor
	}
	return p
}

// WithFeaturedImageExtractor sets a custom featured image extractor for the parser.
// Note: The attachment index will be set automatically during parsing.
// Returns the parser for method chaining.
func (p *Parser) WithFeaturedImageExtractor(extractor *FeaturedImageExtractor) *Parser {
	if extractor == nil {
		p.featuredImageExt = &FeaturedImageExtractor{}
	} else {
		p.featuredImageExt = extractor
	}
	return p
}

// decodeXML decodes and validates the WXR XML document.
func (p *Parser) decodeXML(r io.Reader) (*wxr, error) {
	var wxrDoc wxr
	decoder := xml.NewDecoder(r)

	// Handle CDATA sections properly
	decoder.Strict = false

	if err := decoder.Decode(&wxrDoc); err != nil {
		return nil, fmt.Errorf("wxr: failed to parse WXR XML: %w", err)
	}

	// Validate that we actually got an RSS document
	// XMLName.Local will be empty if the decoder didn't find a matching root element
	if wxrDoc.XMLName.Local != "rss" {
		return nil, fmt.Errorf("wxr: invalid WXR XML: root element is not <rss> (got %q)", wxrDoc.XMLName.Local)
	}

	return &wxrDoc, nil
}

// buildAuthorMap builds a map from author login to display name.
func (p *Parser) buildAuthorMap(ch channel) map[string]string {
	authorMap := make(map[string]string)
	for _, author := range ch.Authors {
		authorMap[author.Login] = author.DisplayName
	}
	return authorMap
}

// transformItem converts an item to a Post using the configured extractors.
func (p *Parser) transformItem(item *item, attachmentIndex *AttachmentIndex) Post {
	// Initialize extractors if needed
	if p.categoryExt == nil {
		p.categoryExt = &CategoryExtractor{}
	}
	if p.modifiedDateExt == nil {
		p.modifiedDateExt = &ModifiedDateExtractor{}
	}
	if p.metaExt == nil {
		p.metaExt = &MetaExtractor{}
	}
	if p.featuredImageExt == nil {
		p.featuredImageExt = &FeaturedImageExtractor{attachmentIndex: attachmentIndex}
	} else {
		p.featuredImageExt.attachmentIndex = attachmentIndex
	}

	categories := p.categoryExt.ExtractCategories(item)
	if categories == nil {
		categories = []string{}
	}
	tags := p.categoryExt.ExtractTags(item)
	if tags == nil {
		tags = []string{}
	}
	meta := p.metaExt.Extract(item)
	if meta == nil {
		meta = make(map[string]string)
	}

	return Post{
		ID:              item.PostID,
		TitleRendered:   item.Title,
		ContentRendered: item.ContentEncoded,
		Excerpt:         p.excerptExt.Extract(item),
		Slug:            item.PostName,
		Link:            item.Link, // Canonical permalink from XML
		Author:          p.authorExt.Extract(item),
		Date:            p.dateExt.Extract(item),
		ModifiedDate:    p.modifiedDateExt.Extract(item),
		Categories:      categories,
		Tags:            tags,
		GUID:            item.GUID,
		ParentID:        item.PostParent,
		Meta:            meta,
		FeaturedImage:   p.featuredImageExt.Extract(item),
	}
}

// Parse parses a WordPress WXR XML export file and converts it into Post instances.
// It filters for published posts only (post_type="post" and status="publish").
// Returns an error if the XML is malformed or cannot be read.
//
// The parser handles:
//   - Attachment URL resolution for featured images
//   - Author name resolution from meta fields or dc:creator
//   - Date normalization to RFC3339 format
//   - Excerpt fallback to subtitle meta field
//   - Featured image resolution from meta fields or attachments
func (p *Parser) Parse(r io.Reader) ([]Post, error) {
	p.logger.Printf("Starting WXR parsing")

	wxrDoc, err := p.decodeXML(r)
	if err != nil {
		return nil, err
	}

	p.logger.Printf("Parsed WXR document, found %d items", len(wxrDoc.Channel.Items))

	posts := make([]Post, 0)
	skippedCount := 0
	skippedByType := make(map[string]int)
	skippedByStatus := make(map[string]int)

	// Build attachment lookups (ID -> URL) and parent->attachments map
	attachmentIndex := buildAttachmentIndex(wxrDoc.Channel)

	// Build author lookup map (currently unused but kept for potential future use)
	_ = p.buildAuthorMap(wxrDoc.Channel)

	for i := range wxrDoc.Channel.Items {
		item := &wxrDoc.Channel.Items[i]

		// Filter: only include posts matching filter criteria
		// Track skipped items separately by type and status to match original behavior
		if item.PostType != "post" {
			skippedByType[item.PostType]++
			skippedCount++
			continue
		}

		if item.Status != "publish" {
			skippedByStatus[item.Status]++
			skippedCount++
			continue
		}

		// Validate essential fields
		if item.PostID == 0 {
			p.logger.Printf("Skipping item with missing post_id")
			skippedCount++
			continue
		}

		if item.Title == "" && item.ContentEncoded == "" && item.ExcerptEncoded == "" {
			p.logger.Printf("Skipping post %d: missing title, content, and excerpt", item.PostID)
			skippedCount++
			continue
		}

		// Transform item to Post
		post := p.transformItem(item, attachmentIndex)

		// Log warning if both Link and Slug are missing (malformed input)
		if post.Link == "" && post.Slug == "" {
			p.logger.Printf("Warning: Post %d has neither Link nor Slug - URL construction may fail", post.ID)
		}

		posts = append(posts, post)
	}

	p.logger.Printf("WXR parsing complete: %d posts extracted, %d items skipped", len(posts), skippedCount)
	if len(skippedByType) > 0 {
		p.logger.Printf("Skipped by type: %+v", skippedByType)
	}
	if len(skippedByStatus) > 0 {
		p.logger.Printf("Skipped by status: %+v", skippedByStatus)
	}

	return posts, nil
}

// Parse is a convenience function that parses a WXR file using the default parser.
// For more control over logging and configuration, use Parser instead.
func Parse(r io.Reader) ([]Post, error) {
	parser := NewParser()
	return parser.Parse(r)
}

// ParseWithContext parses a WXR file with context support for cancellation.
// This is a convenience function that uses the default parser.
func ParseWithContext(ctx context.Context, r io.Reader) ([]Post, error) {
	parser := NewParser()
	return parser.ParseWithContext(ctx, r)
}

// ParseWithContext parses a WordPress WXR XML export file with context support.
// It allows cancellation via the context. If the context is cancelled, parsing stops
// and returns the posts parsed so far along with a context error.
func (p *Parser) ParseWithContext(ctx context.Context, r io.Reader) ([]Post, error) {
	// Check context before starting
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	p.logger.Printf("Starting WXR parsing")

	wxrDoc, err := p.decodeXML(r)
	if err != nil {
		return nil, err
	}

	// Check context after decoding
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	p.logger.Printf("Parsed WXR document, found %d items", len(wxrDoc.Channel.Items))

	posts := make([]Post, 0)
	skippedCount := 0
	skippedByType := make(map[string]int)
	skippedByStatus := make(map[string]int)

	// Build attachment lookups (ID -> URL) and parent->attachments map
	attachmentIndex := buildAttachmentIndex(wxrDoc.Channel)

	// Build author lookup map (currently unused but kept for potential future use)
	_ = p.buildAuthorMap(wxrDoc.Channel)

	for i := range wxrDoc.Channel.Items {
		// Check context periodically during parsing
		select {
		case <-ctx.Done():
			p.logger.Printf("Parsing cancelled, returning %d posts parsed so far", len(posts))
			return posts, ctx.Err()
		default:
		}

		item := &wxrDoc.Channel.Items[i]

		// Filter: only include posts matching filter criteria
		// Track skipped items separately by type and status to match original behavior
		if item.PostType != "post" {
			skippedByType[item.PostType]++
			skippedCount++
			continue
		}

		if item.Status != "publish" {
			skippedByStatus[item.Status]++
			skippedCount++
			continue
		}

		// Validate essential fields
		if item.PostID == 0 {
			p.logger.Printf("Skipping item with missing post_id")
			skippedCount++
			continue
		}

		if item.Title == "" && item.ContentEncoded == "" && item.ExcerptEncoded == "" {
			p.logger.Printf("Skipping post %d: missing title, content, and excerpt", item.PostID)
			skippedCount++
			continue
		}

		// Transform item to Post
		post := p.transformItem(item, attachmentIndex)

		// Log warning if both Link and Slug are missing (malformed input)
		if post.Link == "" && post.Slug == "" {
			p.logger.Printf("Warning: Post %d has neither Link nor Slug - URL construction may fail", post.ID)
		}

		posts = append(posts, post)
	}

	p.logger.Printf("WXR parsing complete: %d posts extracted, %d items skipped", len(posts), skippedCount)
	if len(skippedByType) > 0 {
		p.logger.Printf("Skipped by type: %+v", skippedByType)
	}
	if len(skippedByStatus) > 0 {
		p.logger.Printf("Skipped by status: %+v", skippedByStatus)
	}

	return posts, nil
}
