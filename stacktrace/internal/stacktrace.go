package internal

import (
	"errors"
	"fmt"
	"runtime"
)

// MaxDepth is configurable.
var MaxDepth = 32

func New(msg string, args ...any) error {
	return newCaller(2, msg, args...)
}

func Wrap(err error) error {
	return wrapCaller(2, err)
}

func Annotate(err error, cause string) error {
	if err == nil {
		return nil
	}

	// Skips [Annotate, caller]
	return annotateCaller(2, err, cause)
}

func Unwrap(err error) ([]uintptr, map[uintptr]string) {
	if err == nil {
		return nil, nil
	}

	return unwrap(err)
}

func newCaller(skip int, msg string, args ...any) error {
	stack := callers(skip + 1)
	pc, _ := head(stack)

	return &ErrorTrace{
		err:   fmt.Errorf(msg, args...),
		stack: stack, // Skips [New, caller]
		cause: fmt.Sprintf(msg, args...),
		pc:    pc,
	}
}

func wrapCaller(skip int, err error) error {
	if err == nil {
		return nil
	}

	var t *ErrorTrace
	if errors.As(err, &t) {
		return t
	}

	stack := callers(skip + 1)
	pc, _ := head(stack)

	return &ErrorTrace{
		err:   err,
		stack: stack, // Skips [Wrap, caller]
		cause: err.Error(),
		pc:    pc,
	}
}

func annotateCaller(skip int, err error, cause string) *ErrorTrace {
	if err == nil {
		return nil
	}

	pcs, _ := unwrap(err)
	seen := make(map[runtime.Frame]uintptr)

	for _, pc := range pcs {
		seen[frameKey(pc)] = pc
	}

	stack := callers(skip + 1)

	// In the rare case where the stack is empty, the cause will not be recorded.
	// cause.
	if len(stack) == 0 {
		return &ErrorTrace{
			err: err,
		}
	}

	// The first element in the stack is the PC where we want to annotateCaller the
	// cause.
	// If may already exists in previous frames.
	pc, ok := seen[frameKey(stack[0])]
	if !ok {
		pc = stack[0]
	}

	var count int
	for i, pc := range stack {
		// Only record frames that we have not seen.
		if _, ok := seen[frameKey(pc)]; !ok {
			stack = stack[i:]
			break
		}

		// Track the seen count.
		count++
	}

	// We have seen all the frames, clear it to avoid duplicate stacktrace.
	if count == len(stack) {
		stack = nil
	}

	return &ErrorTrace{
		err:   err,
		stack: stack,
		cause: cause,
		pc:    pc,
	}
}

type ErrorTrace struct {
	err   error
	stack []uintptr

	// The annotated cause at specific program line.
	cause string

	// The PC containing the cause, it can be from previous errors.
	pc uintptr
}

func (e *ErrorTrace) StackTrace() []uintptr {
	pcs := make([]uintptr, len(e.stack))
	copy(pcs, e.stack)
	return pcs
}

func (e *ErrorTrace) Error() string {
	// Wrap the cause. This should be the same behaviour as
	// github.com/pkg/errors.
	if len(e.cause) > 0 && e.cause != e.err.Error() {
		return fmt.Sprintf("%s: %s", e.cause, e.err.Error())
	}

	return e.err.Error()
}

func (e *ErrorTrace) Unwrap() error {
	return e.err
}

func Reverse[T any](s []T) {
	reverse(s)
}

func unwrap(err error) ([]uintptr, map[uintptr]string) {
	if err == nil {
		return nil, nil
	}

	var pcs []uintptr
	cause := make(map[uintptr]string)
	seen := make(map[runtime.Frame]bool)

	for err != nil {
		var t *ErrorTrace
		if !errors.As(err, &t) {
			break
		}

		// Set the frame with the cause.
		if t.pc != 0 && len(t.cause) > 0 {
			cause[t.pc] = t.cause
		}

		var ordered []uintptr
		frames := runtime.CallersFrames(t.StackTrace())
		for {
			f, more := frames.Next()
			if f.Function == "" {
				break
			}

			key := runtime.Frame{
				File:     f.File,
				Function: f.Function,
				Line:     f.Line,
			}
			if seen[key] {
				break
			}

			seen[key] = true
			// The runtime.CallersFrames PC =
			// runtime.callers(skip) PC - 1
			ordered = append(ordered, f.PC+1)

			if !more {
				break
			}
		}

		// Stack is ordered from bottom-up.
		// Reverse it so that it goes top-down.
		reverse(ordered)

		pcs = append(pcs, ordered...)
		err = t.Unwrap()
	}

	// Return in the order as what the original
	// runtime.callers will return, which is bottom-up.
	reverse(pcs)

	return pcs, cause
}

func callers(skip int) []uintptr {
	pcs := make([]uintptr, MaxDepth)
	// skip [runtime.callers, callers]
	n := runtime.Callers(skip+2, pcs)
	if n == 0 {
		return nil
	}

	pcs = pcs[:n]
	return pcs
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func frameKey(pc uintptr) runtime.Frame {
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return runtime.Frame{
		File:     f.File,
		Function: f.Function,
		Line:     f.Line,
	}
}

func head[T any](ts []T) (t T, ok bool) {
	if len(ts) > 0 {
		return ts[0], true
	}

	return
}
