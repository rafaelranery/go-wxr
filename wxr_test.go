package wxr

import (
	"context"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		xml       string
		wantCount int
		wantErr   bool
		validate  func(t *testing.T, posts []Post)
	}{
		{
			name: "valid post with all fields",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:excerpt="http://wordpress.org/export/1.2/excerpt/"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<link>https://example.com</link>
	<item>
		<title><![CDATA[Test Post Title]]></title>
		<link>https://example.com/test-post</link>
		<pubDate>Sun, 01 Jun 2025 14:00:51 +0000</pubDate>
		<dc:creator><![CDATA[John Doe]]></dc:creator>
		<content:encoded><![CDATA[<p>This is the post content</p>]]></content:encoded>
		<excerpt:encoded><![CDATA[This is an excerpt]]></excerpt:encoded>
		<wp:post_id>123</wp:post_id>
		<wp:post_date><![CDATA[2025-06-01 11:00:51]]></wp:post_date>
		<wp:post_date_gmt><![CDATA[2025-06-01 14:00:51]]></wp:post_date_gmt>
		<wp:post_name><![CDATA[test-post]]></wp:post_name>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`,
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, posts []Post) {
				if len(posts) != 1 {
					t.Fatalf("expected 1 post, got %d", len(posts))
				}
				post := posts[0]
				if post.ID != 123 {
					t.Errorf("expected ID 123, got %d", post.ID)
				}
				if post.TitleRendered != "Test Post Title" {
					t.Errorf("expected title 'Test Post Title', got '%s'", post.TitleRendered)
				}
				if post.ContentRendered != "<p>This is the post content</p>" {
					t.Errorf("unexpected content: %s", post.ContentRendered)
				}
				if post.Excerpt != "This is an excerpt" {
					t.Errorf("expected excerpt 'This is an excerpt', got '%s'", post.Excerpt)
				}
				if post.Slug != "test-post" {
					t.Errorf("expected slug 'test-post', got '%s'", post.Slug)
				}
				if post.Link != "https://example.com/test-post" {
					t.Errorf("expected link 'https://example.com/test-post', got '%s'", post.Link)
				}
				if post.Author != "John Doe" {
					t.Errorf("expected author 'John Doe', got '%s'", post.Author)
				}
			},
		},
		{
			name: "uses meta author, subtitle excerpt, and attachment thumbnail",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:excerpt="http://wordpress.org/export/1.2/excerpt/"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<link>https://example.com</link>
	<wp:base_site_url>https://example.com</wp:base_site_url>
	<wp:base_blog_url>https://example.com</wp:base_blog_url>
	<item>
		<title><![CDATA[Post With Meta]]></title>
		<link>https://example.com/post-with-meta</link>
		<content:encoded><![CDATA[<p>Body</p>]]></content:encoded>
		<excerpt:encoded><![CDATA[]]></excerpt:encoded>
		<wp:post_id>200</wp:post_id>
		<wp:post_date><![CDATA[2025-07-01 10:00:00]]></wp:post_date>
		<wp:post_date_gmt><![CDATA[2025-07-01 13:00:00]]></wp:post_date_gmt>
		<wp:post_name><![CDATA[post-with-meta]]></wp:post_name>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
		<wp:postmeta>
			<wp:meta_key><![CDATA[redator]]></wp:meta_key>
			<wp:meta_value><![CDATA[Marco D'Eramo]]></wp:meta_value>
		</wp:postmeta>
		<wp:postmeta>
			<wp:meta_key><![CDATA[subtitulo]]></wp:meta_key>
			<wp:meta_value><![CDATA[Resumo curto]]></wp:meta_value>
		</wp:postmeta>
		<wp:postmeta>
			<wp:meta_key><![CDATA[_thumbnail_id]]></wp:meta_key>
			<wp:meta_value><![CDATA[555]]></wp:meta_value>
		</wp:postmeta>
	</item>
	<item>
		<title><![CDATA[attachment-thumb]]></title>
		<link>https://example.com/attachments/attachment-thumb</link>
		<guid isPermaLink="false">https://example.com/wp-content/uploads/2023/10/thumb.png</guid>
		<wp:post_id>555</wp:post_id>
		<wp:post_parent>200</wp:post_parent>
		<wp:post_type><![CDATA[attachment]]></wp:post_type>
		<wp:status><![CDATA[inherit]]></wp:status>
		<wp:attachment_url><![CDATA[https://cdn.example.com/thumb.png]]></wp:attachment_url>
	</item>
</channel>
</rss>`,
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, posts []Post) {
				post := posts[0]
				if post.Author != "Marco D'Eramo" {
					t.Fatalf("expected author from meta, got %q", post.Author)
				}
				if post.Excerpt != "Resumo curto" {
					t.Fatalf("expected excerpt fallback to meta, got %q", post.Excerpt)
				}
				if post.FeaturedImage != "https://cdn.example.com/thumb.png" {
					t.Fatalf("expected featured image from attachment, got %q", post.FeaturedImage)
				}
			},
		},
		{
			name: "skips non-post items",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Page Title]]></title>
		<wp:post_id>456</wp:post_id>
		<wp:post_type><![CDATA[page]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
	<item>
		<title><![CDATA[Draft Post]]></title>
		<wp:post_id>789</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[draft]]></wp:status>
	</item>
	<item>
		<title><![CDATA[Published Post]]></title>
		<wp:post_id>101</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`,
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, posts []Post) {
				if len(posts) != 1 {
					t.Fatalf("expected 1 post, got %d", len(posts))
				}
				if posts[0].ID != 101 {
					t.Errorf("expected post ID 101, got %d", posts[0].ID)
				}
			},
		},
		{
			name: "skips items with missing post_id",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post Without ID]]></title>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "handles missing optional fields",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Minimal Post]]></title>
		<content:encoded><![CDATA[Content here]]></content:encoded>
		<wp:post_id>999</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`,
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, posts []Post) {
				if len(posts) != 1 {
					t.Fatalf("expected 1 post, got %d", len(posts))
				}
				post := posts[0]
				if post.ID != 999 {
					t.Errorf("expected ID 999, got %d", post.ID)
				}
				if post.Excerpt != "" {
					t.Errorf("expected empty excerpt, got '%s'", post.Excerpt)
				}
			},
		},
		{
			name:      "malformed XML - not XML at all",
			xml:       `not xml at all`,
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "malformed XML - invalid root element",
			xml: `<?xml version="1.0"?>
<notrss>
	<channel></channel>
</notrss>`,
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.xml)
			posts, err := Parse(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(posts) != tt.wantCount {
				t.Errorf("Parse() returned %d posts, want %d", len(posts), tt.wantCount)
			}

			if tt.validate != nil {
				tt.validate(t, posts)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	// Test that Parser.Parse works the same as Parse
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Test Post]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	parser := NewParser()
	reader := strings.NewReader(xml)
	posts, err := parser.Parse(reader)

	if err != nil {
		t.Fatalf("Parser.Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("Parser.Parse() returned %d posts, want 1", len(posts))
	}

	if posts[0].ID != 1 {
		t.Errorf("expected post ID 1, got %d", posts[0].ID)
	}
}

func TestParser_SetLogger(t *testing.T) {
	parser := NewParser()

	// Test setting a nil logger (should use no-op)
	parser.SetLogger(nil)
	if parser.logger == nil {
		t.Error("SetLogger(nil) should set a no-op logger, not nil")
	}

	// Test setting a custom logger
	customLogger := &testLogger{logs: make([]string, 0)}
	parser.SetLogger(customLogger)

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Test Post]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	_, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify that the logger was called
	if len(customLogger.logs) == 0 {
		t.Error("expected logger to be called during parsing")
	}
}

