package cause_test

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

var ErrStorage = cause.New(codes.Internal, "StorageError", "Storage error")

func ExampleError_Wrap() {
	var err error = ErrStorage.Wrap(sql.ErrNoRows)
	fmt.Println("is sql.ErrNoRows?:", errors.Is(err, sql.ErrNoRows))
	fmt.Println("is ErrStorage?:", errors.Is(err, ErrStorage))

	var causeErr *cause.Error
	if errors.As(err, &causeErr) {
		fmt.Println("cause:", causeErr.Unwrap())
	}

	// Output:
	// is sql.ErrNoRows?: true
	// is ErrStorage?: true
	// cause: sql: no rows in result set
}
