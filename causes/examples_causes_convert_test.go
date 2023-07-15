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

	var c *causes.Cause
	if errors.As(err, &c) {
		fmt.Println(c.Code())
		fmt.Println(c.Kind())
		fmt.Println(c.Error())
		fmt.Println(codes.HTTP(c.Code()))
	}

	// Output:
	// true
	// not_found
	// product/not_found
	// The product is not found
	// 404
}
