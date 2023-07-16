package causes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

type PayoutDeclinedErrorDetail struct {
	PayoutID string
	Reason   string
}

// Use alias to shorten the type detail.
type PO01 = PayoutDeclinedErrorDetail

var ErrPayoutDeclined = causes.NewWithHint[PO01](codes.Conflict, "payout/declined", "Payout is declined")

// Alternative is to define a function to ensure the error is always wrapped.

func PayoutDeclinedError(detail PayoutDeclinedErrorDetail) error {
	return causes.WrapDetail(ErrPayoutDeclined, detail)
}

func ExampleNewWithHint() {
	err := PayoutDeclinedError(PayoutDeclinedErrorDetail{
		PayoutID: "PO-42",
		Reason:   "Insufficient balance in account",
	})

	// Or
	err = causes.WrapDetail(ErrPayoutDeclined, PayoutDeclinedErrorDetail{
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

	var d *causes.Detail[PayoutDeclinedErrorDetail]
	if errors.As(err, &d) {
		fmt.Printf("%#v\n", d.Detail())
	}

	fmt.Println(d.Unwrap() == c)
	fmt.Println(d.Unwrap() == ErrPayoutDeclined)

	// Output:
	// true
	// conflict
	// payout/declined
	// Payout is declined
	// causes_test.PayoutDeclinedErrorDetail{PayoutID:"PO-42", Reason:"Insufficient balance in account"}
	// true
	// true
}
