package cause_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/alextanhongpin/errors/cause"
)

func TestCause(t *testing.T) {
	original := cause.ErrUnknown.WithCause(errors.ErrUnsupported)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var restored *cause.Error
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(restored, cause.ErrUnknown) || !errors.Is(restored, errors.ErrUnsupported) {
		t.Fatal("round trip lost comparison")
	}
}