type testLogger struct {
	logs []string
}

func (t *testLogger) Printf(format string, v ...any) {
	t.logs = append(t.logs, format)
}

func TestParse_CategoriesAndTags(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post With Categories]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
		<category domain="category" nicename="tech"><![CDATA[Technology]]></category>
		<category domain="category" nicename="news"><![CDATA[News]]></category>
		<category domain="post_tag" nicename="golang"><![CDATA[Go]]></category>
		<category domain="post_tag" nicename="programming"><![CDATA[Programming]]></category>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	posts, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	post := posts[0]
	if len(post.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d: %v", len(post.Categories), post.Categories)
	}
	if !contains(post.Categories, "Technology") || !contains(post.Categories, "News") {
		t.Errorf("expected categories to contain 'Technology' and 'News', got %v", post.Categories)
	}

	if len(post.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(post.Tags), post.Tags)
	}
	if !contains(post.Tags, "Go") || !contains(post.Tags, "Programming") {
		t.Errorf("expected tags to contain 'Go' and 'Programming', got %v", post.Tags)
	}
}

func TestParse_ModifiedDate(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post With Modified Date]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_date><![CDATA[2025-01-01 10:00:00]]></wp:post_date>
		<wp:post_date_gmt><![CDATA[2025-01-01 13:00:00]]></wp:post_date_gmt>
		<wp:post_modified><![CDATA[2025-01-15 14:30:00]]></wp:post_modified>
		<wp:post_modified_gmt><![CDATA[2025-01-15 17:30:00]]></wp:post_modified_gmt>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	posts, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	post := posts[0]
	if post.ModifiedDate == "" {
		t.Error("expected ModifiedDate to be set")
	}
	// Modified date should be normalized to RFC3339
	if !strings.Contains(post.ModifiedDate, "T") {
		t.Errorf("expected RFC3339 format for ModifiedDate, got %q", post.ModifiedDate)
	}
}

