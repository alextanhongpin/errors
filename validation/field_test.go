package validation_test

import (
	"reflect"
	"testing"

	"github.com/alextanhongpin/errors/validation"
)

type fieldErrors struct{ fields validation.Errors }

func (f fieldErrors) Errors() validation.Errors { return f.fields }

func TestNestedValueErrors(t *testing.T) {
	nested := fieldErrors{validation.Errors{"": {"invalid"}, "name": {"required"}, "[0]": {"invalid item"}}}
	for _, prefix := range []string{"profile", ""} {
		t.Run(prefix, func(t *testing.T) {
			got := make(validation.Errors)
			got.Required(prefix, nested)
			want := validation.Errors{"profile": {"invalid"}, "profile.name": {"required"}, "profile[0]": {"invalid item"}}
			if prefix == "" {
				want = nested.fields
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestSliceNilElements(t *testing.T) {
	got := make(validation.Errors)
	got.Required("users", validation.SliceOf([]*User{nil, {Name: "valid", Age: 18}, {Age: 18}}))
	want := validation.Errors{"users[0]": {"required"}, "users[2].name": {"required"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSliceValueErrors(t *testing.T) {
	got := validation.SliceOf([]fieldErrors{{validation.Errors{"": {"invalid"}, "[0]": {"nested"}}}, {}}).Errors()
	want := validation.Errors{"[0]": {"invalid"}, "[0][0]": {"nested"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
