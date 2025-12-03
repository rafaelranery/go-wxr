# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-12-03

### Added
- **Categories and Tags parsing** - Extract categories and tags from WXR format
- **Post modification dates** - `ModifiedDate` field with RFC3339 normalization
- **Post meta access** - `Meta` map[string]string for all custom fields
- **Context support** - `ParseWithContext()` method for cancellation
- **GUID and ParentID** - Post GUID and parent post relationships
- Comprehensive tests for all new features

### Changed
- `Categories` field now properly populated (was previously empty)
- `Post` struct extended with new fields while maintaining backward compatibility

## [0.1.0] - 2025-12-03

### Added
- Initial release of go-wxr package
- Parse WordPress WXR export files into Go structs
- Attachment URL resolution for featured images
- Date normalization to RFC3339 format
- Metadata extraction (author, excerpt, featured images)
- Configurable logging interface
- Extensible filtering system
- Comprehensive test coverage
- Example programs demonstrating usage
- Full documentation (README, CONTRIBUTING, PACKAGE)

[0.1.1]: https://github.com/rafaelranery/go-wxr/releases/tag/v0.1.1
[0.1.0]: https://github.com/rafaelranery/go-wxr/releases/tag/v0.1.0

