# Go Validator

A flexible validation library for Go, supporting composable validation, error aggregation, and structured error reporting.

## Features
- Validate slices and nested structs with custom logic
- Required/Optional value checks
- Structured error reporting for fields and slices
- Composable error handling
- Table-driven scenario testing

## Installation
Add to your Go module:
```sh
go get github.com/alextanhongpin/errors/validator
```

## Usage

### Required and Optional
```go
err := validator.Required(nil) // returns ErrRequired
err := validator.Required("foo") // returns nil
err := validator.Optional(nil, errors.New("fail")) // returns nil
err := validator.Optional("foo", errors.New("fail")) // returns "fail"
```

### Handling Slices and Nested Validators

#### Slices
To validate each item in a slice, use `ValidateMany`:
```go
type Item struct {
  Value string
}
func (i Item) Validate() error {
  return validator.Required(i.Value, validator.Assert(len(i.Value) >= 3, "value must be at least 3 characters"))
}
items := []Item{{"foo"}, {""}, {"bar"}}
err := validator.ValidateMany(items)
if err != nil {
  fmt.Println(err) // errorSlice with index keys for failed items
}
```
You can also require the slice itself to be non-empty:
```go
err := validator.Required(items, validator.ValidateMany(items))
```

#### Nested Validators
For nested structs, call their `Validate` method inside your parent struct's validation:
```go
type Child struct {
  Name string
}
func (c Child) Validate() error {
  return validator.Required(c.Name)
}
type Parent struct {
  Name   string
  Childs []Child
}
func (p Parent) Validate() error {
  return validator.Map(map[string]error{
    "name": validator.Required(p.Name),
    "childs": validator.Required(p.Childs, validator.ValidateMany(p.Childs)),
  })
}
```
This pattern works for any level of nesting, and errors will be aggregated and mapped to their respective fields and indices.

### Validating Nested Structs
To validate nested structs, simply call their `Validate` method inside the parent struct's validation logic. Errors will be mapped to the corresponding field name and can be nested as deeply as needed.

```go
type Address struct {
  Street string
  City   string
}
func (a Address) Validate() error {
  return validator.Map(map[string]error{
    "street": validator.Required(a.Street),
    "city":   validator.Required(a.City),
  })
}

type User struct {
  Name    string
  Address Address
}
func (u User) Validate() error {
  return validator.Map(map[string]error{
    "name":    validator.Required(u.Name),
    "address": u.Address.Validate(), // Nested validation
  })
}

user := User{Name: "", Address: Address{Street: "", City: "NYC"}}
err := user.Validate()
if err != nil {
  b, _ := json.MarshalIndent(err, "", "  ")
  fmt.Println(string(b))
  // Output:
  // {
  //   "address": {
  //     "street": "required"
  //   },
  //   "name": "required"
  // }
}
```

This approach works for any level of nesting, and errors will be aggregated and mapped to their respective fields.

### Real-world API validation example
```go
// Define your types
 type Product struct {
   Code  string
   Price float64
 }
 func (p Product) Validate() error {
   return validator.Map(map[string]error{
     "code":  validator.Required(p.Code, validator.Assert(len(p.Code) >= 3, "code must be at least 3 characters")),
     "price": validator.Assert(p.Price >= 1, "price must be at least 1"),
   })
 }
 type Order struct {
   ID     string
   Amount float64
   Product Product
 }
 func (o Order) Validate() error {
   return validator.Map(map[string]error{
     "id":     validator.Required(o.ID, validator.Assert(len(o.ID) >= 4, "order id must be at least 4 characters")),
     "amount": validator.Required(o.Amount, validator.Assert(o.Amount >= 1, "amount must be at least 1")),
     "product": o.Product.Validate(),
   })
 }
 type Address struct {
   Street string
   City   string
   Zip    string
 }
 func (a Address) Validate() error {
   return validator.Map(map[string]error{
     "street": validator.Required(a.Street, validator.Assert(len(a.Street) >= 5, "street must be at least 5 characters")),
     "city":   validator.Required(a.City),
     "zip":    validator.Required(a.Zip, validator.Assert(len(a.Zip) == 5, "zip must be 5 characters")),
   })
 }
 type User struct {
   Name     string
   Email    string
   Address  Address
   Orders   []Order
   Password string
 }
 func (u User) Validate() error {
   passwordConds := map[string]bool{
     "must be at least 8 characters": len(u.Password) >= 8,
     "must contain a number":         containsNumber(u.Password),
     "must contain a capital letter": containsCapital(u.Password),
   }
   return validator.Map(map[string]error{
     "name":     validator.Required(u.Name, validator.Assert(len(u.Name) >= 3, "name must be at least 3 characters")),
     "email":    validator.Required(u.Email, validator.Assert(len(u.Email) >= 5, "email must be at least 5 characters")),
     "address":  u.Address.Validate(),
     "orders":   validator.Required(u.Orders, validator.ValidateMany(u.Orders)),
     "password": validator.Required(u.Password, validator.AssertMap(passwordConds)),
   })
 }

// Validate and print errors as JSON
user := User{}
err := user.Validate()
if err != nil {
  b, _ := json.MarshalIndent(err, "", "  ")
  fmt.Println(string(b))
}
```

### Table-driven scenario testing
```go
scenarios := []struct {
  label string
  user  User
}{
  {label: "missing all required fields", user: User{}},
  {label: "missing address", user: User{Name: "John", Email: "john@example.com", Orders: []Order{{ID: "A123", Amount: 100}}, Password: "Abcdef12"}},
  // ...more scenarios...
}
for _, s := range scenarios {
  fmt.Printf("--- %s ---\n", s.label)
  err := s.user.Validate()
  b, _ := json.MarshalIndent(err, "", "  ")
  fmt.Println(string(b))
}
```

### Composing Optional, Required, Assert, and AssertMap

You can nest and compose these helpers for flexible validation logic:

```go
// Required with nested Assert
err := validator.Required(username, validator.Assert(len(username) >= 3, "username must be at least 3 characters"))

// Optional with nested Assert
err := validator.Optional(bio, validator.Assert(len(bio) <= 100, "bio must be at most 100 characters"))

// Required with AssertMap for multiple rules
conds := map[string]bool{
  "must be at least 8 characters": len(password) >= 8,
  "must contain a number":         containsNumber(password),
  "must contain a capital letter": containsCapital(password),
}
err := validator.Required(password, validator.AssertMap(conds))
```

You can nest these as deeply as needed, e.g.:
```go
err := validator.Required(field, validator.Assert(len(field) > 0, "cannot be empty"), validator.AssertMap(conds))
```

### Notes on When and WhenMap

- `When` returns an error if the condition is true (useful for positive assertions):
  ```go
  err := validator.When(isDuplicate, "duplicate value detected")
  ```
- `Assert` returns an error if the condition is false (useful for negative assertions):
  ```go
  err := validator.Assert(isValid, "value is not valid")
  ```
- `WhenMap` aggregates all messages whose condition is true:
  ```go
  conds := map[string]bool{"field is empty": isEmpty, "field is duplicate": isDuplicate}
  err := validator.WhenMap(conds)
  ```
- `AssertMap` aggregates all messages whose condition is false:
  ```go
  conds := map[string]bool{"must be positive": value > 0, "must be even": value%2 == 0}
  err := validator.AssertMap(conds)
  ```

These helpers make it easy to build readable, composable validation logic for any scenario.

## Examples
See `validator_example_test.go` for more usage and scenario examples.

## Testing
Run unit and scenario tests:
```sh
go test -v ./...
```

## License
MIT
