package testassert

import (
	"encoding/json"
	"reflect"
	"testing"
)

// JSONAssertion provides fluent assertions for JSON strings.
type JSONAssertion struct {
	t   testing.TB
	got string
}

// ThatJSON creates a new [JSONAssertion] for a JSON string.
func ThatJSON(t testing.TB, got string) *JSONAssertion {
	t.Helper()
	return &JSONAssertion{t: t, got: got}
}

// Equals asserts that the JSON is semantically equal to want, ignoring formatting.
// Both strings are unmarshalled and compared with reflect.DeepEqual.
func (a *JSONAssertion) Equals(want string) *JSONAssertion {
	a.t.Helper()
	var gotVal, wantVal any
	if err := json.Unmarshal([]byte(a.got), &gotVal); err != nil {
		a.t.Errorf("Equals failed: could not unmarshal got: %v\n  got: %s", err, a.got)
		return a
	}
	if err := json.Unmarshal([]byte(want), &wantVal); err != nil {
		a.t.Errorf("Equals failed: could not unmarshal want: %v\n  want: %s", err, want)
		return a
	}
	if !reflect.DeepEqual(gotVal, wantVal) {
		a.t.Errorf("Equals failed (JSON)\n  got:  %s\n  want: %s", a.got, want)
	}
	return a
}

// Contains asserts that the top-level JSON object contains the given key with the given value.
// The value is compared after JSON unmarshalling for type consistency.
func (a *JSONAssertion) Contains(key string, value any) *JSONAssertion {
	a.t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(a.got), &m); err != nil {
		a.t.Errorf("Contains failed: could not unmarshal got as object: %v\n  got: %s", err, a.got)
		return a
	}
	actual, ok := m[key]
	if !ok {
		a.t.Errorf("Contains failed: key %q not found\n  got: %s", key, a.got)
		return a
	}
	// Marshal and unmarshal the expected value so types match (e.g., int becomes float64).
	expected, err := normalizeJSONValue(value)
	if err != nil {
		a.t.Errorf("Contains failed: could not normalize expected value: %v", err)
		return a
	}
	if !reflect.DeepEqual(actual, expected) {
		a.t.Errorf("Contains failed: key %q value mismatch\n  got:  %v\n  want: %v", key, actual, expected)
	}
	return a
}

// HasKey asserts that the top-level JSON object contains the given key.
func (a *JSONAssertion) HasKey(key string) *JSONAssertion {
	a.t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(a.got), &m); err != nil {
		a.t.Errorf("HasKey failed: could not unmarshal got as object: %v\n  got: %s", err, a.got)
		return a
	}
	if _, ok := m[key]; !ok {
		a.t.Errorf("HasKey failed: key %q not found\n  got: %s", key, a.got)
	}
	return a
}

// normalizeJSONValue round-trips a value through JSON to normalize types.
func normalizeJSONValue(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
