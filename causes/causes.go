// package causes represents an error cause. The name is choosen to avoid
// conflict with the errors package.
package causes

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/codes"
)

// hint hints that an errorHint should be wrapped with detail before it can be
// promoted to an error.
type hint[T any] interface {
	Is(error) bool
	Wrap(T) *errorDetail
	Unwrap(error) (T, bool)
}

// New returns a new errorDetail.
func New(code codes.Code, kind, msg string, args ...any) *errorDetail {
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

// Detail allows replacing the implementation detail.
type Detail interface {
	Code() codes.Code
	Data() any
	Kind() string
	Message() string
	Unwrap() error
}

type errorDetail struct {
	code codes.Code
	kind string
	msg  string
	data any
	err  error
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

func (c *errorDetail) Wrap(err error) error {
	cp := *c
	cp.err = err
	return &cp
}

func (c *errorDetail) Unwrap() error {
	return c.err
}

func (c *errorDetail) String() string {
	return fmt.Sprintf("%s/%s: %s", c.code, c.kind, c.msg)
}

func (c *errorDetail) Is(err error) bool {
	if errors.Is(c.err, err) {
		return true
	}

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

func (e *errorHint[T]) Wrap(t T) *errorDetail {
	cp := *e.err
	cp.data = t
	return &cp
}

func (e *errorHint[T]) Unwrap(err error) (v T, ok bool) {
	var errDetail *errorDetail
	if errors.As(err, &errDetail) {
		v, ok = errDetail.data.(T)
	}

	return
}
