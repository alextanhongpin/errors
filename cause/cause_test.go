package cause_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/alextanhongpin/errors/cause"
)

func TestCause(t *testing.T) {
	err := cause.ErrUnknown.WithCause(errors.ErrUnsupported)
	b, mErr := json.Marshal(err)
	if mErr != nil {
		t.Fatal(mErr)
	}
	var c *cause.Error
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(c, cause.ErrUnknown) {
		t.Fatal("want err unknown")
	}

	if !errors.Is(c, errors.ErrUnsupported) {
		t.Fatal("want err unsupported")
	}
}

type UserError struct {
	Code    string
	Message string
}

func (e *UserError) Is(err error) bool {
	cause, ok := errors.AsType[*UserError](err)
	if !ok {
		return false
	}
	return e.Message == cause.Message && e.Code == cause.Code
}
func (e *UserError) Error() string {
	return e.Message
}
