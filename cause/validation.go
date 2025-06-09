package cause

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

type validatable interface {
	Validate() error
}

type MapError map[string]any

func (me MapError) Error() string {
	return "invalid fields: " + me.String()
}

func (me MapError) String() string {
	return strings.Join(slices.Sorted(maps.Keys(me)), ", ")
}

type MapValidator map[string]any

func NewMapValidator() MapValidator {
	return make(map[string]any)
}

func (m MapValidator) Optional(name string, val any, other ...string) MapValidator {
	return m.Add(name, Optional(val, other...))
}

func (m MapValidator) Required(name string, val any, other ...string) MapValidator {
	return m.Add(name, Required(val, other...))
}

func (v MapValidator) Add(key string, validator validatable) MapValidator {
	if err := validator.Validate(); err != nil {
		switch e := err.(type) {
		case MapError:
			v[key] = e
		case SliceError:
			for i, val := range e {
				v[fmt.Sprintf("%s[%d]", key, i)] = val
			}
		case error:
			v[key] = e.Error()
		default:
			v[key] = e
		}
	}

	return v
}

func (v MapValidator) Validate() error {
	if len(v) == 0 {
		return nil
	}

	me := make(MapError)
	maps.Copy(me, v)

	return me
}

type SliceError map[int]any

func (se SliceError) Error() string {
	return "invalid slice"
}

type SliceValidator []validatable

func Slice[T validatable](vs []T) SliceValidator {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = v
	}

	return SliceValidator(res)
}

func SliceFunc[T any](vs []T, fn func(T) error) SliceValidator {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = &ValueValidator[T]{
			value:    v,
			validate: fn,
		}
	}

	return SliceValidator(res)
}

func (s SliceValidator) Validate() error {
	se := make(SliceError)

	for i, item := range s {
		if err := item.Validate(); err != nil {
			switch e := err.(type) {
			case MapError, SliceError:
				se[i] = e
			case error:
				se[i] = e.Error()
			default:
				se[i] = err
			}
		}
	}

	if len(se) == 0 {
		return nil
	}

	return se
}

type ValueError string

func (ve ValueError) Error() string {
	return string(ve)
}

type ValueValidator[T any] struct {
	value    T
	validate func(T) error
}

func (v *ValueValidator[T]) Validate() error {
	return v.validate(v.value)
}

func Required(val any, other ...string) validatable {
	return &ValueValidator[any]{
		value:    val,
		validate: validateFunc(true, other...),
	}
}

func Optional(val any, other ...string) validatable {
	return &ValueValidator[any]{
		value:    val,
		validate: validateFunc(false, other...),
	}
}

func validateFunc(required bool, other ...string) func(any) error {
	return func(v any) error {
		if isZero(v) {
			if required {
				return ValueError("required")
			}

			return nil
		}

		switch v := v.(type) {
		case validatable:
			return v.Validate()
		default:
			other = filterZero(other)
			if len(other) == 0 {
				return nil
			}

			return ValueError(strings.Join(other, ", "))
		}
	}
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

func When(valid bool, msg string) string {
	if valid {
		return msg
	}

	return ""
}
