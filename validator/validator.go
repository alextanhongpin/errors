package validator

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

// ErrRequired is returned when a required value is missing or zero.
var ErrRequired = errors.New("required")

// validatable represents any type that can validate itself.
// Types implementing this interface can be used in validation chains.
type validatable interface {
	Validate() error
}

// ValidateManyFunc applies a validation function to each item in a slice and collects errors.
// Returns nil if no errors are found, otherwise returns an errorSlice.
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

// ValidateMany validates each item in a slice of validatable types and collects errors.
// Returns nil if no errors are found, otherwise returns an errorSlice.
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

// errorMap is a map of field names to errors, used for structured error reporting.
type errorMap map[string]any

// Map returns the underlying map of the errorMap.
func (e errorMap) Map() map[string]any {
	return e
}

// Error returns a string representation of the invalid fields in the errorMap.
func (e errorMap) Error() string {
	return fmt.Sprintf("invalid fields: %s", strings.Join(slices.Sorted(maps.Keys(e)), ", "))
}

// unwrapMany is an interface for errors that can be unwrapped into multiple errors.
type unwrapMany interface {
	Unwrap() []error
}

// errorField associates a field name with an error.
type errorField struct {
	field string
	err   error
}

// Map converts a map of field errors into a structured errorMap, handling nested errors.
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

// Required returns ErrRequired if the value is zero, otherwise joins additional errors.
func Required(value any, errs ...error) error {
	if isZero(value) {
		return ErrRequired
	}

	return errors.Join(errs...)
}

// Optional returns nil if the value is zero, otherwise joins additional errors.
func Optional(value any, errs ...error) error {
	if isZero(value) {
		return nil
	}

	return errors.Join(errs...)
}

// Validate calls Validate on a validatable type if it is not zero.
func Validate[T validatable](v T) error {
	if isZero(v) {
		return nil
	}

	return v.Validate()
}

// When returns an error with the given message if valid is true, otherwise returns nil.
func When(valid bool, msg string, args ...any) error {
	if valid {
		return fmt.Errorf(msg, args...)
	}

	return nil
}

// Assert returns an error with the given message if valid is false, otherwise returns nil.
func Assert(valid bool, msg string, args ...any) error {
	if valid {
		return nil
	}

	return fmt.Errorf(msg, args...)
}

// isZero checks if a value is the zero value for its type.
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

// errorSlice is a map of index to error, used for reporting errors in slices.
type errorSlice map[int]error

// Error returns a string indicating the slice is invalid.
func (e errorSlice) Error() string {
	return "invalid slice"
}

// WhenMap returns an error listing all messages whose condition is true.
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

// AssertMap returns an error listing all messages whose condition is false.
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

// Min returns an error if value < min.
func Min[T ~int | ~float64](value T, min T, msg string) error {
	if value < min {
		return errors.New(msg)
	}
	return nil
}

// Max returns an error if value > max.
func Max[T ~int | ~float64](value T, max T, msg string) error {
	if value > max {
		return errors.New(msg)
	}
	return nil
}

// Between returns an error if value is not in [min, max].
func Between[T ~int | ~float64](value T, min T, max T, msg string) error {
	if value < min || value > max {
		return errors.New(msg)
	}
	return nil
}

// Length returns an error if the length of value != n.
func Length(value any, n int, msg string) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String || v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() != n {
			return errors.New(msg)
		}
		return nil
	}
	return errors.New("unsupported type for Length")
}

// MinLength returns an error if the length of value < min.
func MinLength(value any, min int, msg string) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String || v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() < min {
			return errors.New(msg)
		}
		return nil
	}
	return errors.New("unsupported type for MinLength")
}

// MaxLength returns an error if the length of value > max.
func MaxLength(value any, max int, msg string) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String || v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() > max {
			return errors.New(msg)
		}
		return nil
	}
	return errors.New("unsupported type for MaxLength")
}

// Any returns an error with msg if valid is false.
func Any(valid bool, msg string) error {
	if !valid {
		return errors.New(msg)
	}
	return nil
}
