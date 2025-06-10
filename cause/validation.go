package cause

import (
	"cmp"
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

type next interface {
	Next() error
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
					case *errorMulti:
						errs = append(errs, e.Errors()...)
					case errorSliceIndex:
						for i, ve := range e {
							em[fmt.Sprintf("%s[%d]", k, i)] = ve
						}
					case errorMap:
						em[k] = e
					case error:
						em[k] = e.Error()
					default:
						em[k] = e
					}
				}
			}
		case error:
			em[k] = t.Error()
		default:
			em[k] = v
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
			v:    v,
			msgs: msgs,
		}
	}

	if v, ok := val.(validatable); ok {
		return &Builder{
			v:    v,
			msgs: msgs,
		}
	}

	return &Builder{
		msgs: msgs,
	}
}

func Required(val any, msgs ...string) *Builder {
	return RequiredMessage(val, "required", msgs...)
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

	em := new(errorMulti)
	if msg := joinStrings(b.msgs...); msg != "" {
		em.err1 = errors.New(msg)
	}

	if b.v != nil {
		em.err2 = b.v.Validate()
	}

	if em.Unwrap() != nil {
		return em
	}

	return nil
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

type errorSliceIndex map[int]any

func (es errorSliceIndex) Error() string {
	return "invalid slice"
}

type sliceValidator []validatable

func (s sliceValidator) Validate() error {
	es := make(errorSliceIndex)
	for i, v := range s {
		if err := v.Validate(); err != nil {
			switch e := err.(type) {
			case errorSliceIndex, errorMap:
				es[i] = e
			case error:
				es[i] = e.Error()
			default:
				es[i] = e
			}
		}
	}

	return es
}

type errorMulti struct {
	err1 error
	err2 error
}

func (e *errorMulti) Error() string {
	return errors.Join(e.err1, e.err2).Error()
}

func (e *errorMulti) Errors() []error {
	var errs []error
	if e.err1 != nil {
		errs = append(errs, e.err1)
	}
	if e.err2 != nil {
		errs = append(errs, e.err2)
	}

	return errs
}

func (e *errorMulti) Unwrap() error {
	return cmp.Or(e.err1, e.err2)
}
