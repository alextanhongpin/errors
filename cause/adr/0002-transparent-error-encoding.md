# ADR 002: Transparent Error Encoding

**Status:** Superseded by [ADR 001](0001-error-serialization-and-comparison.md)

**Date:** 2026-10-08

**Related:** [ADR 001: Error Serialization and Portable Comparison](0001-error-serialization-and-comparison.md)

This historical proposal is not the current API. The revised ADR 001 replaces
codec-backed adapters with JSON methods on `*cause.Error` and `From`. The adapter
implementation has been removed.

## Context

The codec currently returns encoded bytes. Applications must explicitly encode errors before embedding them in JSON responses or passing them to database operations.

We want standard JSON and database APIs to perform that encoding automatically, while retaining ADR 001's error representation, comparison behavior, and safety limits.

A `MarshalJSON` method on `Error` cannot receive the application's codec configuration or determine whether its destination is a client or database.

## Decision

We will add lightweight, codec-backed adapters. Applications select the representation once; standard interfaces handle subsequent encoding.

### JSON integration

Add these public constructors:

```go
func (c *Codec) Client(public *ClientError) ClientDocument
func (c *Codec) Storage(err error) *StorageDocument
```

Both document types implement `json.Marshaler`, delegating directly to the corresponding existing codec method. This supports root values, struct fields, maps, and slices. See the [Go JSON interfaces](https://pkg.go.dev/encoding/json#Marshaler).

```go
response := struct {
    Error cause.ClientDocument `json:"error"`
}{
    Error: codec.Client(publicError),
}

err := json.NewEncoder(w).Encode(response)
data, err := json.Marshal(codec.Storage(originalError))
```

- `Client(nil)` retains the existing generic internal-error fallback.
- `Storage(nil)` encodes as JSON `null`.
- Constructors defer encoding and validation until the standard API invokes the adapter.
- Adapters retain their inputs; they do not promise snapshots. Inputs must not be concurrently mutated during encoding.
- Codec limits apply to each error document, independently of the enclosing response.
- Existing byte-oriented codec methods remain supported.
- No JSON methods are added to `Error`, and no global codec is introduced.

### Database integration

`StorageDocument` also implements `driver.Valuer` and `sql.Scanner`. It exposes:

```go
func (d *StorageDocument) Err() error
```

`Err` returns the original supplied error, or the decoded `*Record` after a successful read.

```go
stored := codec.Storage(originalError)
_, err := db.ExecContext(ctx,
    "INSERT INTO failures (error) VALUES ($1)", stored)

loaded := codec.Storage(nil)
err := db.QueryRowContext(ctx,
    "SELECT error FROM failures WHERE id = $1", id).Scan(loaded)

matches := errors.Is(loaded.Err(), expectedError)
```

- `Value` delegates to `EncodeStorage` and returns JSON as a string. A nil underlying error returns SQL `NULL`.
- `Scan` accepts JSON supplied as `string` or `[]byte`, plus SQL `NULL`.
- JSON input delegates to `Decode`, preserving configured mappings, validation, and limits.
- SQL `NULL` and JSON `null` both produce a nil underlying error.
- Unsupported source types or invalid documents return errors without replacing the previous value.
- Driver-provided byte slices are decoded immediately and are not retained.

These interfaces integrate with `database/sql`; database-specific column types and parameter casts remain application responsibilities. No database driver dependency is added. See [Valuer](https://pkg.go.dev/database/sql/driver#Valuer) and [Scanner](https://pkg.go.dev/database/sql#Scanner).

### Decoding and configuration

`StorageDocument` implements `json.Unmarshaler`, also delegating to `Decode`:

```go
document := codec.Storage(nil)
err := json.Unmarshal(data, document)
```

Successful decoding replaces its underlying error; failed decoding leaves it unchanged. Callers initialize decoding destinations through the codec so sentinel mappings remain explicit.

Unconfigured document zero values return configuration errors when their methods are invoked. `ClientDocument` is encoding-only; callers use `Codec.Decode` when inspecting received client documents.

## Alternatives Considered

- **JSON methods directly on `Error`:** Cannot express destination-specific projection or application-owned codec configuration.
- **Global default codec:** Hides sentinel mappings and introduces shared mutable configuration.
- **Pre-encoded `json.RawMessage`:** Useful today, but still requires manual encoding and provides no database adapter.
- **Adapters backed by the codec:** Chosen because they improve integration while keeping one serialization implementation.

## Consequences

Applications can embed error documents directly in JSON responses and pass storage documents to database operations. All existing serialization and comparison guarantees remain centralized in the codec.

Applications must still explicitly select client or storage intent. Adapter errors surface during encoding or database operations, rather than construction. Storage adapters remain inappropriate for public responses because they preserve diagnostics.

This is an additive API change. It does not change the wire format or restore legacy support.

## Validation and Documentation

Verify:

- JSON output matches existing codec output at the root and inside structs, maps, and slices.
- Client fallback, field disclosure rules, and storage cause trees remain unchanged.
- Cycles, unsupported details, and limit violations propagate through standard encoding.
- Configured sentinel matching survives JSON decoding and `Value`/`Scan` round trips.
- SQL `NULL`, JSON `null`, strings, and byte slices follow the specified behavior.
- Failed reads preserve previous state; unsupported scan types and unconfigured adapters return errors.
- Concurrent adapter encoding follows the codec's existing input-mutation constraints.

Add executable JSON and database-interface examples, and run package tests, race checks, and `go vet` when implementing the adapters.

This proposal is retained as decision history. Its adapter API and validation plan have been superseded by the revised ADR 001.
