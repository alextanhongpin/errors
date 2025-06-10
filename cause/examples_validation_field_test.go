package cause_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alextanhongpin/errors/cause"
)

type Email string

func (e Email) Validate() error {
	if !strings.Contains(string(e), "@") {
		return errors.New("invalid email format")
	}

	return nil
}

type Account struct {
	Email Email
}

func (a *Account) Validate() error {
	return cause.Map{
		"email": cause.Required(a.Email),
	}.Err()
}

func ExampleFields_email_valid() {
	a := &Account{
		Email: Email("john.doe@appleseed.com"),
	}
	validateAccount(a)

	// Output:
	// is map: false
	// is nil: true
	// err: <nil>
	// null
}

func ExampleFields_email_invalid_email() {
	a := &Account{
		Email: Email("invalid-email-format"),
	}
	validateAccount(a)

	// Output:
	// is map: true
	// is nil: false
	// err: invalid fields: email
	// {
	//   "email": "invalid email format"
	// }
}

func ExampleFields_email_empty() {
	a := &Account{}
	validateAccount(a)

	// Output:
	// is map: true
	// is nil: false
	// err: invalid fields: email
	// {
	//   "email": "required"
	// }
}

func validateAccount(a *Account) {
	var err error = a.Validate()
	var me errorMap
	fmt.Println("is map:", errors.As(err, &me))
	fmt.Println("is nil:", err == nil)
	fmt.Println("err:", err)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
