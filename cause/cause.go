package cause

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"runtime/debug"

	"github.com/alextanhongpin/errors/codes"
)

type Error struct {
	Code    codes.Code
	Name    string
	Message string
	Details map[string]any
	Stack   []byte
	Cause   error
}

func New(code codes.Code, name, message string, args ...any) *Error {
	return &Error{
		Code:    code,
		Name:    name,
		Message: fmt.Sprintf(message, args...),
		Details: make(map[string]any),
	}
}

func (e *Error) Is(target error) bool {
	var other *Error
	if !errors.As(target, &other) {
		return false
	}

	return e.Code == other.Code &&
		e.Name == other.Name
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Wrap(err error) *Error {
	e.Cause = err
	return e
}

func (e *Error) Error() string {
	return e.Message
}

func (e Error) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("message", e.Message),
		slog.String("code", e.Code.String()),
		slog.String("name", e.Name),
	}

	if len(e.Details) > 0 {
		attrs = append(attrs, slog.Any("details", e.Details))
	}
	if e.Cause != nil {
		attrs = append(attrs, slog.Any("cause", e.Cause))
	}

	return slog.GroupValue(attrs...)
}

func (e *Error) StackTrace() []byte {
	return e.Stack
}

func (e *Error) Clone() *Error {
	return &Error{
		Code:    e.Code,
		Name:    e.Name,
		Message: e.Message,
		Details: maps.Clone(e.Details),
		Stack:   bytes.Clone(e.Stack),
		Cause:   e.Cause,
	}
}

func (e *Error) WithDetails(details map[string]any) *Error {
	err := e.Clone()
	maps.Copy(err.Details, details)
	return err
}

func (e *Error) WithStack() *Error {
	err := e.Clone()
	err.Stack = debug.Stack()
	return err
}
