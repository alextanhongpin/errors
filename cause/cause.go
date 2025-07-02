// Package cause provides structured error handling with support for error codes,
// attributes, details, stack traces, and validation. It extends Go's standard
// error interface with additional context and structured logging capabilities.
package cause

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"runtime/debug"
	"slices"

	"github.com/alextanhongpin/errors/codes"
)

// Error represents a structured error with additional context including error codes,
// attributes, details, stack traces, and nested causes. It implements the standard
// error interface and provides enhanced logging capabilities.
type Error struct {
	// Attrs contains structured logging attributes
	Attrs []slog.Attr
	// Cause is the underlying wrapped error
	Cause error
	// Code represents the error classification
	Code codes.Code
	// Details contains additional context as key-value pairs
	Details map[string]any
	// Message is the human-readable error message
	Message string
	// Name is a unique identifier for this error type
	Name string
	// Stack contains the stack trace when the error was created
	Stack []byte
}

// New creates a new Error with the specified code, name, and message.
// Additional arguments can include slog.Attr for structured logging attributes.
// The message supports fmt.Sprintf formatting with the provided args.
//
// Example:
//
//	err := New(codes.NotFound, "UserNotFound", "User %s not found", userID)
//	errWithAttrs := New(codes.Invalid, "ValidationError", "Invalid input",
//	    slog.String("field", "email"), slog.Int("length", len(email)))
func New(code codes.Code, name, message string, args ...any) *Error {
	var attrs []slog.Attr
	var formatArgs []any

	// Separate slog.Attr from formatting arguments
	for _, arg := range args {
		if attr, ok := arg.(slog.Attr); ok {
			attrs = append(attrs, attr)
		} else {
			formatArgs = append(formatArgs, arg)
		}
	}

	return &Error{
		Attrs:   attrs,
		Code:    code,
		Details: make(map[string]any),
		Message: fmt.Sprintf(message, formatArgs...),
		Name:    name,
	}
}

// Is reports whether this error matches the target error.
// Two errors match if they have the same Code and Name.
func (e *Error) Is(target error) bool {
	var other *Error
	if !errors.As(target, &other) {
		return false
	}

	return e.Code == other.Code &&
		e.Name == other.Name
}

// Unwrap returns the underlying cause error, supporting Go's error unwrapping.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Wrap sets the cause error and returns the current error for method chaining.
// This allows building error chains while preserving the original error context.
func (e *Error) Wrap(err error) *Error {
	e.Cause = err
	return e
}

// Error returns the error message, implementing the standard error interface.
func (e *Error) Error() string {
	return e.Message
}

// LogValue implements slog.LogValuer for structured logging.
// It returns a grouped slog.Value containing all error context including
// message, code, name, attributes, details, and cause.
func (e Error) LogValue() slog.Value {
	attrs := append([]slog.Attr{
		slog.String("message", e.Message),
		slog.String("code", e.Code.String()),
		slog.String("name", e.Name),
	}, e.Attrs...)

	if len(e.Details) > 0 {
		attrs = append(attrs, slog.Any("details", e.Details))
	}
	if e.Cause != nil {
		attrs = append(attrs, slog.Any("cause", e.Cause))
	}

	return slog.GroupValue(attrs...)
}

// StackTrace returns the captured stack trace.
func (e *Error) StackTrace() []byte {
	return e.Stack
}

// Clone creates a deep copy of the error, allowing safe modification
// without affecting the original error.
func (e *Error) Clone() *Error {
	return &Error{
		Attrs:   slices.Clone(e.Attrs),
		Cause:   e.Cause,
		Code:    e.Code,
		Details: maps.Clone(e.Details),
		Message: e.Message,
		Name:    e.Name,
		Stack:   bytes.Clone(e.Stack),
	}
}

// WithDetails returns a new error with additional details merged in.
// Existing details are preserved, new details override existing keys.
func (e *Error) WithDetails(details map[string]any) *Error {
	err := e.Clone()
	maps.Copy(err.Details, details)
	return err
}

// WithStack returns a new error with the current stack trace captured.
func (e *Error) WithStack() *Error {
	err := e.Clone()
	err.Stack = debug.Stack()
	return err
}

// WithAttrs returns a new error with additional structured logging attributes.
func (e *Error) WithAttrs(attrs ...slog.Attr) *Error {
	err := e.Clone()
	err.Attrs = append(err.Attrs, attrs...)
	return err
}
