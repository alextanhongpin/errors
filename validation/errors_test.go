package validation_test

import (
	"testing"

	"github.com/alextanhongpin/errors/validation"
	"github.com/alextanhongpin/evaltest"
)

func TestValidation(t *testing.T) {
	evaltest.Run(t, func(t *evaltest.T, input *User) (any, error) {
		return nil, input.Validate()
	})
}

func TestNested(t *testing.T) {
	evaltest.Run(t, func(t *evaltest.T, input *Book) (any, error) {
		return nil, input.Validate()
	})
}

func TestSlice(t *testing.T) {
	evaltest.Run(t, func(t *evaltest.T, input *BTOApplication) (any, error) {
		return nil, input.Validate()
	})
}

/* Fixtures */

type Book struct {
	Title  string
	Author *Author
}

func (b *Book) Validate() error {
	f := make(validation.Errors)
	f.If("title", b.Title == "", "required")
	f.Required("author", b.Author)
	return f.Error()
}

type Author struct {
	Name string
}

func (a *Author) Errors() validation.Errors {
	f := make(validation.Errors)
	f.If("name", a.Name == "", "required")
	return f
}

func (a *Author) Validate() error {
	return a.Errors().Error()
}

type User struct {
	Name string
	Age  int
}

func (u *User) Errors() validation.Errors {
	f := make(validation.Errors)
	f.If("name", u.Name == "", "required")
	f.If("name", len(u.Name) > 100, "too long")
	f.If("age", u.Age < 13, "underage")
	return f
}

func (u *User) Validate() error {
	return u.Errors().Error()
}

type BTOApplication struct {
	Spouses  []*User
	Children []*User
}

func (a *BTOApplication) Errors() validation.Errors {
	f := make(validation.Errors)
	f.Required("spouses", validation.SliceOf(a.Spouses))
	f.Optional("children", validation.SliceOf(a.Children))
	return f
}

func (a *BTOApplication) Validate() error {
	return a.Errors().Error()
}
