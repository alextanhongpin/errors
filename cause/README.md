# Error Cause Package

A comprehensive Go error handling package that provides structured error management with support for error codes, attributes, validation, stack traces, and structured logging.

## Features

- **Structured Errors**: Rich error types with codes, names, messages, and additional context
- **Error Chaining**: Support for wrapping and unwrapping errors
- **Validation Framework**: Fluent validation API for complex data structures
- **Structured Logging**: Integration with Go's `slog` package
- **Stack Traces**: Optional stack trace capture for debugging
- **Type Safety**: Strong typing with generic support for validation

## Installation

```bash
go get github.com/alextanhongpin/errors/cause
```

## Quick Start

### Basic Error Creation

```go
package main

import (
    "fmt"
    "github.com/alextanhongpin/errors/cause"
    "github.com/alextanhongpin/errors/codes"
)

// Define error constants
var (
    ErrUserNotFound = cause.New(codes.NotFound, "UserNotFound", "User not found")
    ErrInvalidEmail = cause.New(codes.Invalid, "InvalidEmail", "Invalid email format")
)

func main() {
    err := ErrUserNotFound.WithDetails(map[string]any{
        "user_id": "123",
        "action":  "fetch",
    })
    
    fmt.Println(err.Error()) // Output: User not found
    fmt.Println(err.Code)    // Output: not_found
}
```

### Error with Stack Trace

```go
func riskyOperation() error {
    return ErrUserNotFound.WithStack()
}
```

### Error Chaining

```go
func fetchUser(id string) error {
    if err := validateUserID(id); err != nil {
        return ErrInvalidEmail.Wrap(err)
    }
    // ... fetch logic
    return nil
}
```

## Validation Framework

The package includes a powerful validation framework for validating complex data structures:

### Simple Validation

```go
type User struct {
    Name string
    Age  int
}

func (u *User) Validate() error {
    return cause.Map{
        "name": cause.Required(u.Name),
        "age":  cause.Optional(u.Age).When(u.Age < 0, "must be positive"),
    }.Err()
}
```

### Nested Validation

```go
type Address struct {
    Street string
    City   string
}

func (a *Address) Validate() error {
    return cause.Map{
        "street": cause.Required(a.Street),
        "city":   cause.Required(a.City),
    }.Err()
}

type User struct {
    Name    string
    Address *Address
}

func (u *User) Validate() error {
    return cause.Map{
        "name":    cause.Required(u.Name),
        "address": cause.Required(u.Address),
    }.Err()
}
```

### Slice Validation

```go
func validateUsers(users []User) error {
    return cause.SliceFunc(users, func(u User) error {
        return u.Validate()
    }).Validate()
}
```

## Structured Logging

The package integrates seamlessly with Go's `slog` package:

```go
import "log/slog"

func main() {
    logger := slog.Default()
    
    err := cause.New(codes.Invalid, "ValidationError", "Invalid input").
        WithAttrs(slog.String("field", "email")).
        WithDetails(map[string]any{
            "input": "invalid-email",
            "rule":  "email_format",
        })
    
    logger.Error("Validation failed", "error", err)
}
```

## Advanced Usage

### Custom Validation Conditions

```go
func (u *User) Validate() error {
    return cause.Map{
        "email": cause.Required(u.Email).
            When(!isValidEmail(u.Email), "invalid email format").
            When(isDomainBlocked(u.Email), "email domain not allowed"),
        "age": cause.Optional(u.Age).
            When(u.Age < 13, "under minimum age").
            When(u.Age > 120, "age not realistic"),
    }.Err()
}
```

### Error Type Checking

```go
func handleError(err error) {
    var causeErr *cause.Error
    if errors.As(err, &causeErr) {
        switch causeErr.Code {
        case codes.NotFound:
            // Handle not found
        case codes.Invalid:
            // Handle validation errors
        }
    }
}
```

### Validation Error Handling

```go
func processValidationError(err error) {
    if validationErr, ok := err.(interface{ Map() map[string]any }); ok {
        fieldErrors := validationErr.Map()
        for field, fieldErr := range fieldErrors {
            fmt.Printf("Field %s: %v\n", field, fieldErr)
        }
    }
}
```

## API Reference

### Error Type

The main `Error` type provides:

- `Code`: Error classification using codes package
- `Name`: Unique error type identifier
- `Message`: Human-readable error description
- `Attrs`: Structured logging attributes
- `Details`: Additional context data
- `Cause`: Wrapped underlying error
- `Stack`: Optional stack trace

### Methods

- `New(code, name, message, ...args) *Error`: Create new error
- `WithStack() *Error`: Add stack trace
- `WithDetails(map[string]any) *Error`: Add context details
- `WithAttrs(...slog.Attr) *Error`: Add logging attributes
- `Wrap(error) *Error`: Wrap another error
- `Clone() *Error`: Create deep copy

### Validation Functions

- `Required(val) *Builder`: Validate required field
- `Optional(val) *Builder`: Validate optional field
- `Map{}.Err()`: Validate multiple fields
- `SliceFunc(slice, func) sliceValidator`: Validate slice elements

### Builder Methods

- `When(condition, message) *Builder`: Conditional validation
- `Validate() error`: Execute validation chain

## Error Codes

This package works with the companion `codes` package that provides standard error classifications:

- `codes.OK`: Success
- `codes.Invalid`: Invalid input
- `codes.NotFound`: Resource not found
- `codes.AlreadyExists`: Resource already exists
- `codes.PermissionDenied`: Access denied
- `codes.Internal`: Internal server error
- And more...

## Best Practices

1. **Define Error Constants**: Create package-level error constants for reusable errors
2. **Use Meaningful Names**: Choose descriptive error names and codes
3. **Add Context**: Use `WithDetails()` to provide debugging context
4. **Implement Validation**: Use the validation framework for input validation
5. **Chain Errors**: Use `Wrap()` to preserve error context
6. **Structured Logging**: Leverage `slog` integration for better observability

## Examples

See the `examples_*_test.go` files for comprehensive usage examples including:

- Basic error creation and handling
- Validation patterns
- Logging integration
- Error chaining
- Complex validation scenarios

## Real-World Examples

The package includes comprehensive real-world examples demonstrating validation in different domains:

### Healthcare System (`examples_healthcare_test.go`)
- Patient record management with medical validations
- Healthcare-specific field validation (blood types, medical IDs, vital signs)
- Complex nested structures with emergency contacts and insurance information

### Educational Institution (`examples_education_test.go`)
- Student enrollment system with academic profile validation
- Course registration with prerequisites and scheduling
- Transcript management and GPA calculations

### IoT Device Management (`examples_iot_test.go`)
- Device configuration validation with network settings
- Sensor calibration and threshold management
- Security configuration and API key management
- Power management and location tracking

### Business Applications (`examples_business_test.go`, `examples_api_test.go`)
- Financial transaction validation
- API request/response validation
- E-commerce product and order management
- User registration and authentication

Each example demonstrates:
- Complex nested validation scenarios
- Domain-specific validation rules
- Real-world field constraints and business logic
- Best practices for structuring validation code

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests for any improvements.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
