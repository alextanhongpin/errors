package cause_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

var ErrDuplicateRow = cause.New(codes.Exists, "DuplicateRowError", "Duplicate row")
var ErrPaymentFailed = cause.New(codes.Conflict, "PaymentFailedError", "Duplicate payment attempt")

func ExampleError_LogValue() {
	var err error = ErrPaymentFailed.WithDetails(map[string]any{
		"order_id": "12345",
	}).Wrap(ErrDuplicateRow.Wrap(sql.ErrNoRows))

	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			a.Value = slog.TimeValue(time.Date(2025, 6, 7, 0, 52, 24, 115438000, time.UTC))
		}
		return a
	}

	b := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(b, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replacer}))
	logger.Error("payment failed", slog.Any("error", err))

	data := b.Bytes()
	b.Reset()
	if err := json.Indent(b, data, "", "  "); err != nil {
		panic(err)
	}
	fmt.Println(b.String())

	// Output:
	// {
	//   "time": "2025-06-07T00:52:24.115438Z",
	//   "level": "ERROR",
	//   "source": {
	//     "function": "github.com/alextanhongpin/errors/cause_test.ExampleError_LogValue",
	//     "file": "/Users/alextanhongpin/Documents/go/errors/cause/examples_log_value_test.go",
	//     "line": 32
	//   },
	//   "msg": "payment failed",
	//   "error": {
	//     "message": "Duplicate payment attempt",
	//     "code": "conflict",
	//     "name": "PaymentFailedError",
	//     "details": {
	//       "order_id": "12345"
	//     },
	//     "cause": {
	//       "message": "Duplicate row",
	//       "code": "exists",
	//       "name": "DuplicateRowError",
	//       "cause": "sql: no rows in result set"
	//     }
	//   }
	// }
}
