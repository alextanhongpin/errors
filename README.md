# Errors


[![](https://godoc.org/github.com/alextanhongpin/errors?status.svg)](http://godoc.org/github.com/alextanhongpin/errors)


The `errors` package solves some of the pain points when working with errors in golang.

That includes
- optional stacktrace
- annotate stacktrace with cause
- grouping errors with errors `Code`
- custom errors with errors `Kind`
- mapping errors `Code` to HTTP/gRPC status code
- does not conflict with the standard errors package name

Some things that are considered, but is not included when designing this package is
- localization


## Installation

```bash
$ go get github.com/alextanhongpin/errors
```


## Usage


There are three subpackage in the `errors` package, each fulfilling different usecase:

- `causes`: create custom errors
- `codes`: standard error `codes` that can be mapped to `HTTP/gRPC` codes
- `stacktrace`: add stacktrace to errors and annotate cause

Each folder contains usage examples.
