package cause_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

func ExampleFrom() {
	data, err := json.Marshal(cause.From(sql.ErrNoRows))
	if err != nil {
		panic(err)
	}
	var restored *cause.Error
	if err := json.Unmarshal(data, &restored); err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(errors.Is(restored, sql.ErrNoRows))
	fmt.Println(errors.Is(restored, errors.New(sql.ErrNoRows.Error())))
	// Output:
	// {"code":"unknown","name":"","message":"sql: no rows in result set"}
	// true
	// true
}

func ExampleError_MarshalJSON() {
	original := cause.ErrNotFound.WithCause(sql.ErrNoRows)
	data, err := json.Marshal(original)
	if err != nil {
		panic(err)
	}
	var restored *cause.Error
	if err := json.Unmarshal(data, &restored); err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(errors.Is(restored, original))
	fmt.Println(errors.Is(restored, sql.ErrNoRows))
	// Output:
	// {"code":"not_found","name":"NOT_FOUND","message":"The specified resource was not found","causes":[{"code":"unknown","name":"","message":"sql: no rows in result set"}]}
	// true
	// true
}
