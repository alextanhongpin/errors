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
	return cause.Map{}.
		Required("name", a.Name).
		Required("books", cause.Slice(a.Books)).
		// Required("books", cause.Slice(a.Books, (*Book).Validate)).
		Optional("likes", cause.SliceFunc(a.Likes, validateLike)).
		AsError()
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
	return cause.Map{}.
		Required("title", b.Title).
		Optional("year", b.Year, cause.When(b.Year < 2000, "too old")).
		Optional("languages", len(b.Languages), cause.When(len(b.Languages) > 1, "does not support multilingual")).
		AsError()
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
	validateauthor(a)

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
	validateauthor(a)

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
	validateauthor(a)

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
	validateauthor(a)

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

func validateauthor(a *Author) {
	err := a.Validate()
	fmt.Println("is nil:", err == nil)
	fmt.Println("err:", err)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
