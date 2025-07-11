// Package cause provides validation utilities for structured error handling.
// This file contains validation helpers for building complex validation logic
// with support for nested structures, slices, and conditional validation.
package cause

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

var ErrRequired = errors.New("required")

// validatable represents any type that can validate itself.
// Types implementing this interface can be used in validation chains.
type validatable interface {
	Validate() error
}

func ValidateManyFunc[T any](items []T, fn func(T) error) error {
	sliceErr := make(errorSlice)

	for i, v := range items {
		if err := fn(v); err != nil {
			sliceErr[i] = err
		}
	}
	if len(sliceErr) == 0 {
		return nil
	}

	return sliceErr
}

func ValidateMany[T validatable](value []T) error {
	sliceErr := make(errorSlice)

	for i, v := range value {
		if err := v.Validate(); err != nil {
			sliceErr[i] = err
		}
	}

	if len(sliceErr) == 0 {
		return nil
	}

	return sliceErr
}

type errorMap map[string]any

func (e errorMap) Error() string {
	return fmt.Sprintf("invalid fields: %s", strings.Join(slices.Collect(maps.Keys(e)), ", "))
}

type unwrapMany interface {
	Unwrap() []error
}

type errorField struct {
	field string
	err   error
}

func Map(m map[string]error) error {
	var fieldErrors []errorField
	for field, err := range m {
		if err == nil {
			continue
		}

		fieldErrors = append(fieldErrors, errorField{field: field, err: err})
	}
	result := make(errorMap)

	for len(fieldErrors) > 0 {
		fieldErr := fieldErrors[0]
		fieldErrors = fieldErrors[1:]

		if many, ok := fieldErr.err.(unwrapMany); ok {
			for _, e := range many.Unwrap() {
				fieldErrors = append(fieldErrors, errorField{field: fieldErr.field, err: e})
			}
			continue
		}

		if indexErr, ok := fieldErr.err.(errorSlice); ok {
			for i, e := range indexErr {
				fieldErrors = append(fieldErrors, errorField{field: fmt.Sprintf("%s[%d]", fieldErr.field, i), err: e})
			}
			continue
		}

		if mapErr, ok := fieldErr.err.(errorMap); ok {
			result[fieldErr.field] = mapErr
			continue
		}

		result[fieldErr.field] = fieldErr.err.Error()
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func Required(value any, errs ...error) error {
	if isZero(value) {
		return ErrRequired
	}

	return errors.Join(errs...)
}

func Optional(value any, errs ...error) error {
	if isZero(value) {
		return nil
	}

	return errors.Join(errs...)
}

func Validate[T validatable](v T) error {
	if isZero(v) {
		return nil
	}

	return v.Validate()
}

func When(valid bool, msg string, args ...any) error {
	if valid {
		return fmt.Errorf(msg, args...)
	}

	return nil
}

func Assert(valid bool, msg string, args ...any) error {
	if valid {
		return nil
	}

	return fmt.Errorf(msg, args...)
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

type errorSlice map[int]error

func (e errorSlice) Error() string {
	return "invalid slice"
}

func WhenMap(conds map[string]bool) error {
	var msgs []string
	for msg, valid := range conds {
		if valid {
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	slices.Sort(msgs)
	return errors.New(strings.Join(msgs, ", "))
}

func AssertMap(conds map[string]bool) error {
	var msgs []string
	for msg, valid := range conds {
		if !valid {
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	slices.Sort(msgs)
	return errors.New(strings.Join(msgs, ", "))
}
