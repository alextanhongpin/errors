# ADR 003: Marshal Error Values Consistently

**Status:** Accepted

**Date:** 2026-10-08

**Related:** [ADR 001: Simple Error JSON and Comparison](0001-error-serialization-and-comparison.md)

## Context

`Error.MarshalJSON` currently has a pointer receiver. Go's `encoding/json` v1 only invokes pointer-receiver marshal methods when the value is addressable. A non-addressable value, such as a `cause.Error` stored in `map[string]cause.Error`, can therefore bypass the package's JSON representation and expose the struct's exported fields, including `Attrs` and `Stack`.

This makes the serialized format depend on where the same error value is stored. The Go `encoding/json` documentation describes this pointer-receiver behavior in its [v1 and v2 differences](https://pkg.go.dev/encoding/json#Marshal).

## Decision

We will give `Error.MarshalJSON` a value receiver and route that copied value through the existing error encoder. Both `Error` and `*Error` values will then use the same wire representation in roots, struct fields, maps, and slices.

The method will continue to:

- Encode the existing JSON schema and preserve error causes.
- Omit logging attributes and stack traces.
- Enforce the existing cycle and document-size limits.
- Leave `UnmarshalJSON` on a pointer receiver, since decoding must update its destination.
- Preserve standard `encoding/json` handling of nil `*Error` pointers as JSON `null`.

The versionless wire schema and comparison rules are defined by ADR 001.

## Consequences

Encoding is consistent for addressable and non-addressable error values, so value storage cannot accidentally reveal fields omitted by the error schema.

The value receiver copies the top-level `Error` before traversing its details and causes. Those referenced values must not be mutated concurrently during encoding, matching the existing input-mutation constraint. Directly invoking `MarshalJSON` on a nil `*Error` is not supported; standard `json.Marshal` already emits null for nil pointers.

## Validation

Verify that:

- `json.Marshal` produces identical JSON for equivalent `Error` and `*Error` values.
- Value errors in maps, slices, and structs do not expose `Attrs` or `Stack`.
- Pointer errors retain current output, cause comparison, and nil-pointer behavior.
- Cycles, invalid details, and size limits still return errors through every encoding shape.

## Implementation

`Error.MarshalJSON` uses a value receiver and passes a copy through the error
encoder. `Error.UnmarshalJSON` remains a pointer receiver. README usage now
documents value and pointer encoding. The wire schema is defined by ADR 001.
