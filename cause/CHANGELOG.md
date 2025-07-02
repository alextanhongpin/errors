# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of structured error handling package
- Core Error type with support for codes, names, messages, attributes, and details
- Comprehensive validation framework with fluent API
- Support for error wrapping and unwrapping
- Integration with Go's slog package for structured logging
- Stack trace capture functionality
- Validation for nested structures and slices
- Conditional validation with When() method
- Comprehensive documentation and examples
- Performance benchmarks
- Contributing guidelines

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)

## [1.0.0] - 2025-07-02

### Added
- Core structured error handling functionality
- Validation framework with Map, Required, Optional, and Builder types
- Error enrichment methods: WithDetails, WithStack, WithAttrs
- slog.LogValuer implementation for structured logging
- Support for error chaining and wrapping
- Comprehensive test suite with 100% coverage
- Documentation including README, ADR, and PRD
- Example code demonstrating best practices
- Performance benchmarks
- Contributing guidelines and code of conduct

### Features

#### Error Handling
- **Structured Errors**: Rich error types with codes, names, messages, and context
- **Error Codes**: Integration with codes package for standardized error classification
- **Error Chaining**: Support for wrapping and unwrapping errors with preserved context
- **Stack Traces**: Optional stack trace capture for debugging
- **Context Details**: Arbitrary key-value details for error context
- **Structured Logging**: Native slog integration with structured attributes

#### Validation Framework
- **Fluent API**: Chainable validation methods for readable code
- **Field Validation**: Required and optional field validation
- **Conditional Logic**: When() method for conditional validation rules
- **Nested Structures**: Automatic validation of nested validatable types
- **Slice Support**: Validation of slices with per-element error reporting
- **Custom Messages**: Support for custom validation error messages

#### Developer Experience
- **Type Safety**: Generic support where beneficial
- **Zero Allocation**: Optimized paths for common operations
- **Backward Compatibility**: Compatible with standard Go error interface
- **Comprehensive Documentation**: Extensive examples and API documentation
- **IDE Support**: Full autocomplete and type hints

### Technical Details

#### Performance
- Sub-microsecond error creation time
- Minimal memory allocation for common operations
- Lazy evaluation for expensive operations (stack traces)
- Zero-cost abstractions for unused features

#### Compatibility
- Go 1.21+ support
- Standard library only dependencies
- Full compatibility with existing error handling patterns
- Support for error wrapping and unwrapping

#### Testing
- 100% test coverage
- Comprehensive benchmark suite
- Integration tests for complex scenarios
- Example tests for documentation validation

### Breaking Changes
- N/A (initial release)

### Migration Guide
- N/A (initial release)

### Known Issues
- None at release time

### Acknowledgments
- Thanks to the Go team for the excellent slog package
- Inspired by various error handling patterns in the Go community
- Built on the foundation of Go's simple and effective error model
