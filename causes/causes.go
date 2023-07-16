// package causes represents an error cause. The name is choosen to avoid
// conflict with the errors package.
package causes

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/codes"
)

// Cause represents the error cause.
type Cause struct {
	code codes.Code
	kind string
	msg  string
}

// New returns a new Sentinel error.
func New(code codes.Code, kind, msg string, args ...any) error {
	return &Cause{
		code: code,
		kind: kind,
		msg:  fmt.Sprintf(msg, args...),
	}
}

// New returns a new Sentinel error.
func NewWithHint[T any](code codes.Code, kind, msg string, args ...any) Hint[T] {
	return New(code, kind, msg, args...)
}

func (c *Cause) Code() codes.Code {
	return c.code
}

// Kind is the type of error the cause represents.
//
// Can be of different format:
// - resource based, e.g. user_not_found, payout_declined
// - filepath based, e.g. user/not_found, payout/declined
// - uri based, e.g. http://schema/user/not_found.json
//
// Kind must be unique.
func (c *Cause) Kind() string {
	return c.kind
}

func (c *Cause) Error() string {
	return c.msg
}

func (c *Cause) String() string {
	return fmt.Sprintf("%s/%s: %s", c.code, c.kind, c.msg)
}

func (c *Cause) Is(err error) bool {
	var cause *Cause
	ok := errors.As(err, &cause)

	return ok &&
		c.code == cause.code &&
		c.kind == cause.kind
}

type Hint[T any] error

// Detail wraps an error with details.
type Detail[T any] struct {
	error
	t T
}

func (d *Detail[T]) Detail() T {
	return d.t
}

func (d *Detail[T]) Unwrap() error {
	return d.error
}

// WrapDetail wraps an error Hint with detail.
func WrapDetail[T any](err Hint[T], t T) error {
	return &Detail[T]{
		error: err,
		t:     t,
	}
}