func TestParse_MetaFields(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post With Meta]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
		<wp:postmeta>
			<wp:meta_key><![CDATA[custom_field]]></wp:meta_key>
			<wp:meta_value><![CDATA[custom_value]]></wp:meta_value>
		</wp:postmeta>
		<wp:postmeta>
			<wp:meta_key><![CDATA[another_field]]></wp:meta_key>
			<wp:meta_value><![CDATA[another_value]]></wp:meta_value>
		</wp:postmeta>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	posts, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	post := posts[0]
	if post.Meta == nil {
		t.Fatal("expected Meta map to be initialized")
	}
	if len(post.Meta) != 2 {
		t.Errorf("expected 2 meta fields, got %d", len(post.Meta))
	}
	if post.Meta["custom_field"] != "custom_value" {
		t.Errorf("expected meta['custom_field'] = 'custom_value', got %q", post.Meta["custom_field"])
	}
	if post.Meta["another_field"] != "another_value" {
		t.Errorf("expected meta['another_field'] = 'another_value', got %q", post.Meta["another_field"])
	}
}

func TestParse_GUIDAndParentID(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Child Post]]></title>
		<link>https://example.com/child-post</link>
		<guid isPermaLink="false">https://example.com/?p=2</guid>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>2</wp:post_id>
		<wp:post_parent>1</wp:post_parent>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	posts, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	post := posts[0]
	if post.GUID != "https://example.com/?p=2" {
		t.Errorf("expected GUID 'https://example.com/?p=2', got %q", post.GUID)
	}
	if post.ParentID != 1 {
		t.Errorf("expected ParentID 1, got %d", post.ParentID)
	}
}

func TestParse_BackwardCompatibility(t *testing.T) {
	// Test that existing fields still work as before
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Test Post]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(xml)
	posts, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	post := posts[0]
	// Verify existing fields still work
	if post.ID != 1 {
		t.Errorf("expected ID 1, got %d", post.ID)
	}
	if post.TitleRendered != "Test Post" {
		t.Errorf("expected title 'Test Post', got %q", post.TitleRendered)
	}
	// Verify new fields have sensible defaults
	if post.Categories == nil {
		t.Error("expected Categories to be initialized (empty slice)")
	}
	if post.Tags == nil {
		t.Error("expected Tags to be initialized (empty slice)")
	}
	if post.Meta == nil {
		t.Error("expected Meta to be initialized (empty map)")
	}
}

func TestParseWithContext(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post 1]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
	<item>
		<title><![CDATA[Post 2]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>2</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	ctx := context.Background()
	reader := strings.NewReader(xml)
	posts, err := ParseWithContext(ctx, reader)
	if err != nil {
		t.Fatalf("ParseWithContext() error = %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestParseWithContext_Cancellation(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Test Site</title>
	<item>
		<title><![CDATA[Post 1]]></title>
		<content:encoded><![CDATA[Content]]></content:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	reader := strings.NewReader(xml)
	posts, err := ParseWithContext(ctx, reader)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
	// Should return empty posts when cancelled before parsing
	if len(posts) != 0 {
		t.Errorf("expected 0 posts when cancelled, got %d", len(posts))
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
