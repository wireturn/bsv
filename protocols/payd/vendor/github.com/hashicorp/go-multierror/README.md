# go-multierror

[![Build Status](http://img.shields.io/travis/hashicorp/go-multierror.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[travis]: https://travis-ci.org/hashicorp/go-multierror
[godocs]: https://godoc.org/github.com/hashicorp/go-multierror

`go-multierror` is a package for Go that provides a mechanism for
representing a list of `error` values as a single `error`.

This allows a function in Go to return an `error` that might actually
be a list of errors. If the caller knows this, they can unwrap the
list and access the errors. If the caller doesn't know, the error
formats to a nice human-readable format.

`go-multierror` is fully compatible with the Go standard library
[errors](https://golang.org/pkg/errors/) package, including the
functions `As`, `Is`, and `Unwrap`. This provides a standardized approach
for introspecting on error values.

## Installation and Docs

Install using `go get github.com/hashicorp/go-multierror`.

Full documentation is available at
http://godoc.org/github.com/hashicorp/go-multierror

## Usage

go-multierror is easy to use and purposely built to be unobtrusive in
existing Go applications/libraries that may not be aware of it.

**Building a list of errors**

The `Append` function is used to create a list of errors. This function
behaves a lot like the Go built-in `append` function: it doesn't matter
if the first argument is nil, a `multierror.Error`, or any other `error`,
the function behaves as you would expect.

```go
var result error

if err := step1(); err != nil {
	result = multierror.Append(result, err)
}
if err := step2(); err != nil {
	result = multierror.Append(result, err)
}

return result
```

**Customizing the formatting of the errors**

By specifying a custom `ErrorFormat`, you can customize the format
of the `Error() string` function:

```go
var result *multierror.Error

// ... accumulate errors here, maybe using Append

if result != nil {
	result.ErrorFormat = func([]error) string {
		return "errors!"
	}
}
```

**Accessing the list of errors**

`multierror.Error` implements `error` so if the caller doesn't know about
multierror, it will work just fine. But if you're aware a multierror might
be returned, you can use type switches to access the list of errors:

```go
if err := something(); err != nil {
	if merr, ok := err.(*multierror.Error); ok {
		// Use merr.Errors
	}
}
```

You can also use the standard [`errors.Unwrap`](https://golang.org/pkg/errors/#Unwrap)
function. This will continue to unwrap into subsequent errors until none exist.

**Extracting an error**

The standard library [`errors.As`](https://golang.org/pkg/errors/#As)
function can be used directly with a multierror to extract a specific error:

```go
// Assume err is a multierror value
err := somefunc()

// We want to know if "err" has a "RichErrorType" in it and extract it.
var errRich RichErrorType
if errors.As(err, &errRich) {
	// It has it, and now errRich is populated.
}
```

**Checking for an exact error value**

Some errors are returned as exact errors such as the [`ErrNotExist`](https://golang.org/pkg/os/#pkg-variables)
error in the `os` package. You can check if this error is present by using
the standard [`errors.Is`](https://golang.org/pkg/errors/#Is) function.

```go
// Assume err is a multierror value
err := somefunc()
if errors.Is(err, os.ErrNotExist) {
	// err contains os.ErrNotExist
}
```

**Returning a multierror only if there are errors**

If you build a `multierror.Error`, you can use the `ErrorOrNil` function
to return an `error` implementation only if there are errors to return:

```go
var result *multierror.Error

// ... accumulate errors here

// Return the `error` only if errors were added to the multierror, otherwise
// return nil since there are no errors.
return result.ErrorOrNil()
```
