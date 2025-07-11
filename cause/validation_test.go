package cause_test

import (
	"testing"

	"github.com/alextanhongpin/errors/cause"
)

func TestAssertMap(t *testing.T) {
	t.Run("returns nil when no conditions are true", func(t *testing.T) {
		err := cause.AssertMap(map[string]bool{
			"error 1": true,
			"error 2": true,
		})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns error when conditions are false", func(t *testing.T) {
		err := cause.AssertMap(map[string]bool{
			"error 1": false,
			"error 2": true,
			"error 3": false,
		})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "error 1, error 3" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
