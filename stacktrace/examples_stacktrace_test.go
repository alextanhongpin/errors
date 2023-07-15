package stacktrace_test

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/stacktrace"
)

func root() error {
	err := stacktrace.New("root")
	return err
}

func child() error {
	err := stacktrace.Annotate(root(), "child")
	return err
}

func ExampleStackTrace() {
	err := child()
	b, err := json.MarshalIndent(stacktrace.StackTrace(err), "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	// [
	//  {
	//   "id": 1,
	//   "cause": "",
	//   "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/errors/stacktrace/examples_stacktrace_test.go",
	//   "line": 11,
	//   "function": "github.com/alextanhongpin/errors/stacktrace_test.root"
	//  },
	//  {
	//   "id": 2,
	//   "cause": "child",
	//   "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/errors/stacktrace/examples_stacktrace_test.go",
	//   "line": 16,
	//   "function": "github.com/alextanhongpin/errors/stacktrace_test.child"
	//  },
	//  {
	//   "id": 3,
	//   "cause": "",
	//   "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/errors/stacktrace/examples_stacktrace_test.go",
	//   "line": 21,
	//   "function": "github.com/alextanhongpin/errors/stacktrace_test.ExampleStackTrace"
	//  }
	// ]
}