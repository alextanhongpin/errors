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

// validatable represents any type that can validate itself.
// Types implementing this interface can be used in validation chains.
type validatable interface {
	Validate() error
}

// Map is a validation map that associates field names with validation results.
// It provides a convenient way to collect validation errors for multiple fields
// and convert them into a structured error format.
//
// Example:
//
//	err := Map{
//	    "name": Required(user.Name),
//	    "age":  Optional(user.Age).When(user.Age < 0, "must be positive"),
//	}.Err()
type Map map[string]any

// Err processes the validation map and returns a structured error if any
// validations failed. It handles nested validations, slice validations,
// and converts the results into an errorMap for structured error reporting.
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

// Optional creates a validation builder for optional fields.
// If the value is zero (nil, empty string, 0, empty slice, etc.),
// no validation is performed. Otherwise, the value is validated
// if it implements validatable or if it's a slice of validatable items.
//
// Returns nil if the field is zero-valued, allowing the validation
// chain to be safely ignored.
func Optional(val any) *Builder {
	if isZero(val) {
		return nil
	}

	if v, ok := isSlice(val); ok {
		return &Builder{
			v: v,
		}
	}

	if v, ok := val.(validatable); ok {
		return &Builder{
			v: v,
		}
	}

	return &Builder{}
}

// Required creates a validation builder for required fields.
// If the value is zero (nil, empty string, 0, empty slice, etc.),
// a "required" error message is added. Otherwise, behaves like Optional.
func Required(val any) *Builder {
	return RequiredMessage(val, "required")
}

// RequiredMessage creates a validation builder for required fields with a custom message.
// If the value is zero, the specified message is used instead of the default "required".
func RequiredMessage(val any, msg string) *Builder {
	if isZero(val) {
		return &Builder{
			msgs: []string{msg},
		}
	}

	if v, ok := isSlice(val); ok {
		return &Builder{
			v: v,
		}
	}

	if v, ok := val.(validatable); ok {
		return &Builder{
			v: v,
		}
	}

	return &Builder{}
}

// SliceFunc creates a slice validator by applying a validation function
// to each element in the slice. Each element's validation result is wrapped
// in a validator that can be processed by the validation system.
//
// Example:
//
//	SliceFunc(users, func(u User) error {
//	    return u.Validate()
//	})
func SliceFunc[T any](vs []T, fn func(T) error) sliceValidator {
	res := make([]validatable, len(vs))
	for i, v := range vs {
		res[i] = &validator{
			err: fn(v),
		}
	}

	return res
}

// Builder provides a fluent interface for building validation chains.
// It accumulates validation messages and can validate nested structures.
type Builder struct {
	msgs []string    // Accumulated validation messages
	v    validatable // Nested validatable object
}

// When adds a conditional validation message to the builder.
// If the condition is true, the message is added to the validation errors.
// If the builder is nil (from Optional with zero value), this method
// is safe to call and returns nil.
//
// Example:
//
//	Optional(age).When(age < 0, "must be positive").When(age > 150, "unrealistic")
func (b *Builder) When(cond bool, msg string) *Builder {
	if b == nil {
		return nil
	}

	if cond {
		b.msgs = append(b.msgs, msg)
	}

	return b
}

// Validate executes the validation chain and returns any accumulated errors.
// It combines validation messages and nested validation results into a
// single error. Returns nil if no errors were found.
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

// errorMap represents a collection of field-level validation errors.
// It implements error interface and provides structured access to validation failures.
type errorMap map[string]any

// Map returns the underlying map of validation errors.
func (e errorMap) Map() map[string]any {
	return e
}

// Error returns a human-readable error message listing all failed fields.
func (e errorMap) Error() string {
	return "invalid fields: " + e.String()
}

// String returns a comma-separated list of invalid field names.
func (e errorMap) String() string {
	return strings.Join(slices.Sorted(maps.Keys(e)), ", ")
}

// isSlice checks if the given value is a slice of validatable items.
// Returns a sliceValidator and true if successful, nil and false otherwise.
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

// isZero checks if a value is considered "zero" for validation purposes.
// This includes nil, zero values, and empty slices.
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

// joinStrings joins non-empty strings with commas.
func joinStrings(s ...string) string {
	return strings.Join(filterZero(s), ", ")
}

// filterZero removes zero values from a slice.
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

// validator wraps an error to make it validatable.
type validator struct {
	err error
}

// Validate returns the wrapped error.
func (v *validator) Validate() error {
	return v.err
}

// errorIndex represents an error at a specific index in a slice.
type errorIndex struct {
	pos int
	err error
}

// Error returns a formatted error message including the index position.
func (ei errorIndex) Error() string {
	return fmt.Sprintf("error at index %d: %s", ei.pos, ei.err.Error())
}

// sliceValidator validates each element in a slice.
type sliceValidator []validatable

// Validate validates all elements in the slice and returns accumulated errors.
func (s sliceValidator) Validate() error {
	var em errorMulti
	for i, v := range s {
		if err := v.Validate(); err != nil {
			em = append(em, errorIndex{i, err})
		}
	}

	return em
}

// errorMulti represents multiple errors joined together.
type errorMulti []error

// Error returns a formatted error message for all accumulated errors.
func (e errorMulti) Error() string {
	return errors.Join(e).Error()
}
