package cause

import "github.com/alextanhongpin/errors/codes"

var (
	ErrAborted            = New(codes.Aborted, "ABORTED", "The operation was aborted")
	ErrBadRequest         = New(codes.BadRequest, "BAD_REQUEST", "The request is invalid")
	ErrCanceled           = New(codes.Canceled, "CANCELED", "The operation was canceled")
	ErrConflict           = New(codes.Conflict, "CONFLICT", "The request could not be completed due to a conflict with the current state of the target resource")
	ErrDataLoss           = New(codes.DataLoss, "DATA_LOSS", "Unrecoverable data loss or corruption")
	ErrDeadlineExceeded   = New(codes.DeadlineExceeded, "DEADLINE_EXCEEDED", "The deadline expired before the operation could complete")
	ErrExists             = New(codes.Exists, "EXISTS", "The resource that a client tried to create already exists")
	ErrForbidden          = New(codes.Forbidden, "FORBIDDEN", "The caller does not have permission to execute the specified operation")
	ErrInternal           = New(codes.Internal, "INTERNAL", "Internal server error")
	ErrNotFound           = New(codes.NotFound, "NOT_FOUND", "The specified resource was not found")
	ErrNotImplemented     = New(codes.NotImplemented, "NOT_IMPLEMENTED", "The operation is not implemented or not supported")
	ErrOutOfRange         = New(codes.OutOfRange, "OUT_OF_RANGE", "The operation was attempted past the valid range")
	ErrPreconditionFailed = New(codes.PreconditionFailed, "PRECONDITION_FAILED", "The operation was rejected because the system is not in a state required for the operation's execution")
	ErrTooManyRequests    = New(codes.TooManyRequests, "TOO_MANY_REQUESTS", "The caller has sent too many requests in a given amount of time")
	ErrUnauthorized       = New(codes.Unauthorized, "UNAUTHORIZED", "The request does not have valid authentication credentials for the operation")
	ErrUnavailable        = New(codes.Unavailable, "UNAVAILABLE", "The service is currently unavailable")
	ErrUnknown            = New(codes.Unknown, "UNKNOWN", "An unknown error occurred")
)
