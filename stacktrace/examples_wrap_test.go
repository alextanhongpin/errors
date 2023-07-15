package stacktrace_test

import (
	"database/sql"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func ExampleWrap() {
	err := handlePayout()
	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: product "lego" not found: sql: no rows in result set
	//     Origin is: product "lego" not found
	//         at stacktrace_test.handlePayout (in examples_wrap_test.go:33)
	//     Ends here:
	//         at stacktrace_test.ExampleWrap (in examples_wrap_test.go:11)
	//
	// Reversed:
	// Error: product "lego" not found: sql: no rows in result set
	//     Ends here:
	//         at stacktrace_test.ExampleWrap (in examples_wrap_test.go:11)
	//     Origin is: product "lego" not found
	//         at stacktrace_test.handlePayout (in examples_wrap_test.go:33)
}

func handlePayout() error {
	err := stacktrace.Annotate(sql.ErrNoRows, `product "lego" not found`)
	return err
}
