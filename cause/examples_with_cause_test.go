package cause_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

func ExampleError_WithCause_custom() {
	userErr := &UserError{Code: "duplicate_user_name", Message: "already exists"}
	var err error = userErr
	err = cause.ErrExists.WithCause(err)
	fmt.Println("is &UserError{}?", errors.Is(err, userErr))
	fmt.Println("is cause.ErrExists?", errors.Is(err, cause.ErrExists))
	fmt.Println("is cause.ErrUnknown?", errors.Is(err, cause.ErrUnknown))

	b, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}
	var c *cause.Error
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}

	fmt.Println("is &UserError{}?", errors.Is(c, userErr))
	fmt.Println("is cause.ErrExists?", errors.Is(c, cause.ErrExists))
	fmt.Println("is cause.ErrUnknown?", errors.Is(c, cause.ErrUnknown))

	err = c
	for err != nil {
		fmt.Printf("%T: %v\n", err, err)
		err = errors.Unwrap(err)

		//fmt.Println("is &UserError{}?", errors.Is(c, userErr))
		//fmt.Println("is ErrUnknown?", errors.Is(c, ErrUnknown))
	}

	// Output:
	// is &UserError{}? true
	// is cause.ErrExists? true
	// is cause.ErrUnknown? false
	// is &UserError{}? true
	// is cause.ErrExists? true
	// is cause.ErrUnknown? false
	// *cause.Error: The resource that a client tried to create already exists
	// Caused by: already exists
	// *cause.errorJSON: already exists
}

func ExampleError_WithCause_nested() {
	var err error = cause.ErrAborted.WithCause(cause.ErrBadRequest.WithCause(cause.ErrCanceled.WithCause(errors.ErrUnsupported)))
	fmt.Println("Error:")
	fmt.Println(err)

	fmt.Println("\nBefore marshaling:")
	fmt.Println("is ErrAborted?", errors.Is(err, cause.ErrAborted))
	fmt.Println("is ErrBadRequest?", errors.Is(err, cause.ErrBadRequest))
	fmt.Println("is ErrCanceled?", errors.Is(err, cause.ErrCanceled))
	fmt.Println("is ErrUnsupported?", errors.Is(err, errors.ErrUnsupported))

	b, err := json.MarshalIndent(err, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println("\nMarshal:")
	fmt.Println(string(b))

	var c *cause.Error
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nAfter marshaling:")
	fmt.Println("is ErrAborted?", errors.Is(c, cause.ErrAborted))
	fmt.Println("is ErrBadRequest?", errors.Is(c, cause.ErrBadRequest))
	fmt.Println("is ErrCanceled?", errors.Is(c, cause.ErrCanceled))
	fmt.Println("is ErrUnsupported?", errors.Is(c, errors.ErrUnsupported))

	// Output:
	// Error:
	// The operation was aborted
	// Caused by: The request is invalid
	// Caused by: The operation was canceled
	// Caused by: unsupported operation
	//
	// Before marshaling:
	// is ErrAborted? true
	// is ErrBadRequest? true
	// is ErrCanceled? true
	// is ErrUnsupported? true
	//
	// Marshal:
	// {
	//  "cause": {
	//   "cause": {
	//    "cause": {
	//     "code": 17,
	//     "message": "unsupported operation",
	//     "name": "Unknown"
	//    },
	//    "code": 3,
	//    "message": "The operation was canceled",
	//    "name": "CANCELED"
	//   },
	//   "code": 2,
	//   "message": "The request is invalid",
	//   "name": "BAD_REQUEST"
	//  },
	//  "code": 1,
	//  "message": "The operation was aborted",
	//  "name": "ABORTED"
	// }
	//
	// After marshaling:
	// is ErrAborted? true
	// is ErrBadRequest? true
	// is ErrCanceled? true
	// is ErrUnsupported? true
}
