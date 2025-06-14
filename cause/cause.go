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

type Error struct {
	Attrs   []slog.Attr
	Cause   error
	Code    codes.Code
	Details map[string]any
	Message string
	Name    string
	Stack   []byte
}

func New(code codes.Code, name, message string, args ...any) *Error {
	var attrs []slog.Attr
	for _, arg := range args {
		if _, ok := arg.(slog.Attr); ok {
			attrs = append(attrs, arg.(slog.Attr))
		}
	}

	return &Error{
		Attrs:   attrs,
		Code:    code,
		Details: make(map[string]any),
		Message: fmt.Sprintf(message, args...),
		Name:    name,
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

func (e *Error) StackTrace() []byte {
	return e.Stack
}

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

func (e *Error) WithAttrs(attrs ...slog.Attr) *Error {
	err := e.Clone()
	err.Attrs = append(err.Attrs, attrs...)
	return err
}
