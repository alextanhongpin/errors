package cause_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alextanhongpin/errors/cause"
)

type Author struct {
	Name  string
	Books []*Book
	Likes []string
}

func (a *Author) Validate() error {
	return cause.Map{
		"name":  cause.Required(a.Name),
		"books": cause.Required(a.Books),
		"likes": cause.Optional(cause.SliceFunc(a.Likes, validateLike)),
	}.Err()
}

func validateLike(in string) error {
	in = strings.ToLower(in)
	switch {
	case strings.Contains(in, "drugs"):
		return errors.New("drugs not allowed")
	case strings.Contains(in, "weapons"):
		return errors.New("weapons not allowed")
	default:
		return nil
	}
}

type Book struct {
	Title     string
	Year      int
	Languages []string
}

func (b *Book) Validate() error {
	return cause.Map{
		"title":     cause.Required(b.Title),
		"year":      cause.Optional(b.Year).When(b.Year < 2000, "too old"),
		"languages": cause.Optional(len(b.Languages)).When(len(b.Languages) > 1, "does not support multilingual"),
	}.Err()
}

func ExampleFields_author_valid() {
	a := &Author{
		Name: "John Appleseed",
		Books: []*Book{
			{
				Title:     "The Great Book",
				Year:      2021,
				Languages: []string{"English"},
			},
		},
		Likes: []string{"music", "art"},
	}
	validateAuthor(a)

	// Output:
	// is nil: true
	// err: <nil>
	// null
}

func ExampleFields_author_valid_optional_likes() {
	a := &Author{
		Name: "John Appleseed",
		Books: []*Book{
			{
				Title:     "The Great Book",
				Year:      2021,
				Languages: []string{"English"},
			},
		},
	}
	validateAuthor(a)

	// Output:
	// is nil: true
	// err: <nil>
	// null
}

func ExampleFields_author_invalid() {
	a := &Author{
		Name:  "",
		Books: nil,
		Likes: []string{"drugs", "weapons"},
	}
	validateAuthor(a)

	// Output:
	// is nil: false
	// err: invalid fields: books, likes[0], likes[1], name
	// {
	//   "books": "required",
	//   "likes[0]": "drugs not allowed",
	//   "likes[1]": "weapons not allowed",
	//   "name": "required"
	// }
}

func ExampleFields_author_invalid_books() {
	a := &Author{
		Name: "John Appleseed",
		Books: []*Book{
			{},
			{Year: 1999},
			{Title: "A book", Year: 2000, Languages: []string{"English", "French"}},
		},
	}
	validateAuthor(a)

	// Output:
	// is nil: false
	// err: invalid fields: books[0], books[1], books[2]
	// {
	//   "books[0]": {
	//     "title": "required"
	//   },
	//   "books[1]": {
	//     "title": "required",
	//     "year": "too old"
	//   },
	//   "books[2]": {
	//     "languages": "does not support multilingual"
	//   }
	// }
}

func validateAuthor(a *Author) {
	err := a.Validate()
	fmt.Println("is nil:", err == nil)
	fmt.Println("err:", err)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
