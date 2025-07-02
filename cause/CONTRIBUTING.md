# Contributing to Error Cause Package

We welcome contributions to the Error Cause package! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Review Process](#review-process)

## Code of Conduct

This project adheres to a code of conduct that promotes a welcoming and inclusive environment. By participating, you agree to uphold this standard.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes and test them
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for convenience)

### Local Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/errors
cd errors/cause

# Install dependencies
go mod download

# Run tests to ensure everything works
go test ./...

# Run benchmarks
go test -bench=. -benchmem
```

## Making Changes

### Branch Naming

Use descriptive branch names that indicate the type of change:

- `feature/validation-improvements`
- `fix/error-wrapping-bug`
- `docs/api-documentation`
- `refactor/validation-internals`

### Code Style

We follow standard Go conventions:

- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Keep line length under 100 characters when possible

### API Design Principles

When adding new features, follow these principles:

1. **Simplicity**: APIs should be easy to use for common cases
2. **Composability**: Features should work well together
3. **Type Safety**: Leverage Go's type system
4. **Performance**: Maintain minimal overhead
5. **Backward Compatibility**: Don't break existing APIs

### Error Handling

- All public functions should handle errors gracefully
- Use structured errors from this package for internal errors
- Include appropriate context in error messages
- Test error conditions thoroughly

## Testing

### Test Coverage

- Maintain 100% test coverage for new code
- Include unit tests for all public APIs
- Add integration tests for complex scenarios
- Include benchmark tests for performance-critical code

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem

# Run tests with race detection
go test -race ./...
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Example Tests**: Runnable examples in documentation
4. **Benchmark Tests**: Performance measurements

### Writing Good Tests

- Use table-driven tests when appropriate
- Test both success and failure cases
- Include edge cases and boundary conditions
- Use descriptive test names
- Keep tests simple and focused

Example test structure:

```go
func TestErrorValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    User
        wantErr  bool
        errField string
    }{
        {
            name: "valid user",
            input: User{
                ID:    "123",
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: false,
        },
        {
            name: "missing email",
            input: User{
                ID:   "123",
                Name: "Test User",
            },
            wantErr:  true,
            errField: "email",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Add godoc comments for all public APIs
- Include examples in documentation
- Document complex algorithms and design decisions
- Keep comments up to date with code changes

### Documentation Standards

- Use complete sentences in godoc comments
- Start comments with the function/type name
- Include usage examples for complex APIs
- Document error conditions and return values

Example documentation:

```go
// New creates a new Error with the specified code, name, and message.
// Additional arguments can include slog.Attr for structured logging attributes.
// The message supports fmt.Sprintf formatting with the provided args.
//
// Example:
//   err := New(codes.NotFound, "UserNotFound", "User %s not found", userID)
//   errWithAttrs := New(codes.Invalid, "ValidationError", "Invalid input",
//       slog.String("field", "email"), slog.Int("length", len(email)))
func New(code codes.Code, name, message string, args ...any) *Error {
    // implementation
}
```

### README Updates

Update the README.md when:
- Adding new features
- Changing APIs
- Adding new examples
- Updating installation instructions

## Submitting Changes

### Pull Request Guidelines

1. **Clear Description**: Explain what changes you made and why
2. **Issue Reference**: Link to related issues if applicable
3. **Breaking Changes**: Clearly mark any breaking changes
4. **Testing**: Confirm all tests pass
5. **Documentation**: Update documentation as needed

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Documentation
- [ ] Code comments updated
- [ ] README updated (if needed)
- [ ] Examples added/updated (if needed)

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Changes are backward compatible (or breaking changes documented)
```

### Commit Message Format

Use conventional commit format:

```
type(scope): description

body (optional)

footer (optional)
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test updates
- `chore`: Build/tooling changes

Examples:
```
feat(validation): add support for custom validation messages
fix(error): resolve stack trace memory leak
docs(readme): update validation examples
```

## Review Process

### Review Criteria

Pull requests are reviewed for:

1. **Correctness**: Code works as intended
2. **Performance**: No unnecessary performance degradation
3. **Style**: Follows Go conventions and project style
4. **Tests**: Adequate test coverage
5. **Documentation**: Appropriate documentation
6. **Compatibility**: Maintains backward compatibility

### Review Timeline

- Initial review within 48 hours
- Response to feedback within 24 hours
- Final approval within 1 week for standard changes

### Addressing Feedback

- Respond to all review comments
- Make requested changes promptly
- Ask for clarification if feedback is unclear
- Update the PR description if scope changes

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version tagged
- [ ] Release notes published

## Getting Help

If you need help or have questions:

1. Check existing issues and documentation
2. Search previous discussions
3. Open a new issue with the `question` label
4. Join our community discussions

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Special recognition for significant contributions

Thank you for contributing to the Error Cause package!
