package cause_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

var ErrUnknown = cause.New(codes.Unknown, "Unknown error", "This needs to be fixed")

func ExampleError_marshal() {
	var err error = ErrUnknown.WithCause(sql.ErrNoRows)
	b, err := json.Marshal(err)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	var causeErr *cause.Error
	if err := json.Unmarshal(b, &causeErr); err != nil {
		log.Fatal(err)
	}

	fmt.Println("is sql.ErrNoRows?:", errors.Is(causeErr, sql.ErrNoRows))
	fmt.Println("is ErrUnknown?:", errors.Is(causeErr, ErrUnknown))

	// Output:
	// {"code":17,"message":"This needs to be fixed","name":"Unknown error"}
	// is sql.ErrNoRows?: false
	// is ErrUnknown?: true
}
