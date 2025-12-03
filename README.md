# go-wxr

A Go library for parsing WordPress WXR (WordPress eXtended RSS) export files.

[![CI](https://github.com/rafaelranery/go-wxr/actions/workflows/ci.yml/badge.svg)](https://github.com/rafaelranery/go-wxr/actions/workflows/ci.yml)

## Overview

`go-wxr` provides a simple and efficient way to parse WordPress WXR XML export files into Go structs. It handles the complexities of WordPress export formats, including attachment resolution, date normalization, and metadata extraction.

## Features

- Parse WordPress WXR export files
- Filter for published posts only (configurable)
- Resolve attachment URLs for featured images
- Normalize dates to RFC3339 format
- Extract metadata (author, excerpt, featured images)
- Configurable logging (no-op by default)
- Comprehensive error handling

## Installation

```bash
go get github.com/rafaelrapnery/go-wxr
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/rafaelrapnery/go-wxr"
)

func main() {
    file, err := os.Open("export.xml")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    posts, err := wxr.Parse(file)
    if err != nil {
        panic(err)
    }
    
    for _, post := range posts {
        fmt.Printf("Post: %s (%s)\n", post.TitleRendered, post.Link)
    }
}
```

### Using a Parser with Custom Logging

```go
package main

import (
    "log"
    "os"
    
    "github.com/rafaelrapnery/go-wxr"
)

func main() {
    // Create a parser with custom logging
    logger := log.New(os.Stdout, "[wxr] ", log.LstdFlags)
    parser := wxr.NewParserWithStdLogger(logger)
    
    file, err := os.Open("export.xml")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    posts, err := parser.Parse(file)
    if err != nil {
        panic(err)
    }
    
    // Use posts...
}
```

### Using a Custom Logger Interface

```go
package main

import (
    "github.com/rafaelrapnery/go-wxr"
)

type myLogger struct{}

func (l *myLogger) Printf(format string, v ...any) {
    // Custom logging implementation
}

func main() {
    parser := wxr.NewParserWithLogger(&myLogger{})
    
    // Use parser...
}
```

## API Reference

### Types

#### Post

Represents a WordPress post parsed from a WXR export file.

```go
type Post struct {
    ID              int      // WordPress post ID
    TitleRendered   string   // Post title (may contain HTML)
    ContentRendered string   // Full post content (HTML)
    Slug            string   // URL-friendly post slug
    Link            string   // Canonical permalink URL
    Excerpt         string   // Post excerpt or summary
    Author          string   // Post author name
    Categories      []string // List of category names or IDs
    Date            string   // Publication date in RFC3339 format
    FeaturedImage   string   // URL of the featured image
}
```

#### Logger

Interface for logging operations. Implementations should handle log messages for debugging and informational purposes.

```go
type Logger interface {
    Printf(format string, v ...any)
}
```

### Functions

#### Parse

Convenience function that parses a WXR file using the default parser (no logging).

```go
func Parse(r io.Reader) ([]Post, error)
```

### Parser

#### NewParser

Creates a new Parser with the default no-op logger.

```go
func NewParser() *Parser
```

#### NewParserWithLogger

Creates a new Parser with a custom logger.

```go
func NewParserWithLogger(logger Logger) *Parser
```

#### NewParserWithStdLogger

Creates a new Parser using the standard log.Logger.

```go
func NewParserWithStdLogger(stdLogger *log.Logger) *Parser
```

#### SetLogger

Sets the logger for the parser.

```go
func (p *Parser) SetLogger(logger Logger)
```

#### Parse

Parses a WordPress WXR XML export file and converts it into Post instances.

```go
func (p *Parser) Parse(r io.Reader) ([]Post, error)
```

The parser:
- Filters for published posts only (`post_type="post"` and `status="publish"`)
- Resolves attachment URLs for featured images
- Handles author name resolution from meta fields or `dc:creator`
- Normalizes dates to RFC3339 format
- Extracts excerpts from meta fields if not present in the excerpt field
- Resolves featured images from meta fields or attachments

## Behavior

### Filtering

By default, the parser only includes posts with:
- `post_type="post"`
- `status="publish"`

All other items (pages, drafts, attachments, etc.) are skipped.

### Date Normalization

Dates are normalized to RFC3339 format. The parser attempts to parse dates in various WordPress formats:
- RFC3339
- `2006-01-02 15:04:05`
- `2006-01-02T15:04:05`
- RFC822/RFC822Z
- Other common WordPress date formats

If parsing fails, the original date string is returned.

### Featured Image Resolution

Featured images are resolved in the following order:
1. Custom meta fields: `banner_da_materia`, `banner_old`, `link_do_banner`
2. WordPress thumbnail ID (`_thumbnail_id`) pointing to an attachment
3. First attachment associated with the post

### Author Resolution

Author names are resolved in the following order:
1. Custom meta fields: `redator`, `autor`, `author_name`
2. `dc:creator` field from the RSS item

### Excerpt Fallback

If the excerpt field is empty, the parser attempts to use the `subtitulo` meta field as a fallback.

## Error Handling

The parser returns errors in the following cases:
- XML parsing fails (malformed XML)
- Root element is not `<rss>`
- I/O errors reading from the input

Posts with missing required fields (like `post_id`) are skipped and logged, but do not cause the parser to return an error.

## Logging

By default, the parser uses a no-op logger that discards all log messages. This prevents noisy output in library code. You can enable logging by:

1. Using `NewParserWithStdLogger` with a standard `log.Logger`
2. Implementing the `Logger` interface and using `NewParserWithLogger`
3. Calling `SetLogger` on an existing parser instance

Log messages include:
- Parsing start/completion
- Number of items found
- Skipped items (by type and status)
- Warnings about malformed data

## Testing

Run tests with:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## WordPress WXR Format

This library targets WordPress WXR export format version 1.2. It should work with exports from WordPress 2.1+.

For more information about the WXR format, see the [WordPress Codex](https://codex.wordpress.org/WordPress_Export).

