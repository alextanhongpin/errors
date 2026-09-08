package validation

// Package validation provides utilities for collecting and formatting validation errors
// for Go structs. It supports nested validation, optional/required fields, and slice
// validation with path-based error keys.

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

var _ error = make(ErrorMap)

// filter removes empty strings from a slice of error messages.
func filter(ss []string) []string {
	var res []string
	for _, s := range ss {
		if s == "" {
			continue
		}
		res = append(res, s)
	}
	return res
}

// ErrorMap is a map from field path to a list of error messages.
// It implements the error interface and formats errors as "field: msg1, msg2".
type ErrorMap map[string][]string

// Error returns a human-readable string representation of the error map.
// Keys are sorted alphabetically for deterministic output.
func (e ErrorMap) Error() string {
	keys := slices.Sorted(maps.Keys(e))
	res := make([]string, len(keys))
	for i, key := range keys {
		res[i] = fmt.Sprintf("%s: %s", key, strings.Join(e[key], ", "))
	}

	return strings.Join(res, "\n")
}
