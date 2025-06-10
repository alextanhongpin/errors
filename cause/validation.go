package cause

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

var validatableType = reflect.TypeOf((*validatable)(nil)).Elem()

type validatable interface {
	Validate() error
}

type Map map[string]any

func (m Map) Err() error {
	em := make(errorMap)
	for k, v := range m {
		if isZero(v) {
			continue
		}

		switch t := v.(type) {
		case validatable:
			if err := t.Validate(); err != nil {
				var errs = []error{err}
				for len(errs) > 0 {
					err = errs[0]
					errs = errs[1:]

					switch e := err.(type) {
					case errorMulti:
						errs = append(errs, e...)
					case errorIndex:
						em[fmt.Sprintf("%s[%d]", k, e.pos)] = e.err
					default:
						em[k] = e
					}
				}
			}
		default:
			em[k] = v
		}
	}

	for k, v := range em {
		switch e := v.(type) {
		case errorMulti, errorIndex, errorMap:
			// Do nothing, we already have the error in the map.
		case error:
			// If it's a single error, we can just use it directly.
			em[k] = e.Error()
		}
	}
	if len(em) == 0 {
		return nil
	}

	return em
}

func Optional(val any, msgs ...string) *Builder {
	if isZero(val) {
		return nil
	}

	if v, ok := isSlice(val); ok {
		return &Builder{
			msgs: msgs,
			v:    v,
		}
	}

	if v, ok := val.(validatable); ok {
		return &Builder{
			msgs: msgs,
			v:    v,
		}
	}

	return &Builder{
		msgs: msgs,
	}
}

func Required(val any, msgs ...string) *Builder {
	return RequiredMessage(val, "required", msgs...)
}

func RequiredMessage(val any, msg string, msgs ...string) *Builder {
	if isZero(val) {
		return &Builder{
			msgs: []string{msg},
		}
	}

	if v, ok := isSlice(val); ok {
		return &Builder{
			msgs: msgs,
			v:    v,
		}
	}

	if v, ok := val.(validatable); ok {
		return &Builder{
			msgs: msgs,
			v:    v,
		}
	}

	return &Builder{
		msgs: msgs,
	}
}

func SliceFunc[T any](vs []T, fn func(T) error) sliceValidator {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = &validator{
			err: fn(v),
		}
	}

	return res
}

type Builder struct {
	msgs []string
	v    validatable
}

func (b *Builder) When(cond bool, msg string) *Builder {
	if b == nil {
		return nil
	}

	if cond {
		b.msgs = append(b.msgs, msg)
	}

	return b
}

func (b *Builder) Validate() error {
	if b == nil {
		return nil
	}

	var errs errorMulti
	if msg := joinStrings(b.msgs...); msg != "" {
		errs = append(errs, errors.New(msg))
	}

	if b.v != nil {
		if err := b.v.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type errorMap map[string]any

func (e errorMap) Map() map[string]any {
	return e
}

func (e errorMap) Error() string {
	return "invalid fields: " + e.String()
}

func (e errorMap) String() string {
	return strings.Join(slices.Sorted(maps.Keys(e)), ", ")
}

func isSlice(t any) (sliceValidator, bool) {
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Slice {
		result := make([]validatable, v.Len())
		for i := range v.Len() {
			e, ok := v.Index(i).Interface().(validatable)
			if !ok {
				return nil, false
			}
			result[i] = e
		}

		return result, true
	}

	return nil, false
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

func joinStrings(s ...string) string {
	return strings.Join(filterZero(s), ", ")
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

type validator struct {
	err error
}

func (v *validator) Validate() error {
	return v.err
}

type errorIndex struct {
	pos int
	err error
}

func (ei errorIndex) Error() string {
	return fmt.Sprintf("error at index %d: %s", ei.pos, ei.err.Error())
}

type sliceValidator []validatable

func (s sliceValidator) Validate() error {
	var em errorMulti
	for i, v := range s {
		if err := v.Validate(); err != nil {
			em = append(em, errorIndex{i, err})
		}
	}

	return em
}

type errorMulti []error

func (e errorMulti) Error() string {
	return errors.Join(e).Error()
}
