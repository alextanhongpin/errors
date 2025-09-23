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

var (
	ErrUnknown = cause.New(codes.Unknown, "Unknown error", "This needs to be fixed")
	ErrNested  = cause.New(codes.Unknown, "Nested error", "One level of nesting")
)

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
	// {"cause":{"code":17,"message":"sql: no rows in result set","name":"Unknown"},"code":17,"message":"This needs to be fixed","name":"Unknown error"}
	// is sql.ErrNoRows?: true
	// is ErrUnknown?: true
}

func ExampleError_marshal_nested() {
	var err error = ErrUnknown.WithCause(ErrNested.WithCause(sql.ErrNoRows))
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
	fmt.Println("is ErrNested?:", errors.Is(causeErr, ErrNested))

	// Output:
	// {"cause":{"cause":{"code":17,"message":"sql: no rows in result set","name":"Unknown"},"code":17,"message":"One level of nesting","name":"Nested error"},"code":17,"message":"This needs to be fixed","name":"Unknown error"}
	// is sql.ErrNoRows?: true
	// is ErrUnknown?: true
	// is ErrNested?: true
}
