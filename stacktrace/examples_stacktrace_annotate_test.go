package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func ExampleStackTraceAnnotate() {
	err := findProduct()
	err = stacktrace.Annotate(err, "product-123")

	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: product not found
	//     Origin is: findProduct
	//         at stacktrace_test.findProduct (in examples_stacktrace_annotate_test.go:37)
	//         at stacktrace_test.ExampleStackTraceAnnotate (in examples_stacktrace_annotate_test.go:11)
	//     Ends here: product-123
	//         at stacktrace_test.ExampleStackTraceAnnotate (in examples_stacktrace_annotate_test.go:12)
	//
	// Reversed:
	// Error: product not found
	//     Ends here: product-123
	//         at stacktrace_test.ExampleStackTraceAnnotate (in examples_stacktrace_annotate_test.go:12)
	//         at stacktrace_test.ExampleStackTraceAnnotate (in examples_stacktrace_annotate_test.go:11)
	//     Origin is: findProduct
	//         at stacktrace_test.findProduct (in examples_stacktrace_annotate_test.go:37)
}

func findProduct() error {
	err := stacktrace.Annotate(errors.New("product not found"), "findProduct")
	return err
}
