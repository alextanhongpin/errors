// Package cause provides structured error handling with support for error codes,
// attributes, details, stack traces, and validation. It extends Go's standard
// error interface with additional context and structured logging capabilities.
package cause

import (
	"fmt"
	"log/slog"
	"maps"
	"runtime"
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
	// Name is a stable, case-sensitive machine identifier for this error type
	Name string
	// Stack contains the stack trace when the error was created
	Stack string
}

// New creates a new Error with the specified code, name, and message.
// Additional arguments can include slog.Attr for structured logging attributes.
// The message supports fmt.Sprintf formatting with the provided args.
//
// Example:
//
//	err := New(codes.NotFound, "UserNotFound", "User %s not found", userID)
func New(code codes.Code, name, message string, args ...any) *Error {
	return &Error{
		Code:    code,
		Details: make(map[string]any),
		Message: fmt.Sprintf(message, args...),
		Name:    name,
	}
}

// Is reports whether this error matches the target error.
// Named nodes match by Code and Name; unnamed nodes also compare Message.
// Ordinary targets normalize to the unknown code, empty name, and their text.
// Identical ordinary error text intentionally matches. The target is not unwrapped.
func (e *Error) Is(target error) bool {
	if e == nil || nilError(target) || !e.Code.Valid() {
		return false
	}
	code, name, message := codes.Unknown, "", ""
	if other, ok := target.(*Error); ok {
		code, name, message = other.Code, other.Name, other.Message
	} else {
		message = target.Error()
	}
	return e.Code == code && e.Name == name && (e.Name != "" || e.Message == message)
}

// Unwrap returns the underlying cause error, supporting Go's error unwrapping.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Error returns the error message, implementing the standard error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil && e.Cause.Error() != "" {
		if len(e.Stack) > 0 {
			return fmt.Sprintf("%s\n\t%s\nCaused by: %s", e.Message, e.Stack, e.Cause)
		}
		return fmt.Sprintf("%s\nCaused by: %s", e.Message, e.Cause)
	}
	if len(e.Stack) > 0 {
		return fmt.Sprintf("%s\n\t%s", e.Message, e.Stack)
	}

	return e.Message
}

// LogValue implements slog.LogValuer for structured logging.
// It returns a grouped slog.Value containing all error context including
// message, code, name, attributes, details, and cause.
func (e Error) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("message", e.Message),
		slog.String("code", e.Code.String()),
		slog.String("name", e.Name),
	}

	if len(e.Attrs) > 0 {
		attrs = append(attrs, slog.GroupAttrs("data", e.Attrs...))
	}
	if e.Cause != nil {
		attrs = append(attrs, slog.Any("cause", e.Cause))
	}
	if len(e.Details) > 0 {
		attrs = append(attrs, slog.Any("details", e.Details))
	}

	return slog.GroupValue(attrs...)
}

// Clone copies the error and its top-level attribute slice and details map.
// Nested detail values and the underlying cause are shared.
func (e *Error) Clone() *Error {
	return &Error{
		Attrs:   slices.Clone(e.Attrs),
		Cause:   e.Cause,
		Code:    e.Code,
		Details: maps.Clone(e.Details),
		Message: e.Message,
		Name:    e.Name,
		Stack:   e.Stack,
	}
}

// WithDetails returns a new error with additional details merged in.
// Existing details are preserved, new details override existing keys.
func (e *Error) WithDetails(details map[string]any) *Error {
	err := e.Clone()
	if err.Details == nil {
		err.Details = make(map[string]any)
	}
	maps.Copy(err.Details, details)
	return err
}

// WithStack returns a new error with the current stack trace captured.
func (e *Error) WithStack() *Error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}

	err := e.Clone()
	err.Stack = fmt.Sprintf("at %s.%d", file, line)
	return err
}

// WithAttrs returns a new error with additional structured logging attributes.
func (e *Error) WithAttrs(attrs ...slog.Attr) *Error {
	err := e.Clone()
	err.Attrs = append(err.Attrs, attrs...)
	return err
}

// WithMessage returns a new error with the specified message formatted with args.
func (e *Error) WithMessage(message string, args ...any) *Error {
	err := e.Clone()
	err.Message = fmt.Sprintf(message, args...)
	return err
}

// WithCause returns a new error with the specified cause error.
func (e *Error) WithCause(cause error) *Error {
	err := e.Clone()
	err.Cause = cause
	return err
}
