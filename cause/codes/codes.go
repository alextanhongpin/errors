// Package codes provides standard error codes for structured error handling.
// This is a minimal implementation to satisfy the dependency.
package codes

// Code represents an error classification code.
type Code int

const (
	// OK indicates success (no error).
	OK Code = iota

	// Invalid indicates invalid input or arguments.
	Invalid

	// NotFound indicates a resource was not found.
	NotFound

	// AlreadyExists indicates a resource already exists.
	AlreadyExists

	// PermissionDenied indicates insufficient permissions.
	PermissionDenied

	// Unauthenticated indicates authentication is required.
	Unauthenticated

	// Unavailable indicates the service is unavailable.
	Unavailable

	// Internal indicates an internal error.
	Internal

	// DeadlineExceeded indicates the operation timed out.
	DeadlineExceeded

	// Aborted indicates the operation was aborted.
	Aborted

	// BadRequest indicates a bad request.
	BadRequest

	// Conflict indicates a conflict with the current state.
	Conflict

	// Exists is an alias for AlreadyExists.
	Exists = AlreadyExists
)

// String returns the string representation of the error code.
func (c Code) String() string {
	switch c {
	case OK:
		return "ok"
	case Invalid:
		return "invalid"
	case NotFound:
		return "not_found"
	case AlreadyExists: // This handles both AlreadyExists and Exists since they're the same value
		return "exists"
	case PermissionDenied:
		return "permission_denied"
	case Unauthenticated:
		return "unauthenticated"
	case Unavailable:
		return "unavailable"
	case Internal:
		return "internal"
	case DeadlineExceeded:
		return "deadline_exceeded"
	case Aborted:
		return "aborted"
	case BadRequest:
		return "bad_request"
	case Conflict:
		return "conflict"
	default:
		return "unknown"
	}
}
