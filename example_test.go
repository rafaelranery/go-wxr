package wxr_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/rafaelranery/go-wxr"
)

func ExampleParse() {
	// Example WXR XML content
	wxrXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:excerpt="http://wordpress.org/export/1.2/excerpt/"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Example Site</title>
	<link>https://example.com</link>
	<item>
		<title><![CDATA[Hello World]]></title>
		<link>https://example.com/hello-world</link>
		<content:encoded><![CDATA[<p>Welcome to WordPress!</p>]]></content:encoded>
		<excerpt:encoded><![CDATA[Welcome post]]></excerpt:encoded>
		<wp:post_id>1</wp:post_id>
		<wp:post_date_gmt><![CDATA[2025-01-01 12:00:00]]></wp:post_date_gmt>
		<wp:post_name><![CDATA[hello-world]]></wp:post_name>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	reader := strings.NewReader(wxrXML)
	posts, err := wxr.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		fmt.Printf("Post ID: %d\n", post.ID)
		fmt.Printf("Title: %s\n", post.TitleRendered)
		fmt.Printf("Link: %s\n", post.Link)
	}

	// Output:
	// Post ID: 1
	// Title: Hello World
	// Link: https://example.com/hello-world
}

func ExampleParser_Parse() {
	// Example WXR XML content
	wxrXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wp="http://wordpress.org/export/1.2/">
<channel>
	<title>Example Site</title>
	<item>
		<title><![CDATA[Example Post]]></title>
		<link>https://example.com/example-post</link>
		<content:encoded><![CDATA[<p>Content here</p>]]></content:encoded>
		<wp:post_id>42</wp:post_id>
		<wp:post_name><![CDATA[example-post]]></wp:post_name>
		<wp:post_type><![CDATA[post]]></wp:post_type>
		<wp:status><![CDATA[publish]]></wp:status>
	</item>
</channel>
</rss>`

	// Create a parser (uses no-op logger by default)
	parser := wxr.NewParser()

	// To enable logging, use:
	// logger := log.New(os.Stdout, "[wxr] ", log.LstdFlags)
	// parser := wxr.NewParserWithStdLogger(logger)

	reader := strings.NewReader(wxrXML)
	posts, err := parser.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed %d posts\n", len(posts))
	// Output: Parsed 1 posts
}
