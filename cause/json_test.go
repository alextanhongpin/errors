package cause_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/alextanhongpin/errors/cause"
	"github.com/alextanhongpin/errors/codes"
)

func roundTrip(t *testing.T, original error) *cause.Error {
	t.Helper()
	data, err := json.Marshal(cause.From(original))
	if err != nil {
		t.Fatal(err)
	}
	var restored *cause.Error
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	return restored
}

func TestJSONComparison(t *testing.T) {
	named := cause.New(codes.NotFound, "accounts.user_not_found", "original")
	for _, original := range []error{named, named.WithMessage("translated").WithDetails(map[string]any{"id": "123"}), sql.ErrNoRows, cause.ErrInternal.WithCause(sql.ErrNoRows)} {
		restored := roundTrip(t, original)
		if !errors.Is(restored, original) {
			t.Fatalf("decoded error does not match original %v", original)
		}
	}
	restored := roundTrip(t, named.WithMessage("changed"))
	if !errors.Is(restored, named) {
		t.Fatal("named identity depends on text")
	}
	for _, wrong := range []error{cause.New(codes.Internal, named.Name, "%s", named.Message), cause.New(codes.NotFound, "other", "%s", named.Message), fmt.Errorf("wrapper: %w", named)} {
		if errors.Is(restored, wrong) || errors.Is(named, wrong) {
			t.Fatalf("unexpected match for %v", wrong)
		}
	}
	plain := roundTrip(t, sql.ErrNoRows)
	if !errors.Is(plain, errors.New(sql.ErrNoRows.Error())) {
		t.Fatal("same-text ordinary errors must match")
	}
	if errors.Is(plain, errors.New("different")) || errors.Is(plain, cause.New(codes.Unknown, "named", "%s", plain.Message)) {
		t.Fatal("different ordinary or named identity matched")
	}
	if errors.Is(sql.ErrNoRows, plain) {
		t.Fatal("reverse external comparison is not preserved")
	}
	unnamed := cause.New(codes.NotFound, "", "one")
	if !errors.Is(roundTrip(t, unnamed), unnamed) || errors.Is(unnamed, unnamed.WithMessage("two")) {
		t.Fatal("unnamed structured comparison omitted message")
	}
	empty := errors.New("")
	if !errors.Is(roundTrip(t, empty), empty) {
		t.Fatal("empty ordinary text lost")
	}
}

func TestJSONTreesAndContext(t *testing.T) {
	domain := cause.New(codes.Conflict, "accounts.conflict", "conflict")
	original := fmt.Errorf("lookup context: %w", cause.ErrInternal.WithCause(errors.Join(sql.ErrNoRows, domain, errors.ErrUnsupported)))
	restored := roundTrip(t, original)
	if restored.Message != original.Error() || restored.Name != "" || restored.Cause.(*cause.Error).Name != "INTERNAL" {
		t.Fatal("wrapper context lost")
	}
	for _, target := range []error{original, sql.ErrNoRows, domain, errors.ErrUnsupported, cause.ErrInternal} {
		if !errors.Is(restored, target) {
			t.Fatalf("tree lost %v", target)
		}
	}
	// All actual wire nodes decode as *Error; native joins link their branches.
	var inspect func(error)
	inspect = func(err error) {
		if e, ok := err.(*cause.Error); ok {
			if e.Cause != nil {
				inspect(e.Cause)
			}
		} else if joined, ok := err.(interface{ Unwrap() []error }); ok {
			for _, child := range joined.Unwrap() {
				inspect(child)
			}
		} else {
			t.Fatalf("unexpected decoded node %T", err)
		}
	}
	inspect(restored)
	for _, input := range []error{original, errors.Join(sql.ErrNoRows, domain), cause.ErrInternal.WithCause(errors.Join(domain, domain))} {
		before, err := json.Marshal(cause.From(input))
		if err != nil {
			t.Fatal(err)
		}
		after, err := json.Marshal(roundTrip(t, input))
		if err != nil || !bytes.Equal(before, after) {
			t.Fatalf("tree shape changed:\n%s\n%s\n%v", before, after, err)
		}
	}
}

