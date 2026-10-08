package cause_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

func ExampleError_UnmarshalJSON_nested() {
	outer := cause.New(codes.Unknown, "example.outer", "Outer error")
	nested := cause.New(codes.Unknown, "example.nested", "Nested error")
	data, err := json.Marshal(outer.WithCause(nested.WithCause(sql.ErrNoRows)))
	if err != nil {
		panic(err)
	}
	var restored *cause.Error
	err = json.Unmarshal(data, &restored)
	if err != nil {
		panic(err)
	}
	fmt.Println("is sql.ErrNoRows?:", errors.Is(restored, sql.ErrNoRows))
	fmt.Println("is outer?:", errors.Is(restored, outer))
	fmt.Println("is nested?:", errors.Is(restored, nested))
	for node := restored; node != nil; {
		fmt.Printf("%s: %s\n", node.Name, node.Message)
		child, _ := node.Cause.(*cause.Error)
		node = child
	}
	// Output:
	// is sql.ErrNoRows?: true
	// is outer?: true
	// is nested?: true
	// example.outer: Outer error
	// example.nested: Nested error
	// : sql: no rows in result set
}
