package cause_test

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

func ExampleError_WithStack() {
	err := bar()
	fmt.Println(err)
	fmt.Println()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		}}))
	logger.Error("An error occurred", "error", err)

	// Output:
	// An error in bar
	// 	at /Users/alextanhongpin/Documents/go/errors/cause/examples_stack_trace_test.go.46
	// Caused by: An error with stack trace
	// 	at /Users/alextanhongpin/Documents/go/errors/cause/examples_stack_trace_test.go.41
	// Caused by: sql: no rows in result set
	//
	// {"level":"ERROR","source":{"function":"github.com/alextanhongpin/errors/cause_test.ExampleError_WithStack","file":"/Users/alextanhongpin/Documents/go/errors/cause/examples_stack_trace_test.go","line":26},"msg":"An error occurred","error":{"message":"An error in bar","code":"internal","name":"BarError","cause":{"message":"An error with stack trace","code":"internal","name":"StackError","cause":"sql: no rows in result set"}}}
}

func foo() error {
	return cause.New(codes.Internal, "StackError", "An error with stack trace").
		WithCause(sql.ErrNoRows).
		WithStack()
}

func bar() error {
	return cause.New(codes.Internal, "BarError", "An error in bar").
		WithStack().
		WithCause(foo())
}
