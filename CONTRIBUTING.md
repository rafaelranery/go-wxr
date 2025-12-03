# Contributing to go-wxr

Thank you for your interest in contributing to `go-wxr`!

## Package Structure

The `go-wxr` package is organized as follows:

```
go-wxr/
├── wxr.go              # Main parser API and orchestration
├── post.go              # Post struct (public API)
├── logger.go           # Logger interface and implementations
├── filter.go           # Filter interface and default implementation
├── extractor.go        # Field extractors (author, excerpt, date, featured image)
├── attachments.go      # Attachment resolution logic
├── metadata.go         # Metadata extraction utilities
├── date.go             # Date normalization utilities
├── xml.go              # XML structs and decoding (internal)
├── wxr_test.go        # Test suite
├── example_test.go     # Example code (visible in GoDoc)
├── examples/           # Standalone example programs
│   ├── basic/         # Basic usage example
│   └── custom-logger/ # Custom logger example
├── testdata/           # Test fixtures (if needed)
└── .github/           # CI/CD configuration
```

## Code Organization

### Public API
- **`wxr.go`**: Main `Parser` struct and `Parse()` function
- **`post.go`**: `Post` struct representing parsed WordPress posts
- **`logger.go`**: `Logger` interface for logging
- **`filter.go`**: `Filter` interface for custom filtering

### Internal Implementation
- **`extractor.go`**: Extractors for transforming WXR items to Post fields
- **`attachments.go`**: Attachment URL resolution
- **`metadata.go`**: Metadata value extraction
- **`date.go`**: Date parsing and normalization
- **`xml.go`**: Internal XML structs (unexported)

### Testing
- **`wxr_test.go`**: Comprehensive test suite
- **`example_test.go`**: Example functions visible in GoDoc

## Adding New Features

1. **Keep the public API minimal**: Only export what users need
2. **Use interfaces for extensibility**: Allow users to customize behavior
3. **Maintain backward compatibility**: Don't break existing APIs
4. **Add tests**: All new features should have corresponding tests
5. **Update documentation**: Update README.md and add GoDoc comments

## Testing

Run tests with:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Code Style

- Follow standard Go formatting (`gofmt`)
- Run `go vet` before committing
- Keep functions focused and small
- Add GoDoc comments for exported types and functions
- Use meaningful variable names

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation if needed
7. Submit a pull request

## Questions?

Feel free to open an issue for questions or discussions about contributions.

