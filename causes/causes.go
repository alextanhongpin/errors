// package causes represents an error cause. The name is choosen to avoid
// conflict with the errors package.
package causes

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/codes"
)

// Detail needs to be implemented by errors that returns detail.
type Detail interface {
	Detail() detail
}

// Using an interface instead of struct allows for detail to be chained later
// for different operations.
type detail interface {
	Code() codes.Code
	Kind() string
	Message() string
	Data() any
}

// hint hints that an errorHint should be wrapped with detail before it can be
// promoted to an error.
type hint[T any] interface {
	Is(error) bool
	Wrap(T) error
	Unwrap(detail) (T, bool)
}

// New returns a new errorDetail.
func New(code codes.Code, kind, msg string, args ...any) error {
	return &errorDetail{
		code: code,
		kind: kind,
		msg:  fmt.Sprintf(msg, args...),
	}
}

// NewHint returns a partial error that needs to be fulfilled with the hinted
// type.
func NewHint[T any](code codes.Code, kind, msg string, args ...any) hint[T] {
	return &errorHint[T]{
		err: &errorDetail{
			code: code,
			kind: kind,
			msg:  fmt.Sprintf(msg, args...),
		},
	}
}

type errorDetail struct {
	code codes.Code
	kind string
	msg  string
	data any
}

func (c *errorDetail) Detail() detail {
	return c
}

func (c *errorDetail) Code() codes.Code {
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
func (c *errorDetail) Kind() string {
	return c.kind
}

func (c *errorDetail) Error() string {
	return c.msg
}

func (c *errorDetail) Message() string {
	return c.msg
}

func (c *errorDetail) Data() any {
	return c.data
}

func (c *errorDetail) String() string {
	return fmt.Sprintf("%s/%s: %s", c.code, c.kind, c.msg)
}

func (c *errorDetail) Is(err error) bool {
	var cause *errorDetail
	ok := errors.As(err, &cause)

	return ok &&
		c.code == cause.code &&
		c.kind == cause.kind
}

type errorHint[T any] struct {
	err *errorDetail
}

func (e *errorHint[T]) Is(err error) bool {
	return errors.Is(err, e.err)
}

func (e *errorHint[T]) Wrap(t T) error {
	cp := *e.err
	cp.data = t
	return &cp
}

func (e *errorHint[T]) Unwrap(det detail) (T, bool) {
	t, ok := det.Data().(T)
	return t, ok
}