func TestJSONFieldsAndEmbedding(t *testing.T) {
	original := cause.New(codes.NotFound, "accounts.missing", "missing").WithStack().WithAttrs(slog.String("secret", "hidden")).WithDetails(map[string]any{"id": json.Number("9007199254740993")})
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("stack")) || bytes.Contains(data, []byte("secret")) {
		t.Fatal("diagnostic fields leaked")
	}
	restored := roundTrip(t, original)
	if restored.Stack != "" || len(restored.Attrs) != 0 || restored.Details["id"] != json.Number("9007199254740993") {
		t.Fatal("fields or precision changed")
	}
	// Details remain enrichable even when an absent details field decoded to nil.
	noDetails := roundTrip(t, cause.ErrNotFound).WithDetails(map[string]any{"field": "id"})
	if noDetails.Details["field"] != "id" {
		t.Fatal("decoded details cannot be enriched")
	}
	for _, value := range []any{struct {
		Error *cause.Error `json:"error"`
	}{original}, map[string]error{"error": original}, []error{original}} {
		encoded, err := json.Marshal(value)
		if err != nil || !bytes.Contains(encoded, data) {
			t.Fatalf("embedding failed: %s %v", encoded, err)
		}
	}
	reordered := []byte(`{"message":"translated","name":"accounts.missing","code":"not_found","future":true}`)
	var decoded cause.Error
	if err := json.Unmarshal(reordered, &decoded); err != nil || !errors.Is(&decoded, original) {
		t.Fatal("formatting changed identity", err)
	}
	plain, err := json.Marshal(cause.From(sql.ErrNoRows))
	if err != nil || string(plain) != `{"code":"unknown","name":"","message":"sql: no rows in result set"}` {
		t.Fatalf("unexpected ordinary representation: %s %v", plain, err)
	}
}

func TestJSONValidationIsAtomic(t *testing.T) {
	for _, payload := range []string{
		`{}`, `[]`, `true`, `{"code":17,"name":"UNKNOWN","message":"old"}`,
		`{"code":17,"name":"x","message":"m"}`,
		`{"code":"future","name":"x","message":"m"}`,
		`{"code":"unknown","name":null,"message":"m"}`,
		`{"code":"unknown","name":"x"}`,
		`{"code":"unknown","name":"x","message":null}`,
		`{"code":"unknown","name":"x","message":"m","details":null}`,
		`{"code":"unknown","name":"x","message":"m","details":[]}`,
		`{"code":"unknown","name":"x","message":"m","causes":null}`,
		`{"code":"unknown","name":"x","message":"m","causes":[null]}`,
		`{"code":"unknown","name":"x","message":"m"} {}`,
	} {
		t.Run(payload, func(t *testing.T) {
			original := cause.ErrConflict.WithStack().WithCause(sql.ErrNoRows).WithDetails(map[string]any{"safe": "keep"})
			previousCause := original.Cause
			if err := json.Unmarshal([]byte(payload), original); err == nil {
				t.Fatalf("accepted %s", payload)
			}
			if original.Code != codes.Conflict || original.Name != "CONFLICT" || original.Cause != previousCause || original.Stack == "" || original.Details["safe"] != "keep" {
				t.Fatal("failed decoding mutated receiver")
			}
		})
	}
	existing := cause.ErrInternal.WithStack().WithAttrs(slog.String("private", "data")).WithCause(sql.ErrNoRows)
	if err := json.Unmarshal([]byte(`{"code":"not_found","name":"new","message":"m"}`), existing); err != nil {
		t.Fatal(err)
	}
	if existing.Cause != nil || existing.Stack != "" || len(existing.Attrs) != 0 {
		t.Fatal("successful decoding retained stale fields")
	}
}

type cycleError struct{ child error }

func (e *cycleError) Error() string { return "cycle" }
func (e *cycleError) Unwrap() error { return e.child }

type branches []error

func (e branches) Error() string   { return "branches" }
func (e branches) Unwrap() []error { return []error(e) }

