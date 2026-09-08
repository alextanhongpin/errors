package validation

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

var _ error = make(ErrorMap)

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

type ErrorMap map[string][]string

func (e ErrorMap) Error() string {
	keys := slices.Sorted(maps.Keys(e))
	res := make([]string, len(keys))
	for i, key := range keys {
		res[i] = fmt.Sprintf("%s: %s", key, strings.Join(e[key], ", "))
	}

	return strings.Join(res, "\n")
}
