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

var ErrPayoutDeclined = causes.NewHint[PO01](codes.Conflict, "payout/declined", "Payout is declined")

func ExampleNewHint() {
	var err error = ErrPayoutDeclined.Wrap(PayoutDeclinedErrorDetail{
		PayoutID: "PO-42",
		Reason:   "Insufficient balance in account",
	})

	fmt.Println(ErrPayoutDeclined.Is(err))

	var d causes.Detail
	if errors.As(err, &d) {
		fmt.Println(d.Code())
		fmt.Println(d.Kind())
		fmt.Println(d.Message())

		// Data is of type `any`.
		fmt.Printf("%#v\n", d.Data())

		// Assert as type `PayoutDeclinedErrorDetail`.
		t, ok := ErrPayoutDeclined.Unwrap(err)
		fmt.Printf("%#v\n", t)
		fmt.Println(ok)
	}

	// Output:
	// true
	// conflict
	// payout/declined
	// Payout is declined
	// causes_test.PayoutDeclinedErrorDetail{PayoutID:"PO-42", Reason:"Insufficient balance in account"}
	// causes_test.PayoutDeclinedErrorDetail{PayoutID:"PO-42", Reason:"Insufficient balance in account"}
	// true
}
