package cause

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type validatable interface {
	Validate() error
}

type validateMany interface {
	Validate() map[int]error
}

type Fields map[string]any

func (ve Fields) AsError() error {
	if len(ve) == 0 {
		return nil
	}

	fe := make(fieldError)
	maps.Copy(fe, ve)
	return fe
}

func (ve Fields) Optional(key string, value any, other ...string) Fields {
	return ve.validate(key, value, false, other...)
}

func (ve Fields) Required(key string, value any, other ...string) Fields {
	return ve.validate(key, value, true, other...)
}

func (ve Fields) validate(key string, value any, required bool, other ...string) Fields {
	if isZero(value) {
		if required {
			ve[key] = "required"
		}

		return ve
	}

	if v, ok := value.(validatable); ok {
		if len(other) > 0 {
			panic("cannot use validatable with other messages")
		}

		ve[key] = v.Validate()
	} else if v, ok := value.(validateMany); ok {
		if len(other) > 0 {
			panic("cannot use validateMany with other messages")
		}

		for i, err := range v.Validate() {
			if err != nil {
				ve[fmt.Sprintf("%s[%d]", key, i)] = err
			}
		}
	} else {
		if err := stringSliceError(other).Validate(); err != nil {
			ve[key] = err
		}
	}

	return ve
}

func Cond(valid bool, msg string) string {
	if valid {
		return msg
	}

	return ""
}

type validateFunc[T any] struct {
	val T
	fn  func(T) error
}

func (vf *validateFunc[T]) Validate() error {
	if isZero(vf.val) {
		return nil
	}

	return vf.fn(vf.val)
}

func Slice[T comparable](s []T, fn func(T) error) validateMany {
	if len(s) == 0 {
		return nil
	}

	res := make([]validatable, len(s))
	for i, item := range s {
		res[i] = &validateFunc[T]{val: item, fn: fn}
	}

	return Collect(res)
}

func Collect[T validatable](s []T) validateMany {
	if len(s) == 0 {
		return nil
	}

	ei := make(errorIndex)

	for i, item := range s {
		if err := item.Validate(); err != nil {
			var fe fieldError
			if errors.As(err, &fe) {
				ei[i] = fe
			} else {
				ei[i] = stringSliceError{err.Error()}
			}
		}
	}

	return ei
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

type fieldError map[string]any

func (fe fieldError) Error() string {
	keys := slices.Sorted(maps.Keys(fe))
	return fmt.Sprintf("invalid fields: %s", strings.Join(keys, ", "))
}

type stringSliceError []string

func (s stringSliceError) Validate() error {
	var res []string
	for _, v := range s {
		if v != "" {
			res = append(res, v)
		}
	}
	if len(res) == 0 {
		return nil
	}

	return stringSliceError(res)
}

func (s stringSliceError) Error() string {
	return strings.Join(s, ", ")
}

func (s stringSliceError) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(s.Error())), nil
}

type errorIndex map[int]error

func (ei errorIndex) Validate() map[int]error {
	return ei
}
