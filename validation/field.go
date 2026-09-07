package validation

import (
	"fmt"
	"maps"
	"slices"
)

type validatable interface {
	Errors() Errors
}

type Slice[T validatable] []T

func SliceOf[T validatable](s []T) Slice[T] {
	return Slice[T](s)
}

func (ss Slice[T]) Errors() Errors {
	field := make(map[string][]string)
	for i, s := range ss {
		for k, v := range s.Errors() {
			field[fmt.Sprintf("[%d].%s", i, k)] = v
		}
	}
	return field
}

type Errors map[string][]string

func (e Errors) If(key string, value bool, val string) {
	if !value {
		return
	}
	e[key] = append(e[key], val)
}

func (e Errors) Optional[T validatable](key string, val T) {
	if IsNilOrZero(val) {
		return
	}

	e.nest(key, val)
}

func (e Errors) Required[T validatable](key string, val T) {
	if IsNilOrZero(val) {
		e[key] = append(e[key], "required")
		return
	}

	e.nest(key, val)
}

func (e Errors) nest[T validatable](key string, val T) {
	for k, v := range val.Errors() {
		if k[0] == '[' {
			k = fmt.Sprintf("%s%s", key, k)
		} else {
			k = fmt.Sprintf("%s.%s", key, k)
		}
		e[k] = append(e[k], v...)
	}
}

func (e Errors) Add(key string, vals string) {
	e[key] = append(e[key], vals)
}

func (e Errors) Error() error {
	keys := slices.Sorted(maps.Keys(e))
	errs := make(errorMap)
	for _, key := range keys {
		vals := filter(e[key])
		if len(vals) == 0 {
			continue
		}
		errs[key] = vals
	}
	if len(errs) == 0 {
		return nil
	}

	return errs
}
