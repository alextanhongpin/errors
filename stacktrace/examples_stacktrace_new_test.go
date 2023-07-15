package stacktrace_test

import (
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func ExampleStackTraceNew() {
	err := baz()
	err = fmt.Errorf("wrapped: %w", err)

	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: wrapped: foo
	//     Origin is:
	//         at stacktrace_test.foo (in examples_stacktrace_new_test.go:42)
	//     Caused by: bar
	//         at stacktrace_test.bar (in examples_stacktrace_new_test.go:47)
	//     Caused by: baz
	//         at stacktrace_test.baz (in examples_stacktrace_new_test.go:51)
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceNew (in examples_stacktrace_new_test.go:10)
	//
	// Reversed:
	// Error: wrapped: foo
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceNew (in examples_stacktrace_new_test.go:10)
	//     Caused by: baz
	//         at stacktrace_test.baz (in examples_stacktrace_new_test.go:51)
	//     Caused by: bar
	//         at stacktrace_test.bar (in examples_stacktrace_new_test.go:47)
	//     Origin is:
	//         at stacktrace_test.foo (in examples_stacktrace_new_test.go:42)
}

func foo() error {
	err := stacktrace.New("foo")
	return err
}

func bar() error {
	return stacktrace.Annotate(foo(), "bar")
}

func baz() error {
	return stacktrace.Annotate(bar(), "baz")
}
