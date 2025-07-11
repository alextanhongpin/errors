//go:build go1.18
// +build go1.18

package validator

import (
	"errors"
	"testing"
)

func FuzzIsZero_Int(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Fuzz(func(t *testing.T, v int) {
		_ = isZero(v)
	})
}

func FuzzIsZero_String(f *testing.F) {
	f.Add("")
	f.Add("foo")
	f.Fuzz(func(t *testing.T, v string) {
		_ = isZero(v)
	})
}

func FuzzRequired_Int(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Fuzz(func(t *testing.T, v int) {
		_ = Required(v, errors.New("fail"))
	})
}

func FuzzRequired_String(f *testing.F) {
	f.Add("")
	f.Add("foo")
	f.Fuzz(func(t *testing.T, v string) {
		_ = Required(v, errors.New("fail"))
	})
}

func FuzzOptional_Int(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Fuzz(func(t *testing.T, v int) {
		_ = Optional(v, errors.New("fail"))
	})
}

func FuzzOptional_String(f *testing.F) {
	f.Add("")
	f.Add("foo")
	f.Fuzz(func(t *testing.T, v string) {
		_ = Optional(v, errors.New("fail"))
	})
}

func FuzzValidateManyFunc_IntSlice(f *testing.F) {
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, b []byte) {
		items := make([]int, len(b))
		for i, v := range b {
			items[i] = int(v)
		}
		_ = ValidateManyFunc(items, func(i int) error {
			if i%2 == 0 {
				return errors.New("even")
			}
			return nil
		})
	})
}
