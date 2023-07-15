package stacktrace_test

import (
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func ExampleStackTraceNew() {
	err := baz()

	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: baz: bar: to foo or not: foo
	//     Origin is: foo
	//         at stacktrace_test.foo (in examples_stacktrace_new_test.go:41)
	//     Caused by: bar
	//         at stacktrace_test.bar (in examples_stacktrace_new_test.go:47)
	//     Caused by: baz
	//         at stacktrace_test.baz (in examples_stacktrace_new_test.go:51)
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceNew (in examples_stacktrace_new_test.go:10)
	//
	// Reversed:
	// Error: baz: bar: to foo or not: foo
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceNew (in examples_stacktrace_new_test.go:10)
	//     Caused by: baz
	//         at stacktrace_test.baz (in examples_stacktrace_new_test.go:51)
	//     Caused by: bar
	//         at stacktrace_test.bar (in examples_stacktrace_new_test.go:47)
	//     Origin is: foo
	//         at stacktrace_test.foo (in examples_stacktrace_new_test.go:41)
}

func foo() error {
	err := stacktrace.New("foo")
	err = fmt.Errorf("to foo or not: %w", err)
	return err
}

func bar() error {
	return stacktrace.Annotate(foo(), "bar")
}

func baz() error {
	return stacktrace.Annotate(bar(), "baz")
}
