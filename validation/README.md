# validation

Package validation provides utilities for collecting and formatting validation errors for Go structs. It supports nested validation, optional/required fields, and slice validation with path-based error keys.

The core type is `Errors`, a `map[string][]string` that collects errors per field path. Validation is performed by implementing the `validatable` interface:

```go
type validatable interface {
    Errors() Errors
}
```

## Installation

```bash
go get github.com/alextanhongpin/errors/validation
```

## Usage

### Simple validation

```go
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
```

```go
u := &User{Name: "", Age: 10}
err := u.Validate()
// err != nil
// err.Error() => "age: underage\nname: required"
```

### Nested validation

```go
type Author struct {
    Name string
}

func (a *Author) Errors() validation.Errors {
    f := make(validation.Errors)
    f.If("name", a.Name == "", "required")
    return f
}

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
```

If `Author` is nil, the error is `author: required`. If `Author` is present but `Name` is empty, the error key is `author.name`.

### Optional vs Required

```go
f.Optional("address", user.Address)  // skip if nil/zero
f.Required("email", user.Email)      // adds "required" if nil/zero
```

### Slice validation

```go
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
```

Slice errors are prefixed with the index, e.g. `spouses[0].name`, `children[1].age`.

`SliceOf` converts a slice to a validatable slice:

```go
validation.SliceOf(slice)
```

## API

- `type Errors map[string][]string`
  - `If(key string, value bool, val string)` – add error if condition is true
  - `Optional(key string, val T)` – validate nested value if non-nil/non-zero
  - `Required(key string, val T)` – validate nested value, add "required" if nil/zero
  - `Add(key string, vals string)` – append error message
  - `Error() error` – returns nil if no errors, otherwise `ErrorMap`

- `type ErrorMap map[string][]string` implements `error` with sorted, human-readable output.

- `type Slice[T validatable] []T` with `SliceOf[T]` constructor for slice validation.

- `IsNilOrZero(x any) bool` – helper to check nil or zero values.

## Testing

Tests use `evaltest` with YAML fixtures in `testdata/`.

```bash
go test ./...
```

## License

MIT
