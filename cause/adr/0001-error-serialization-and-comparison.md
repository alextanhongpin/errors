# ADR 001: Simple Error JSON and Comparison

**Status:** Accepted

**Date:** 2026-10-08

**Related:** [ADR 002](0002-transparent-error-encoding.md) is superseded by this simpler design. [ADR 003](0003-json-value-marshaler.md) makes JSON encoding consistent for `Error` values and pointers.

## Context

We need to marshal errors to JSON, unmarshal them, and compare the decoded error with the original using `errors.Is`.

We accept a simpler guarantee: ordinary errors are equivalent when their error messages match. This means an unrelated `errors.New` with the same text as `sql.ErrNoRows` will also match.

The earlier codec, sentinel registry, separate projections, and database adapters are unnecessary for this requirement.

## Decision

We implement `MarshalJSON` on `cause.Error` values and `UnmarshalJSON` on `*cause.Error`, using one representation for clients and database storage.

### Normalize errors into one representation

Provide one helper:

```go
func From(err error) *Error
```

- A `*cause.Error` is returned directly.
- An ordinary error becomes a `*cause.Error` with `Code: codes.Unknown`, an empty `Name`, and `Message: err.Error()`.
- Nil and typed-nil errors return nil.
- Normalization preserves wrapped causes and joined branches, inspecting each node directly rather than searching through wrappers.

Top-level ordinary errors require this helper because the package cannot add JSON methods to external types such as `sql.ErrNoRows`:

```go
original := sql.ErrNoRows

data, err := json.Marshal(cause.From(original))
if err != nil {
    return err
}

var restored *cause.Error
if err := json.Unmarshal(data, &restored); err != nil {
    return err
}

fmt.Println(errors.Is(restored, original)) // true
```

Ordinary errors nested inside `cause.Error` are normalized automatically:

```go
original := cause.ErrNotFound.WithCause(sql.ErrNoRows)
data, err := json.Marshal(original)
```

### Compare normalized fields

`Error.Is` compares the current node against the normalized fields of the supplied target:

| Representation | Comparison |
|---|---|
| Nonempty `Name` | Exact `Code + Name` |
| Empty `Name` | Exact `Code + Name + Message` |

Consequently:

- Named structured errors remain comparable despite message changes.
- Ordinary errors compare by their messages through the common normalized representation.
- Details and causes do not affect a node's identity.
- `errors.Is` traverses the source's causes; `Error.Is` does not unwrap the target.
- Comparison uses fields, not serialized JSON bytes or object key order.

The guarantee is `errors.Is(restored, original)`. Reverse comparison with an external Go sentinel is not guaranteed, because its implementation remains outside this package.

### JSON contract

Use string `code`, `name`, `message`, optional `details`, and ordered `causes`. The document has no version field.

An ordinary error becomes:

```json
{
  "code": "unknown",
  "name": "",
  "message": "sql: no rows in result set"
}
```

- All decoded wire nodes use `*cause.Error`; there is no separate `Record` type.
- Single causes populate `Error.Cause`; multiple causes use `errors.Join` behind an internal wrapper that keeps wire branches stable across repeated round trips.
- Stacks and logging attributes are excluded from JSON.
- Details retain JSON-compatible values; decoded numbers use `json.Number`.
- Nil pointers encode as `null`. Decoding into a `*Error` variable supports nil; decoding `null` into an existing `Error` value resets it.
- Invalid documents return errors without changing the destination. Unknown fields, including any `version` field, are ignored.
- Retain fixed internal safeguards of 64 error levels, 1,024 node occurrences, 1 MiB per document, and cycle detection. Limits apply to the error document before outer formatting, not its envelope or all encoding allocations.

## Consequences

Serialization becomes ordinary Go JSON usage, and comparison works without configuring mappings.

Message collisions are an accepted limitation for unnamed errors. Changing an ordinary error's message changes its identity. Original concrete types and arbitrary custom `Is` or `As` behavior are not reconstructed.

One JSON representation serves both clients and storage. Applications are responsible for selecting safe messages, details, and causes before returning them to clients.

## Validation

Verify:

- Named errors round-trip and match despite differing messages or details.
- Different structured codes or names do not match.
- `sql.ErrNoRows` round-trips and matches the original sentinel.
- A different error with identical text intentionally matches; different text does not.
- Wrapped targets are not searched during shallow comparison.
- Nested wrappers and every joined branch remain searchable and stable across repeated round trips.
- Nil values, malformed documents, unsupported details, cycles, and limits behave correctly.
- JSON omits stacks and logging attributes.
- Concurrent encoding of immutable input and decoding into separate destinations is race-free.

## Compatibility

The codec and adapter APIs, including mappings, identity interfaces, `Record`, client/storage documents, and database interfaces, are removed. Examples, tests, README, and changelog use direct JSON methods and `From`.

Database persistence stores the resulting JSON bytes; database-specific adapters remain outside this package. Legacy numeric-code documents are rejected because `code` must be a string. The JSON schema has no version discriminator, so compatible schema evolution must preserve existing field meanings. Applications using earlier text-free sentinel matching must account for the newly accepted text collisions.
