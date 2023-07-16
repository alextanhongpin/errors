# Errors


[![](https://godoc.org/github.com/alextanhongpin/errors?status.svg)](http://godoc.org/github.com/alextanhongpin/errors)


The `errors` package solves some of the usecases that I discovered when working with errors in golang.

That includes
- conflict in package name (it's annoying)
- no stacktrace
- no standardized way of grouping errors, especially when mapping the error back to HTTP/gRPC status code
- custom error design

Some things that are considered, but is not included when designing this package is
- localization
