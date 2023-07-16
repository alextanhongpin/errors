package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

var ErrUserNotFound = errors.New("user not found")

func ExampleStackTraceWrap() {
	err := findUser()

	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: user not found
	//     Origin is: user not found
	//         at stacktrace_test.findUser (in examples_stacktrace_wrap_test.go:36)
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceWrap (in examples_stacktrace_wrap_test.go:13)
	//
	// Reversed:
	// Error: user not found
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceWrap (in examples_stacktrace_wrap_test.go:13)
	//     Origin is: user not found
	//         at stacktrace_test.findUser (in examples_stacktrace_wrap_test.go:36)
}

func findUser() error {
	err := stacktrace.Wrap(ErrUserNotFound)
	return err
}
