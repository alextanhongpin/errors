package cause_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

type User struct {
	Age  int
	Name string
}

func (u *User) Validate() error {
	return cause.Map{
		"age":  cause.Optional(u.Age).When(u.Age < 13, "under age limit"),
		"name": cause.Required(u.Name),
	}.Err()
}

func ExampleFields_user_valid() {
	u := &User{
		Name: "John Appleseed",
		Age:  0,
	}
	validateUser(u)

	// Output:
	// is map: false
	// is nil: true
	// err: <nil>
	// null
}

func ExampleFields_user_invalid_name() {
	u := &User{
		Name: "",
		Age:  0,
	}
	validateUser(u)

	// Output:
	// is map: true
	// is nil: false
	// err: invalid fields: name
	// {
	//   "name": "required"
	// }
}

func ExampleFields_user_invalid_age() {
	u := &User{
		Name: "John Appleseed",
		Age:  12,
	}
	validateUser(u)

	// Output:
	// is map: true
	// is nil: false
	// err: invalid fields: age
	// {
	//   "age": "under age limit"
	// }
}

func ExampleFields_user_invalid_age_and_name() {
	u := &User{
		Name: "",
		Age:  12,
	}
	validateUser(u)

	// Output:
	// is map: true
	// is nil: false
	// err: invalid fields: age, name
	// {
	//   "age": "under age limit",
	//   "name": "required"
	// }
}

type errorMap interface {
	Map() map[string]any
}

func validateUser(u *User) {
	var err error = u.Validate()
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
