package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func ExampleUnwrap() {
	err := stacktrace.Wrap(errors.New("bad request"))

	// Unwrap using errors.As.
	var errTrace *stacktrace.ErrorTrace
	fmt.Println(errors.As(err, &errTrace))

	// Returns the raw unfiltered stacktrace.
	fmt.Println(len(errTrace.StackTrace()))

	// Returned the raw stacktraces, together with cause annotation at
	// specific PCs.
	pcs, causes := stacktrace.Unwrap(err)
	fmt.Println(len(pcs))
	fmt.Println(len(causes))
	fmt.Println(stacktrace.Sprint(err))

	// Output:
	// true
	// 7
	// 7
	// 1
	// Error: bad request
	//         at stacktrace_test.ExampleUnwrap (in examples_unwrap_test.go:11)
}
