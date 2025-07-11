package stacktrace_test

import (
	"database/sql"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func Example_annotateMore() {
	err := findProduct(42)
	err = stacktrace.Annotate(err, "failed to find product")

	fmt.Println(stacktrace.Sprint(err))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.SprintReversed(err))

	// Output:
	// Error: failed to find product: product id "42": sql: no rows in result set
	//     Origin is: product id "42"
	//         at stacktrace_test.findProduct (in examples_annotate_more_test.go:37)
	//         at stacktrace_test.Example_annotateMore (in examples_annotate_more_test.go:11)
	//     Ends here: failed to find product
	//         at stacktrace_test.Example_annotateMore (in examples_annotate_more_test.go:12)
	//
	// Reversed:
	// Error: failed to find product: product id "42": sql: no rows in result set
	//     Ends here: failed to find product
	//         at stacktrace_test.Example_annotateMore (in examples_annotate_more_test.go:12)
	//         at stacktrace_test.Example_annotateMore (in examples_annotate_more_test.go:11)
	//     Origin is: product id "42"
	//         at stacktrace_test.findProduct (in examples_annotate_more_test.go:37)
}

func findProduct(id int64) error {
	err := stacktrace.Annotate(sql.ErrNoRows, `product id "%d"`, id)
	return err
}
