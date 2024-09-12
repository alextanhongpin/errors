package causes_test

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

var ErrDocumentNotFound = causes.New(codes.NotFound, "document/not_found", "The document does not exists or may have been deleted")

func ExampleWrap() {
	var err error = ErrDocumentNotFound.Wrap(sql.ErrNoRows)
	fmt.Println(errors.Is(err, ErrDocumentNotFound))
	fmt.Println(errors.Is(err, sql.ErrNoRows))

	var d causes.Detail
	if errors.As(err, &d) {
		fmt.Println(d.Code())
		fmt.Println(d.Kind())
		fmt.Println(d.Message())
		fmt.Println(d.Data())
		fmt.Println(d.Unwrap())
	}

	unwrapErr := errors.Unwrap(err)
	fmt.Println(errors.Is(unwrapErr, ErrDocumentNotFound))
	fmt.Println(errors.Is(unwrapErr, sql.ErrNoRows))

	// Should not affect original.
	fmt.Println(ErrDocumentNotFound.Unwrap() == nil)

	// Output:
	// true
	// true
	// not_found
	// document/not_found
	// The document does not exists or may have been deleted
	// <nil>
	// sql: no rows in result set
	// false
	// true
	// true
}
