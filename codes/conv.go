package codes

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var httpStatusByCode = map[Code]int{
	Aborted:            http.StatusConflict,
	BadRequest:         http.StatusBadRequest,
	Canceled:           499, // client closed request.
	Conflict:           http.StatusConflict,
	DataLoss:           http.StatusInternalServerError,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	Exists:             http.StatusConflict,
	Forbidden:          http.StatusForbidden,
	Internal:           http.StatusInternalServerError,
	NotFound:           http.StatusNotFound,
	NotImplemented:     http.StatusNotImplemented,
	OutOfRange:         http.StatusBadRequest,
	PreconditionFailed: http.StatusBadRequest,
	TooManyRequests:    http.StatusTooManyRequests,
	Unauthorized:       http.StatusUnauthorized,
	Unavailable:        http.StatusServiceUnavailable,
	Unknown:            http.StatusInternalServerError,
}

// ̱HTTP returns the HTTP status code for the given error code.
func HTTP(code Code) int {
	status, ok := httpStatusByCode[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}

// https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md
var grpcByCode = map[Code]codes.Code{
	Aborted:            codes.Aborted,
	BadRequest:         codes.InvalidArgument,
	Canceled:           codes.Canceled,
	Conflict:           codes.Aborted,
	DataLoss:           codes.DataLoss,
	DeadlineExceeded:   codes.DeadlineExceeded,
	Exists:             codes.AlreadyExists,
	Forbidden:          codes.PermissionDenied,
	Internal:           codes.Internal,
	NotFound:           codes.NotFound,
	NotImplemented:     codes.Unimplemented,
	OutOfRange:         codes.OutOfRange,
	PreconditionFailed: codes.FailedPrecondition,
	TooManyRequests:    codes.ResourceExhausted,
	Unauthorized:       codes.Unauthenticated,
	Unavailable:        codes.Unavailable,
	Unknown:            codes.Unknown,
}

// GRPC returns the gRPC code for the given error code.
func GRPC(code Code) codes.Code {
	c, ok := grpcByCode[code]
	if !ok {
		return codes.Internal
	}
	return c
}

var codeByGRPC = func() map[codes.Code]Code {
	m := make(map[codes.Code]Code)
	for k, v := range grpcByCode {
		m[v] = k
	}
	return m
}()

// GRPCToHTTP returns the HTTP code for the given grpc code.
func GRPCToHTTP(code codes.Code) int {
	c, ok := codeByGRPC[code]
	if !ok {
		return http.StatusInternalServerError
	}

	return HTTP(c)
}
