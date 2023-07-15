package causes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

var ErrUserNotFound = causes.New(codes.NotFound, "user/not_found", "The user does not exists or may have been deleted")

type UserNotFoundError struct {
	error
	UserID string
}

func NewUserNotFoundError(userID string) error {
	return &UserNotFoundError{
		error:  ErrUserNotFound,
		UserID: userID,
	}
}

// Unwrap enables comparison with errors.Is.
func (e *UserNotFoundError) Unwrap() error {
	return e.error
}

// Error is an customised error message.
func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("%s: %q", e.Error(), e.UserID)
}

func ExampleNewCustom() {
	var err error = NewUserNotFoundError("user-3173")
	fmt.Println(errors.Is(err, ErrUserNotFound))

	var c *causes.Cause
	if errors.As(err, &c) {
		fmt.Println(c.Code())
		fmt.Println(c.Kind())
		fmt.Println(c.Error())
	}

	var d *UserNotFoundError
	if errors.As(err, &d) {
		fmt.Printf("%s", d.UserID)
	}

	// Output:
	// true
	// not_found
	// user/not_found
	// The user does not exists or may have been deleted
	// user-3173
}
