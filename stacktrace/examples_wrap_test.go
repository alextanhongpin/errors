package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

var ErrPayoutDeclined = errors.New("Payout cannot be processed. Please contact customer support for more information")

func ExampleWrap() {
	err := handlePayout()
	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Origin is: account is actually frozen
	//         at stacktrace_test.handlePayout (in examples_wrap_test.go:35)
	//     Ends here:
	//         at stacktrace_test.ExampleWrap (in examples_wrap_test.go:13)
	//
	// Reversed:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Ends here:
	//         at stacktrace_test.ExampleWrap (in examples_wrap_test.go:13)
	//     Origin is: account is actually frozen
	//         at stacktrace_test.handlePayout (in examples_wrap_test.go:35)
}

func handlePayout() error {
	err := stacktrace.Annotate(ErrPayoutDeclined, "account is actually frozen")
	return err
}
