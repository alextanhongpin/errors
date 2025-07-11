package validator

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// mockValidatable is a helper for testing validatable interface
type mockValidatable struct {
	err error
}

func (m mockValidatable) Validate() error {
	return m.err
}

func TestValidateManyFunc(t *testing.T) {
	items := []int{1, 2, 3}
	err := ValidateManyFunc(items, func(i int) error {
		if i%2 == 0 {
			return fmt.Errorf("even: %d", i)
		}
		return nil
	})
	if err == nil {
		t.Errorf("expected error for even numbers")
	}
}

func TestValidateManyFunc_EmptySlice(t *testing.T) {
	items := []int{}
	err := ValidateManyFunc(items, func(i int) error { return nil })
	if err != nil {
		t.Errorf("expected nil for empty slice")
	}
}

func TestValidateMany(t *testing.T) {
	items := []mockValidatable{{nil}, {errors.New("fail")}, {nil}}
	err := ValidateMany(items)
	if err == nil {
		t.Errorf("expected error for failed validation")
	}
}

func TestValidateMany_EmptySlice(t *testing.T) {
	items := []mockValidatable{}
	err := ValidateMany(items)
	if err != nil {
		t.Errorf("expected nil for empty slice")
	}
}

func TestRequired(t *testing.T) {
	if Required(nil) != ErrRequired {
		t.Errorf("expected ErrRequired for nil")
	}
	if Required(0) != ErrRequired {
		t.Errorf("expected ErrRequired for zero int")
	}
	if Required("nonzero") != nil {
		t.Errorf("expected nil for nonzero string")
	}
}

func TestOptional(t *testing.T) {
	if Optional(nil, errors.New("fail")) != nil {
		t.Errorf("expected nil for nil value")
	}
	if Optional(0, errors.New("fail")) != nil {
		t.Errorf("expected nil for zero int")
	}
	if Optional("nonzero", errors.New("fail")) == nil {
		t.Errorf("expected error for nonzero string")
	}
}

func TestValidate(t *testing.T) {
	var v mockValidatable
	if Validate(v) != nil {
		t.Errorf("expected nil for zero validatable")
	}
	v = mockValidatable{errors.New("fail")}
	if Validate(v) == nil {
		t.Errorf("expected error for failed validation")
	}
}

func TestWhen(t *testing.T) {
	err := When(true, "should fail")
	if err == nil {
		t.Errorf("expected error when valid is true")
	}
	if When(false, "should not fail") != nil {
		t.Errorf("expected nil when valid is false")
	}
}

func TestAssert(t *testing.T) {
	err := Assert(false, "should fail")
	if err == nil {
		t.Errorf("expected error when valid is false")
	}
	if Assert(true, "should not fail") != nil {
		t.Errorf("expected nil when valid is true")
	}
}

func TestIsZero(t *testing.T) {
	if !isZero(nil) {
		t.Errorf("expected true for nil")
	}
	if !isZero(0) {
		t.Errorf("expected true for zero int")
	}
	if isZero(1) {
		t.Errorf("expected false for nonzero int")
	}
	if !isZero("") {
		t.Errorf("expected true for empty string")
	}
	if isZero("foo") {
		t.Errorf("expected false for nonzero string")
	}
	if !isZero([]int{}) {
		t.Errorf("expected true for empty slice")
	}
	if isZero([]int{1}) {
		t.Errorf("expected false for non-empty slice")
	}
	type S struct{ X int }
	var s *S
	if !isZero(s) {
		t.Errorf("expected true for nil struct pointer")
	}
	s = &S{}
	if !isZero(*s) {
		t.Errorf("expected true for zero struct value")
	}
	s.X = 1
	if isZero(*s) {
		t.Errorf("expected false for nonzero struct value")
	}
}

func TestErrorSlice_Error(t *testing.T) {
	err := errorSlice{0: errors.New("fail")}
	if err.Error() != "invalid slice" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestWhenMap(t *testing.T) {
	conds := map[string]bool{"foo": true, "bar": false}
	err := WhenMap(conds)
	if err == nil || err.Error() != "foo" {
		t.Errorf("expected error for true condition")
	}
	if WhenMap(map[string]bool{}) != nil {
		t.Errorf("expected nil for empty map")
	}
}

func TestAssertMap(t *testing.T) {
	conds := map[string]bool{"foo": true, "bar": false}
	err := AssertMap(conds)
	if err == nil || err.Error() != "bar" {
		t.Errorf("expected error for false condition")
	}
	if AssertMap(map[string]bool{}) != nil {
		t.Errorf("expected nil for empty map")
	}
}

func TestMap(t *testing.T) {
	m := map[string]error{"field": errors.New("fail")}
	err := Map(m)
	if err == nil {
		t.Errorf("expected error for field")
	}
	if err.Error() != "invalid fields: field" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestMap_NestedErrorMapAndSlice(t *testing.T) {
	m := map[string]error{
		"outer": errorMap{"inner": errorSlice{0: errors.New("fail")}},
	}
	err := Map(m)
	if err == nil {
		t.Errorf("expected error for nested errorMap and errorSlice")
	}
	if !strings.Contains(err.Error(), "outer") {
		t.Errorf("expected error to contain 'outer'")
	}
}

func TestRequired_JoinErrors(t *testing.T) {
	err := Required("nonzero", errors.New("err1"), errors.New("err2"))
	if err == nil || !strings.Contains(err.Error(), "err1") || !strings.Contains(err.Error(), "err2") {
		t.Errorf("expected joined errors in output")
	}
}

func TestOptional_JoinErrors(t *testing.T) {
	err := Optional("nonzero", errors.New("err1"), errors.New("err2"))
	if err == nil || !strings.Contains(err.Error(), "err1") || !strings.Contains(err.Error(), "err2") {
		t.Errorf("expected joined errors in output")
	}
}

func TestAssert_EmptyMessage(t *testing.T) {
	err := Assert(false, "")
	if err == nil || err.Error() != "" {
		t.Errorf("expected empty error message")
	}
}

func TestWhen_EmptyMessage(t *testing.T) {
	err := When(true, "")
	if err == nil || err.Error() != "" {
		t.Errorf("expected empty error message")
	}
}

func TestValidateManyFunc_AllNilErrors(t *testing.T) {
	items := []int{1, 2, 3}
	err := ValidateManyFunc(items, func(i int) error { return nil })
	if err != nil {
		t.Errorf("expected nil when all errors are nil")
	}
}

func TestMap_EmptyMap(t *testing.T) {
	m := map[string]error{}
	err := Map(m)
	if err != nil {
		t.Errorf("expected nil for empty error map")
	}
}
