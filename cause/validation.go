package cause

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

type validatable interface {
	Validate() error
}

type mapError map[string]any

func (me mapError) Error() string {
	return fmt.Sprintf("invalid fields: %s", strings.Join(slices.Sorted(maps.Keys(me)), ", "))
}

type Map map[string]any

func (m Map) Required(name string, val any, other ...string) Map {
	return m.Add(name, Required(val, other...))
}

func (m Map) Optional(name string, val any, other ...string) Map {
	return m.Add(name, Optional(val, other...))
}

func (m Map) Add(name string, val validatable) Map {
	if isZero(val) {
		return m
	}

	if err := val.Validate(); err != nil {
		switch e := err.(type) {
		case sliceError:
			// If the error is a slice error, we convert it to a map.
			for i, item := range e {
				if isZero(item) {
					continue
				}
				itemName := fmt.Sprintf("%s[%d]", name, i)
				m[itemName] = item
			}
		case mapError:
			// If the error is a map error, we merge it into the current map.
			m[name] = e
		case error:
			m[name] = e.Error()
		default:
			// If the error is not a string, we wrap it in a map.
			m[name] = e
		}
	}
	return m
}

func (m Map) AsError() error {
	me := make(mapError)
	for name, val := range m {
		if isZero(val) {
			continue
		}
		me[name] = val
	}
	if len(me) == 0 {
		return nil
	}
	return me
}

func Optional(value any, msgs ...string) validatable {
	if isZero(value) {
		return &errorStrings{}
	}

	switch v := value.(type) {
	case validatable:
		return v
	default:
		return &errorStrings{msgs}
	}
}

func Required(value any, msgs ...string) validatable {
	if isZero(value) {
		return &errorStrings{[]string{"required"}}
	}
	switch v := value.(type) {
	case validatable:
		return v
	default:
		return &errorStrings{msgs}
	}
}

type errorStrings struct {
	msgs []string
}

func (err *errorStrings) Validate() error {
	msgs := filterZero(err.msgs)
	if len(msgs) == 0 {
		return nil
	}

	return errors.New(strings.Join(msgs, ", "))
}

type sliceError map[int]any

func (se sliceError) Error() string {
	return "SliceError"
}

func When(valid bool, msg string) string {
	if valid {
		return msg
	}

	return ""
}

type SliceError []validatable

func SliceFunc[T any](vs []T, fn func(T) error) SliceError {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = &value[T]{val: v, fn: fn}
	}

	return res

}

func Slice[T validatable](vs []T) SliceError {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = v
	}

	return res
}

func (s SliceError) Validate() error {
	if len(s) == 0 {
		return nil
	}

	se := make(sliceError)

	for i, item := range s {
		if err := item.Validate(); err != nil {
			switch e := err.(type) {
			case sliceError, mapError:
				// If the error is a slice error, we convert it to a map.
				se[i] = e
			case error:
				se[i] = e.Error()
			default:
				// If the error is not a string, we wrap it in a map.
				se[i] = e
			}
		}
	}

	if len(se) == 0 {
		return nil
	}

	return se
}

func filterZero[T comparable](s []T) []T {
	var res []T
	var zero T
	for _, v := range s {
		if v != zero {
			res = append(res, v)
		}
	}
	return res
}

type value[T any] struct {
	val T
	fn  func(T) error
}

func (v *value[T]) Validate() error {
	return v.fn(v.val)
}

func isZero(t any) bool {
	if t == nil {
		return true
	}

	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Slice {
		return v.Len() == 0
	}

	return v.IsZero()
}
