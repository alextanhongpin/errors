package internal

import (
	"fmt"
	"sync/atomic"
)

func Caller(skip int) stacktraceToggler {
	t := new(atomic.Bool)
	t.Store(true)

	return &stackTraceToggler{
		// Skip caller.[New, Wrap, Annotate]
		caller:  &caller{skip: skip + 1},
		null:    &null{},
		enabled: t,
	}
}

type stackTraceToggler struct {
	caller  *caller
	null    *null
	enabled *atomic.Bool
}

func (st *stackTraceToggler) Enable() {
	st.enabled.Swap(true)
}

func (st *stackTraceToggler) Disable() {
	st.enabled.Swap(false)
}

func (st *stackTraceToggler) New(msg string, args ...any) error {
	if st.enabled.Load() {
		return st.caller.New(msg, args...)
	}

	return st.null.New(msg, args...)
}

func (st *stackTraceToggler) Wrap(err error) error {
	if st.enabled.Load() {
		return st.caller.Wrap(err)
	}

	return st.null.Wrap(err)
}

func (st *stackTraceToggler) Annotate(err error, cause string) error {
	if st.enabled.Load() {
		return st.caller.Annotate(err, cause)
	}

	return st.null.Annotate(err, cause)
}

type stacktraceToggler interface {
	toggler
	stacktrace
}

type toggler interface {
	Enable()
	Disable()
}

type stacktrace interface {
	New(msg string, args ...any) error
	Wrap(error) error
	Annotate(err error, cause string) error
}

type caller struct {
	skip int
}

func (c *caller) New(msg string, args ...any) error {
	// Skip [New]
	return newCaller(c.skip+1, msg, args...)
}

func (c *caller) Wrap(err error) error {
	// Skip [Wrap]
	return wrapCaller(c.skip+1, err)
}

func (c *caller) Annotate(err error, cause string) error {
	// Skip [Annotate]
	return annotateCaller(c.skip+1, err, cause)
}

func Null() stacktrace {
	return &null{}
}

type null struct {
	skip int
}

func (n *null) New(msg string, args ...any) error {
	return fmt.Errorf(msg, args...)
}

func (n *null) Wrap(err error) error {
	return err
}

func (n *null) Unwrap(err error) ([]uintptr, map[uintptr]string) {
	return nil, nil
}

func (n *null) Annotate(err error, cause string) error {
	return fmt.Errorf("%s: %w", cause, err)
}