func TestJSONSafetyLimits(t *testing.T) {
	cycle := &cycleError{}
	cycle.child = cycle
	structuredCycle := cause.ErrInternal.Clone()
	structuredCycle.Cause = structuredCycle
	sliceCycle := make(branches, 1)
	sliceCycle[0] = sliceCycle
	invalid := []*cause.Error{structuredCycle, cause.ErrInternal.WithCause(cycle), cause.ErrInternal.WithCause(sliceCycle), cause.ErrInternal.WithCause(branches{nil}), cause.ErrInternal.WithCause((*cause.Error)(nil)), cause.New(codes.Code(999), "bad", "bad"), cause.ErrInternal.WithDetails(map[string]any{"bad": make(chan int)})}
	detailCycle := map[string]any{}
	detailCycle["self"] = detailCycle
	invalid = append(invalid, cause.ErrInternal.WithDetails(detailCycle))
	for _, original := range invalid {
		if _, err := json.Marshal(original); err == nil {
			t.Fatal("invalid tree encoded")
		}
	}
	deep := cause.ErrNotFound.Clone()
	for range 63 {
		deep = cause.ErrInternal.WithCause(deep)
	}
	data, err := json.Marshal(deep)
	if err != nil {
		t.Fatal("64-level boundary rejected", err)
	}
	var restored cause.Error
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal("64-level decode rejected", err)
	}
	if _, err := json.Marshal(cause.ErrInternal.WithCause(deep)); err == nil {
		t.Fatal("depth limit bypassed")
	}
	child := string(data)
	overDeep := `{"code":"internal","name":"INTERNAL","message":"m","causes":[` + child + `]}`
	if err := json.Unmarshal([]byte(overDeep), &restored); err == nil {
		t.Fatal("decode depth limit bypassed")
	}
	causes := make([]error, 1023)
	for i := range causes {
		causes[i] = errors.New("leaf")
	}
	data, err = json.Marshal(cause.From(errors.Join(causes...)))
	if err != nil {
		t.Fatal("1024-node boundary rejected", err)
	}
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal("1024-node decode rejected", err)
	}
	causes = append(causes, errors.New("leaf"))
	if _, err := json.Marshal(cause.From(errors.Join(causes...))); err == nil {
		t.Fatal("node limit bypassed")
	}
	leaf := `{"code":"unknown","name":"","message":"leaf"}`
	overNodes := `{"code":"unknown","name":"","message":"branches","causes":[` + strings.TrimSuffix(strings.Repeat(leaf+",", 1024), ",") + `]}`
	if err := json.Unmarshal([]byte(overNodes), &restored); err == nil {
		t.Fatal("decode node limit bypassed")
	}
	if _, err := json.Marshal(cause.ErrInternal.WithMessage("%s", strings.Repeat("x", 1<<20))); err == nil {
		t.Fatal("byte limit bypassed")
	}
	oversized := `{"code":"unknown","name":"","message":"` + strings.Repeat("x", 1<<20) + `"}`
	if err := json.Unmarshal([]byte(oversized), &restored); err == nil {
		t.Fatal("decode byte limit bypassed")
	}
}

func TestFromNilAndConcurrentJSON(t *testing.T) {
	original := cause.ErrNotFound.Clone()
	if cause.From(original) != original {
		t.Fatal("From must preserve structured pointer")
	}
	for _, input := range []error{nil, (*cause.Error)(nil)} {
		if cause.From(input) != nil {
			t.Fatal("typed nil normalized to nonnil")
		}
		data, err := json.Marshal(cause.From(input))
		if err != nil || string(data) != "null" {
			t.Fatal("nil encoding failed", err)
		}
		restored := cause.ErrInternal.Clone()
		if err := json.Unmarshal(data, &restored); err != nil || restored != nil {
			t.Fatal("null pointer decoding failed", err)
		}
	}
	var nilError *cause.Error
	if nilError.Is(sql.ErrNoRows) || errors.Is(nilError, sql.ErrNoRows) {
		t.Fatal("nil matched")
	}
	if err := nilError.UnmarshalJSON([]byte("null")); err == nil {
		t.Fatal("nil receiver accepted")
	}
	value := *original
	if err := json.Unmarshal([]byte("null"), &value); err != nil || value.Cause != nil || value.Name != "" || value.Code != 0 {
		t.Fatal("null value not reset", err)
	}
	var wg sync.WaitGroup
	immutable := cause.ErrInternal.WithCause(sql.ErrNoRows)
	for range 16 {
		wg.Go(func() {
			data, err := json.Marshal(immutable)
			if err != nil {
				t.Error(err)
				return
			}
			var restored *cause.Error
			if err := json.Unmarshal(data, &restored); err != nil || !errors.Is(restored, sql.ErrNoRows) {
				t.Error("concurrent round trip failed", err)
			}
		})
	}
	wg.Wait()
}

func FuzzErrorJSON(f *testing.F) {
	f.Add([]byte(`{"code":"unknown","name":"","message":"m"}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, data []byte) {
		var original *cause.Error
		if err := json.Unmarshal(data, &original); err != nil {
			return
		}
		encoded, err := json.Marshal(original)
		if err != nil {
			return
		} // Canonical JSON may exceed the output byte limit.
		var restored *cause.Error
		if err := json.Unmarshal(encoded, &restored); err != nil {
			t.Fatal(err)
		}
		if original == nil {
			if restored != nil {
				t.Fatal("nil changed")
			}
			return
		}
		if !errors.Is(restored, original) {
			t.Fatal("round trip lost comparison")
		}
	})
}
