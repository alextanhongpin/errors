package causes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

type PayoutDeclinedDetail struct {
	PayoutID string
	Reason   string
}

var ErrPayoutDeclined = causes.NewWithHint[PayoutDeclinedDetail](codes.Conflict, "payout/declined", "Payout is declined")

func ExampleNewWithHint() {
	err := causes.WrapDetail(ErrPayoutDeclined, PayoutDeclinedDetail{
		PayoutID: "PO-42",
		Reason:   "Insufficient balance in account",
	})

	fmt.Println(errors.Is(err, ErrPayoutDeclined))

	var c *causes.Cause
	if errors.As(err, &c) {
		fmt.Println(c.Code())
		fmt.Println(c.Kind())
		fmt.Println(c.Error())
	}

	var d *causes.Detail[PayoutDeclinedDetail]
	if errors.As(err, &d) {
		fmt.Printf("%#v", d.Detail())
	}

	// Output:
	// true
	// conflict
	// payout/declined
	// Payout is declined
	// causes_test.PayoutDeclinedDetail{PayoutID:"PO-42", Reason:"Insufficient balance in account"}
}
