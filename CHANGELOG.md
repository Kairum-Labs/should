# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - Initial Release

### Added

- **Core Assertions**: Support for equality, boolean, nil, type, and length checks.
- **Numeric Comparisons**: Assertions for greater/less than, with detailed difference reporting.
- **String Assertions**: Checks for prefixes, suffixes, and substrings with case-insensitive options and typo detection.
- **Collection Assertions**: Assertions for element containment, set membership, and duplicate detection.
- **Map Assertions**: Checks for the presence of keys and values.
- **Detailed Error Messages**: Rich, contextual feedback for failed assertions to simplify debugging.
- **Functional Options**: Support for custom failure messages using `should.WithMessage()`.
- **Panic Handling**: Assertions to verify when a function panics or not.
