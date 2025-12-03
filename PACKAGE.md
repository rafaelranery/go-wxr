# Package Structure

This document describes the internal organization of the `go-wxr` package.

## File Organization

### Core API Files (Root Package)

- **`wxr.go`**: Main parser implementation, public API (`Parser`, `Parse()`)
- **`post.go`**: Public `Post` struct representing parsed WordPress posts
- **`logger.go`**: Public `Logger` interface and implementations
- **`filter.go`**: Public `Filter` interface and `DefaultFilter` implementation

### Implementation Files (Root Package)

- **`extractor.go`**: Field extractors for transforming WXR items:
  - `AuthorExtractor`: Resolves author names
  - `ExcerptExtractor`: Extracts excerpts
  - `DateExtractor`: Extracts and normalizes dates
  - `FeaturedImageExtractor`: Resolves featured image URLs

- **`attachments.go`**: Attachment resolution logic:
  - `AttachmentIndex`: Maps attachment IDs to URLs
  - `buildAttachmentIndex()`: Builds the attachment index from WXR items

- **`metadata.go`**: Metadata extraction utilities:
  - `getMetaValue()`: Searches for meta values by key
  - `cleanMetaValue()`: Cleans and validates meta values

- **`date.go`**: Date normalization:
  - `normalizeWXRDate()`: Converts WordPress dates to RFC3339

- **`xml.go`**: Internal XML structs (unexported):
  - `wxr`, `channel`, `item`, `wpAuthor`, `postMeta`

### Test Files

- **`wxr_test.go`**: Comprehensive test suite
- **`example_test.go`**: Example functions (visible in GoDoc)

### Documentation

- **`README.md`**: User-facing documentation
- **`CONTRIBUTING.md`**: Contribution guidelines and package structure
- **`PACKAGE.md`**: This file - internal package documentation
- **`LICENSE`**: MIT License

### Examples

- **`examples/basic/`**: Basic usage example
- **`examples/custom-logger/`**: Custom logger example

### Test Data

- **`testdata/`**: Directory for test fixtures (currently empty, tests use inline XML)

## Design Principles

1. **Separation of Concerns**: Each file has a single, focused responsibility
2. **Public API**: Only essential types and functions are exported
3. **Extensibility**: Interfaces allow customization without modifying core code
4. **Testability**: Components can be tested in isolation
5. **Backward Compatibility**: Public API remains stable

## Adding New Features

When adding new features:

1. **Public API**: Add to `wxr.go` or create a new file if it's a major feature
2. **Internal Logic**: Add to appropriate implementation file or create new one
3. **Tests**: Add tests to `wxr_test.go` or create `*_test.go` for new files
4. **Documentation**: Update README.md and add GoDoc comments

## File Naming Conventions

- Use lowercase with underscores: `extractor.go`, `wxr_test.go`
- Test files: `*_test.go`
- Example files: `example_test.go` (for GoDoc) or `examples/` directory (standalone)

