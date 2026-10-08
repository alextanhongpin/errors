package cause_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

type userError struct{ message string }

func (e *userError) Error() string { return e.message }

func ExampleError_WithCause_custom() {
	userErr := &userError{message: "already exists"}
	input := cause.ErrExists.WithCause(userErr)
	fmt.Println("before: matches custom error?", errors.Is(input, userErr))
	data, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	var restored *cause.Error
	err = json.Unmarshal(data, &restored)
	if err != nil {
		panic(err)
	}
	fmt.Println("after: matches original error?", errors.Is(restored, userErr))
	fmt.Println("after: matches ErrExists?", errors.Is(restored, cause.ErrExists))
	fmt.Println("after: matches same text?", errors.Is(restored, &userError{message: "already exists"}))
	// Output:
	// before: matches custom error? true
	// after: matches original error? true
	// after: matches ErrExists? true
	// after: matches same text? true
}

func ExampleError_WithCause_nested() {
	input := cause.ErrAborted.WithCause(cause.ErrBadRequest.WithCause(cause.ErrCanceled.WithCause(errors.ErrUnsupported)))
	data, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	var restored *cause.Error
	err = json.Unmarshal(data, &restored)
	if err != nil {
		panic(err)
	}
	for _, target := range []error{cause.ErrAborted, cause.ErrBadRequest, cause.ErrCanceled, errors.ErrUnsupported} {
		fmt.Println(errors.Is(input, target), errors.Is(restored, target))
	}
	// Output:
	// true true
	// true true
	// true true
	// true true
}
