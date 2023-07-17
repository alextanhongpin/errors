package causes_test

import (
	"database/sql"
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
	return fmt.Sprintf("%s: %q", e.error.Error(), e.UserID)
}

func ExampleNewCustom() {
	var err error = fmt.Errorf("%w: %w", NewUserNotFoundError("user-3173"), sql.ErrNoRows)
	fmt.Println(errors.Is(err, ErrUserNotFound))

	var c causes.Detail
	if errors.As(err, &c) {
		d := c.Detail()
		fmt.Println(d.Code())
		fmt.Println(d.Kind())
		fmt.Println(d.Message())
		fmt.Println(d.Data())
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
	// <nil>
	// user-3173
}
