package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// validatableItem wraps myType and implements validatable
type myType struct{ valid bool }
type validatableItem struct{ myType }

func (v validatableItem) Validate() error {
	if !v.valid {
		return errors.New("not valid")
	}
	return nil
}

// Example of using Required and Optional
func Example_required_optional() {
	fmt.Println(Required(nil))
	fmt.Println(Required("foo"))
	fmt.Println(Optional(nil, errors.New("fail")))
	fmt.Println(Optional("foo", errors.New("fail")))
	// Output:
	// required
	// <nil>
	// <nil>
	// fail
}

// Example of ValidateManyFunc
func ExampleValidateManyFunc() {
	items := []int{1, 2, 3}
	err := ValidateManyFunc(items, func(i int) error {
		if i%2 == 0 {
			return fmt.Errorf("even: %d", i)
		}
		return nil
	})
	fmt.Println(err)
	// Output:
	// invalid slice
}

// Example of ValidateMany with validatable
func ExampleValidateMany() {
	items := []validatableItem{{myType{true}}, {myType{false}}, {myType{true}}}
	err := ValidateMany(items)
	fmt.Println(err)
	// Output:
	// invalid slice
}

// Example of Map for field errors
func ExampleMap() {
	m := map[string]error{"field": errors.New("fail")}
	err := Map(m)
	fmt.Println(err)
	// Output:
	// invalid fields: field
}

type Child struct {
	Name string
}

func (c Child) Validate() error {
	return Required(c.Name)
}

type Parent struct {
	Name     string
	Age      int
	Children []Child
}

func (p Parent) Validate() error {
	fieldErrs := map[string]error{
		"name":     Required(p.Name, Assert(len(p.Name) >= 3, "name must be at least 3 characters")),
		"age":      Required(p.Age, Assert(p.Age >= 18, "age must be at least 18")),
		"children": ValidateMany(p.Children),
	}
	return Map(fieldErrs)
}

// Example of nested struct validation and error aggregation
func Example_nested_validation() {
	p := Parent{
		Name:     "",
		Age:      0,
		Children: []Child{{Name: ""}, {Name: "Alice"}},
	}
	err := p.Validate()
	fmt.Println(err)
	// Output:
	// invalid fields: age, children[0], name
}

type Address struct {
	Street string
	City   string
	Zip    string
}

func (a Address) Validate() error {
	return Map(map[string]error{
		"street": Required(a.Street, Assert(len(a.Street) >= 5, "street must be at least 5 characters")),
		"city":   Required(a.City),
		"zip":    Required(a.Zip, Assert(len(a.Zip) == 5, "zip must be 5 characters")),
	})
}

type Product struct {
	Code  string
	Price float64
}

func (p Product) Validate() error {
	return Map(map[string]error{
		"code":  Required(p.Code, Assert(len(p.Code) >= 3, "code must be at least 3 characters")),
		"price": Assert(p.Price >= 1, "price must be at least 1"),
	})
}

type Order struct {
	ID      string
	Amount  float64
	Product Product
}

func (o Order) Validate() error {
	return Map(map[string]error{
		"id":      Required(o.ID, Assert(len(o.ID) >= 4, "order id must be at least 4 characters")),
		"amount":  Required(o.Amount, Assert(o.Amount >= 1, "amount must be at least 1")),
		"product": o.Product.Validate(),
	})
}

type User struct {
	Name     string
	Email    string
	Address  Address
	Orders   []Order
	Bio      string // Optional field
	Password string // New field for AssertMap demo
}

func (u User) Validate() error {
	return Map(map[string]error{
		"name":    Required(u.Name, Assert(len(u.Name) >= 3, "name must be at least 3 characters")),
		"email":   Required(u.Email, Assert(len(u.Email) >= 5, "email must be at least 5 characters")),
		"address": Required(u.Address, Validate(u.Address)),
		"orders":  Required(u.Orders, ValidateMany(u.Orders)),
		"bio":     Optional(u.Bio, Assert(len(u.Bio) <= 100, "bio must be at most 100 characters")),
		"password": Required(u.Password, AssertMap(map[string]bool{
			"must be at least 8 characters": len(u.Password) >= 8,
			"must contain a number":         containsNumber(u.Password),
			"must contain a capital letter": containsCapital(u.Password),
		})),
	})
}

// Example of validating a real-world API payload with nested structs and slices
func Example_api_payload_validation() {
	payload := User{
		Name:     "",
		Email:    "",
		Address:  Address{Street: "", City: "", Zip: "12345"},
		Orders:   []Order{{ID: "", Amount: 0}, {ID: "A123", Amount: 100}},
		Password: "abc",
	}
	err := payload.Validate()
	if err != nil {
		b, _ := json.MarshalIndent(err, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println("<nil>")
	}
	// Output:
	// {
	//   "address.city": "required",
	//   "address.street": "required",
	//   "email": "required",
	//   "name": "required",
	//   "orders[0].amount": "required",
	//   "orders[0].id": "required",
	//   "orders[0].product.code": "required",
	//   "orders[0].product.price": "price must be at least 1",
	//   "orders[1].product.code": "required",
	//   "orders[1].product.price": "price must be at least 1",
	//   "password": "must be at least 8 characters, must contain a capital letter, must contain a number"
	// }
}

func containsNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func containsCapital(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func printValidation(label string, u User) {
	err := u.Validate()
	fmt.Printf("--- %s ---\n", label)
	if err != nil {
		b, _ := json.MarshalIndent(err, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println("<nil>")
	}
}

func TestUserValidationScenarios(t *testing.T) {
	scenarios := []struct {
		label string
		user  User
	}{
		{
			label: "missing all required fields",
			user:  User{},
		},
		{
			label: "missing address",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Orders:   []Order{{ID: "A123", Amount: 100}},
				Password: "Abcdef12",
			},
		},
		{
			label: "missing orders",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{},
				Password: "Abcdef12",
			},
		},
		{
			label: "missing password",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{{ID: "A123", Amount: 100}},
				Password: "",
			},
		},
		{
			label: "missing address fields",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "", City: "", Zip: ""},
				Orders:   []Order{{ID: "A123", Amount: 100}},
				Password: "Abcdef12",
			},
		},
		{
			label: "missing order fields",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{{ID: "", Amount: 0}},
				Password: "Abcdef12",
			},
		},
		{
			label: "missing product fields",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{{ID: "A123", Amount: 100, Product: Product{Code: "", Price: 0}}},
				Password: "Abcdef12",
			},
		},
		{
			label: "invalid product",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{{ID: "A123", Amount: 100, Product: Product{Code: "AB", Price: 0.5}}},
				Password: "Abcdef12",
			},
		},
		{
			label: "valid user",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Address:  Address{Street: "123 Main St", City: "NYC", Zip: "12345"},
				Orders:   []Order{{ID: "A123", Amount: 100, Product: Product{Code: "ABC123", Price: 10}}},
				Password: "Abcdef12",
			},
		},
	}
	for _, s := range scenarios {
		printValidation(s.label, s.user)
	}
}
