package causes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

var ErrPayoutFrozen = causes.New(codes.Conflict, "payout/frozen", "Your payout is frozen due to suspicious transactions.")

func ExampleNew() {
	var err error = ErrPayoutFrozen
	fmt.Println(errors.Is(err, ErrPayoutFrozen))

	var c *causes.Cause
	if errors.As(err, &c) {
		fmt.Println(c.Code())
		fmt.Println(c.Kind())
		fmt.Println(c.Error())
	}

	// Output:
	// true
	// conflict
	// payout/frozen
	// Your payout is frozen due to suspicious transactions.
}
