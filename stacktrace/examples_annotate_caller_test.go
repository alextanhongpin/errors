package stacktrace_test

import (
	"database/sql"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

var strace = stacktrace.Caller(0)

func ExampleCaller() {
	err := handleProduct()
	fmt.Println(stacktrace.Sprint(err))
	fmt.Println()

	fmt.Println("Reversed:")
	fmt.Println(stacktrace.SprintReversed(err))

	// Output:
	// Error: product "lego" not found: sql: no rows in result set
	//     Origin is: product "lego" not found
	//         at stacktrace_test.handleProduct (in examples_annotate_caller_test.go:36)
	//     Ends here:
	//         at stacktrace_test.ExampleCaller (in examples_annotate_caller_test.go:13)
	//
	// Reversed:
	// Error: product "lego" not found: sql: no rows in result set
	//     Ends here:
	//         at stacktrace_test.ExampleCaller (in examples_annotate_caller_test.go:13)
	//     Origin is: product "lego" not found
	//         at stacktrace_test.handleProduct (in examples_annotate_caller_test.go:36)
}

func handleProduct() error {
	err := strace.Annotate(sql.ErrNoRows, `product "lego" not found`)
	return err
}
