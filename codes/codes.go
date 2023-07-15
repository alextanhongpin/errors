package codes

import "strings"

type Code int

//go:generate stringer -type=Code -linecomment
const (
	unknown Code = iota

	Aborted            // aborted
	BadRequest         // bad_request
	Canceled           // canceled
	Conflict           // conflict
	DataLoss           // data_loss
	DeadlineExceeded   // deadline_exceeded
	Exists             // exists
	Forbidden          // forbidden
	Internal           // internal
	NotFound           // not_found
	NotImplemented     // not_implemented
	OutOfRange         // out_of_range
	PreconditionFailed // precondition_failed
	TooManyRequests    // too_many_requests
	Unauthorized       // unauthorized
	Unavailable        // unavailable
	Unknown            // unknown
)

func (c Code) Valid() bool {
	return c > unknown && c <= Unknown
}

func Canonical(c Code) string {
	return strings.ToUpper(c.String())
}

var textByCode = map[Code]string{
	Aborted:            "Aborted",
	BadRequest:         "Bad Request",
	Canceled:           "Canceled",
	Conflict:           "Conflict",
	DataLoss:           "Data Loss",
	DeadlineExceeded:   "Deadline Exceeded",
	Exists:             "Exists",
	Forbidden:          "Forbidden",
	Internal:           "Internal",
	NotFound:           "Not Found",
	NotImplemented:     "Not Implemented",
	OutOfRange:         "Out of Range",
	PreconditionFailed: "Precondition Failed",
	TooManyRequests:    "Too Many Requests",
	Unauthorized:       "Unauthorized",
	Unavailable:        "Unavailable",
	Unknown:            "Unknown",
}

func Text(c Code) string {
	v, ok := textByCode[c]
	if ok {
		return v
	}
	return c.String()
}
