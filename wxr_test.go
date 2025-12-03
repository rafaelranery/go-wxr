package wxr

import (
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
