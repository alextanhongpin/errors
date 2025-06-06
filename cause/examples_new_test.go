package cause_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

var ErrUserNotFound = cause.New(codes.NotFound, "UserNotFoundError", "User not found")

func ExampleNew() {
	var err error = ErrUserNotFound
	fmt.Println("err:", err)
	fmt.Println("is:", errors.Is(err, ErrUserNotFound))

	var causeErr *cause.Error
	if errors.As(err, &causeErr) {
		fmt.Println("code:", causeErr.Code)
		fmt.Println("details:", causeErr.Details)
		fmt.Println("message:", causeErr.Message)
		fmt.Println("name:", causeErr.Name)
	}

	// Output:
	// err: User not found
	// is: true
	// code: not_found
	// details: map[]
	// message: User not found
	// name: UserNotFoundError
}
