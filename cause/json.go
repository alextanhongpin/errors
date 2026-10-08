package cause

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"

	"github.com/alextanhongpin/errors/codes"
)

const (
	maxErrorDepth = 64
	maxErrorNodes = 1024
	maxErrorBytes = 1 << 20
)

// From returns a structured error ready for JSON encoding. A *Error is returned
// unchanged. Other errors use the unknown code, an empty name, and err.Error()
// as their message, preserving their wrapped causes. Nil and typed-nil errors
// return nil. Ordinary errors with the same message intentionally compare equal
// after normalization; their original concrete types are not reconstructed.
func From(err error) *Error {
	if nilError(err) {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Code: codes.Unknown, Message: err.Error(), Cause: joinCauses(errorChildren(err))}
}

// MarshalJSON encodes the error representation, normalizing all causes.
// Stacks and logging attributes are omitted. Applications select safe messages,
// details, and causes before sending this same representation to clients.
func (e Error) MarshalJSON() ([]byte, error) {
	state := jsonTraversal{active: make(map[error]bool)}
	node, err := state.encode(&e, 1)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("cause: encode error: %w", err)
	}
	if len(data) > maxErrorBytes {
		return nil, errors.New("cause: error document byte limit exceeded")
	}
	return data, nil
}

// UnmarshalJSON validates an error document and replaces the receiver only after
// success. Detail numbers use json.Number. Decode into a *Error variable to
// represent JSON null as nil; null decoded into an Error value resets that value.
func (e *Error) UnmarshalJSON(data []byte) error {
	if e == nil {
		return errors.New("cause: cannot decode into nil error receiver")
	}
	if len(data) > maxErrorBytes {
		return errors.New("cause: error document byte limit exceeded")
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*e = Error{}
		return nil
	}
	state := jsonTraversal{}
	decoded, err := state.decode(data, 1)
	if err != nil {
		return err
	}
	*e = *decoded
	return nil
}

// jsonError is a wire-only node. It has no custom marshaler, so the entire tree
// is validated and counted once rather than resetting limits for each child.
type jsonError struct {
	Code    string         `json:"code"`
	Name    string         `json:"name"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	Causes  []*jsonError   `json:"causes,omitempty"`
}

type jsonTraversal struct {
	nodes  int
	active map[error]bool
}

func (t *jsonTraversal) visit(depth int) error {
	t.nodes++
	if depth > maxErrorDepth {
		return errors.New("cause: error depth limit exceeded")
	}
	if t.nodes > maxErrorNodes {
		return errors.New("cause: error node limit exceeded")
	}
	return nil
}

func (t *jsonTraversal) encode(err error, depth int) (*jsonError, error) {
	if nilError(err) {
		return nil, errors.New("cause: nil cause branch")
	}
	if err := t.visit(depth); err != nil {
		return nil, err
	}
	if reflect.ValueOf(err).Comparable() {
		if t.active[err] {
			return nil, errors.New("cause: cyclic error tree")
		}
		t.active[err] = true
		defer delete(t.active, err)
	}
	node := &jsonError{Code: codes.Unknown.String()}
	// Check the whole tree before calling Error, which can recursively render it.
	for _, child := range errorChildren(err) {
		encoded, err := t.encode(child, depth+1)
		if err != nil {
			return nil, err
		}
		node.Causes = append(node.Causes, encoded)
	}
	if e, ok := err.(*Error); ok {
		if !e.Code.Valid() {
			return nil, errors.New("cause: invalid error code")
		}
		node.Code, node.Name, node.Message = e.Code.String(), e.Name, e.Message
		node.Details = maps.Clone(e.Details)
	} else {
		node.Message = err.Error()
	}
	return node, nil
}

func (t *jsonTraversal) decode(data []byte, depth int) (*Error, error) {
	if err := t.visit(depth); err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("cause: decode error node: %w", err)
	}
	if fields == nil {
		return nil, errors.New("cause: expected error object")
	}
	var code string
	e := &Error{}
	for _, field := range []struct {
		name string
		out  *string
	}{{"code", &code}, {"name", &e.Name}, {"message", &e.Message}} {
		if err := readString(fields, field.name, field.out, true); err != nil {
			return nil, err
		}
	}
	var err error
	e.Code, err = parseCode(code)
	if err != nil {
		return nil, err
	}
	if raw, ok := fields["details"]; ok {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&e.Details); err != nil || e.Details == nil {
			return nil, errors.New("cause: details must be an object")
		}
	}
	if raw, ok := fields["causes"]; ok {
		var children []json.RawMessage
		if err := json.Unmarshal(raw, &children); err != nil || children == nil {
			return nil, errors.New("cause: causes must be an array")
		}
		var causes []error
		for _, child := range children {
			decoded, err := t.decode(child, depth+1)
			if err != nil {
				return nil, err
			}
			causes = append(causes, decoded)
		}
		e.Cause = joinCauses(causes)
	}
	return e, nil
}

func errorChildren(err error) []error {
	if e, ok := err.(*Error); ok {
		if e == nil || e.Cause == nil {
			return nil
		}
		// Flatten only our synthetic join; preserve original multi-cause wrapper nodes.
		if joined, ok := e.Cause.(*joinedCause); ok {
			return joined.children
		}
		return []error{e.Cause}
	}
	switch e := err.(type) {
	case interface{ Unwrap() []error }:
		return e.Unwrap()
	case interface{ Unwrap() error }:
		if child := e.Unwrap(); child != nil {
			return []error{child}
		}
	}
	return nil
}

// joinedCause uses errors.Join for Go traversal while keeping the wire branches
// stable across repeated round trips, without adding a synthetic JSON node.
type joinedCause struct {
	error
	children []error
}

func (e *joinedCause) Unwrap() []error { return e.children }
func joinCauses(children []error) error {
	switch len(children) {
	case 0:
		return nil
	case 1:
		return children[0]
	default:
		return &joinedCause{error: errors.Join(children...), children: children}
	}
}

func readString(fields map[string]json.RawMessage, name string, out *string, required bool) error {
	raw, ok := fields[name]
	if !ok && !required {
		return nil
	}
	if !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return fmt.Errorf("cause: missing or null %s", name)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("cause: %s must be a string", name)
	}
	return nil
}

func parseCode(value string) (codes.Code, error) {
	for code := codes.Aborted; code <= codes.Unknown; code++ {
		if code.String() == value {
			return code, nil
		}
	}
	return 0, fmt.Errorf("cause: unknown code %q", value)
}

func nilError(err error) bool {
	if err == nil {
		return true
	}
	v := reflect.ValueOf(err)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}
