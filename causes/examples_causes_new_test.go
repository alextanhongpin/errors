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

	var d causes.Detail
	if errors.As(err, &d) {
		fmt.Println(d.Code())
		fmt.Println(d.Kind())
		fmt.Println(d.Message())
		fmt.Println(d.Data())
	}

	// Output:
	// true
	// conflict
	// payout/frozen
	// Your payout is frozen due to suspicious transactions.
	// <nil>
}
