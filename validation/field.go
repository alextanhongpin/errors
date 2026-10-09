package validation

import (
	"fmt"
	"maps"
	"slices"
)

// validatable is implemented by types that can provide validation errors.
type validatable interface {
	Errors() Errors
}

// Slice is a wrapper for validating a slice of validatable items.
// It is created via SliceOf.
type Slice[T validatable] []T

// SliceOf converts a slice to a Slice[T] for validation.
func SliceOf[T validatable](s []T) Slice[T] {
	return Slice[T](s)
}

// Errors returns a map of field errors for each element in the slice.
// Keys are formatted as "[i].field". Nil elements report "[i]: required";
// an empty field key reports an error on the element itself.
func (ss Slice[T]) Errors() Errors {
	field := make(map[string][]string)
	for i, s := range ss {
		prefix := fmt.Sprintf("[%d]", i)
		if isNil(s) {
			field[prefix] = []string{"required"}
			continue
		}
		for k, v := range s.Errors() {
			field[nestedKey(prefix, k)] = v
		}
	}
	return field
}

// Errors is a map from field path to a list of error messages.
// Use the helper methods If, Optional, Required, and Add to populate it.
type Errors map[string][]string

// If adds an error message for a key if the condition is true.
func (e Errors) If(key string, value bool, val string) {
	if !value {
		return
	}
	e[key] = append(e[key], val)
}

// Optional validates a nested validatable value only if it is non-nil/non-zero.
// If the value is nil or zero, validation is skipped.
func (e Errors) Optional[T validatable](key string, val T) {
	if IsNilOrZero(val) {
		return
	}

	e.nest(key, val)
}

// Required validates a nested validatable value and adds a "required" error if it is nil or zero.
// Otherwise, nested errors are merged under the given key.
func (e Errors) Required[T validatable](key string, val T) {
	if IsNilOrZero(val) {
		e[key] = append(e[key], "required")
		return
	}

	e.nest(key, val)
}

// nest merges errors from a validatable value under the given key prefix.
func (e Errors) nest[T validatable](key string, val T) {
	for k, v := range val.Errors() {
		k = nestedKey(key, k)
		e[k] = append(e[k], v...)
	}
}

// An empty child key denotes an error on the nested value itself.
func nestedKey(parent, child string) string {
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}
	if child[0] == '[' {
		return parent + child
	}
	return parent + "." + child
}

// Add appends an error message for the given key.
func (e Errors) Add(key string, vals string) {
	e[key] = append(e[key], vals)
}

// Error converts the Errors map into an error.
// Returns nil if there are no errors.
// The resulting error is an ErrorMap with sorted keys.
func (e Errors) Error() error {
	keys := slices.Sorted(maps.Keys(e))
	errs := make(ErrorMap)
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
