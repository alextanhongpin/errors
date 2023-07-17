package causes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

var ErrProductNotFound = causes.New(codes.NotFound, "product/not_found", "The product is not found")

func ExampleCausesConvert() {
	var err error = ErrProductNotFound
	fmt.Println(errors.Is(err, ErrProductNotFound))

	var c causes.Detail
	if errors.As(err, &c) {
		d := c.Detail()
		fmt.Println(d.Code())
		fmt.Println(d.Kind())
		fmt.Println(d.Message())
		fmt.Println(d.Data())
		fmt.Println(codes.HTTP(d.Code()))
	}

	// Output:
	// true
	// not_found
	// product/not_found
	// The product is not found
	// <nil>
	// 404
}
