package stacktrace_test

import (
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func dive(depth int) error {
	var do func(n int) error
	do = func(n int) error {
		if n == 0 {
			return stacktrace.New("depth %d", n)
		}

		// Unfortunately due to the current implementation, the same line of code
		// is treated as duplicate.
		switch n {
		case 1:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 2:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 3:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 4:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 5:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 6:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 7:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 8:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 9:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 10:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 11:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		case 12:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		default:
			return stacktrace.Annotate(do(n-1), "at depth %d", n-1)
		}
	}

	return do(depth)
}

func ExampleSetMaxDepth() {
	stacktrace.SetMaxDepth(4)
	defer stacktrace.SetMaxDepth(32)

	err := dive(12)
	fmt.Println(stacktrace.Sprint(err, false))

	// Output:
	// Error: at depth 11: at depth 10: at depth 9: at depth 8: at depth 7: at depth 6: at depth 5: at depth 4: at depth 3: at depth 2: at depth 1: at depth 0: depth 0
	//     Origin is: depth 0
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:13)
	//     Caused by: at depth 0
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:20)
	//     Caused by: at depth 1
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:22)
	//     Caused by: at depth 2
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:24)
	//     Caused by: at depth 3
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:26)
	//     Caused by: at depth 4
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:28)
	//     Caused by: at depth 5
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:30)
	//     Caused by: at depth 6
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:32)
	//     Caused by: at depth 7
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:34)
	//     Caused by: at depth 8
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:36)
	//     Caused by: at depth 9
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:38)
	//     Caused by: at depth 10
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:40)
	//     Caused by: at depth 11
	//         at stacktrace_test.dive.func1 (in examples_set_max_depth_test.go:42)
	//         at stacktrace_test.dive (in examples_set_max_depth_test.go:48)
	//     Ends here:
	//         at stacktrace_test.ExampleSetMaxDepth (in examples_set_max_depth_test.go:55)
}
