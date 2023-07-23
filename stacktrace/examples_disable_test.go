package stacktrace_test

import (
	"database/sql"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

// Skip []
var straceToggle = stacktrace.Caller(0)

func ExampleDisable() {
	err := handleStackTrace()
	fmt.Println(stacktrace.Sprint(err))
	fmt.Println()

	// Disable stacktrace.
	straceToggle.Disable()

	err = handleStackTrace()
	fmt.Println(stacktrace.Sprint(err))

	// Output:
	// Error: stacktrace is enabled: sql: no rows in result set
	//     Origin is: stacktrace is enabled
	//         at stacktrace_test.handleStackTrace (in examples_disable_test.go:35)
	//     Ends here:
	//         at stacktrace_test.ExampleDisable (in examples_disable_test.go:14)
	//
	// Error: stacktrace is enabled: sql: no rows in result set
}

func handleStackTrace() error {
	err := straceToggle.Annotate(sql.ErrNoRows, `stacktrace is enabled`)
	return err
}
