# stacktrace

Stracktrace package provides a convenient way of adding stacktrace to errors.

## Why is it not part if errors.New

We want to define our errors mostly as sentinel error if possible, to allow comparison. 

Creating a stack trace at the location where `New` is called is not useful.

That is why the caller has to decide when to add stack trace to the error by calling `stacktrace.New` or `stacktrace.Cause`.

## Why are errosStack and errorCause separate implementation

Most tools like Sentry extract the stacktrace from the error by checking if the error implements the `StackTrace()` method.

We want to avoid wrapping the errors multiple times with the stacktrace to avoid creating duplicates. The stacktrace should only originate at the root cause of the error

However, at the same time, we want to add annotations to certain paths when wrapping errors. 

The difference is `errorStack` fulfills the `StackTracer` interface, but `errorsCause` doesn't.
