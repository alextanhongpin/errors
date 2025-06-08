package cause_test

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

type User struct {
	Name string
	Age  int
}

func (u *User) Validate() error {
	return cause.Fields{}.
		Required("name", u.Name).
		Optional("age", u.Age, cause.Cond(u.Age < 13, "under age limit")).
		AsError()
}

func ExampleFields_user_valid() {
	u := &User{
		Name: "John Appleseed",
		Age:  0,
	}
	validateUser(u)

	// Output:
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
	// is nil: false
	// err: invalid fields: age, name
	// {
	//   "age": "under age limit",
	//   "name": "required"
	// }
}

func validateUser(u *User) {
	err := u.Validate()
	fmt.Println("is nil:", err == nil)
	fmt.Println("err:", err)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
