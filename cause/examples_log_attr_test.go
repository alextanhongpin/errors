package cause_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

func ExampleError_log_attr() {
	var err error = cause.New(codes.BadRequest, "BadRequest", "email=%s is invalid", "xyz@mail.com").
		WithAttrs(slog.String("email", "xyz@mail.com"))

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
	//     "function": "github.com/alextanhongpin/errors/cause_test.ExampleError_log_attr",
	//     "file": "/Users/alextanhongpin/Documents/go/errors/cause/examples_log_attr_test.go",
	//     "line": 27
	//   },
	//   "msg": "payment failed",
	//   "error": {
	//     "message": "email=xyz@mail.com is invalid",
	//     "code": "bad_request",
	//     "name": "BadRequest",
	//     "email": "xyz@mail.com"
	//   }
	// }
}
