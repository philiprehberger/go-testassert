package testassert

import (
	"testing"
)

// MapAssertion provides fluent assertions for map values.
type MapAssertion[K comparable, V any] struct {
	t   testing.TB
	got map[K]V
}

// ThatMap creates a new [MapAssertion] for a map value.
func ThatMap[K comparable, V any](t testing.TB, got map[K]V) *MapAssertion[K, V] {
	t.Helper()
	return &MapAssertion[K, V]{t: t, got: got}
}

// HasKey asserts that the map contains the given key.
func (a *MapAssertion[K, V]) HasKey(key K) *MapAssertion[K, V] {
	a.t.Helper()
	if _, ok := a.got[key]; !ok {
		a.t.Errorf("HasKey failed\n  key: %v\n  map: %v", key, a.got)
	}
	return a
}

// HasLen asserts that the map has length n.
func (a *MapAssertion[K, V]) HasLen(n int) *MapAssertion[K, V] {
	a.t.Helper()
	if len(a.got) != n {
		a.t.Errorf("HasLen failed\n  got length: %d\n  want length: %d\n  value: %v", len(a.got), n, a.got)
	}
	return a
}

// IsEmpty asserts that the map is empty.
func (a *MapAssertion[K, V]) IsEmpty() *MapAssertion[K, V] {
	a.t.Helper()
	if len(a.got) != 0 {
		a.t.Errorf("IsEmpty failed\n  got length: %d\n  value: %v", len(a.got), a.got)
	}
	return a
}

// IsNotEmpty asserts that the map is not empty.
func (a *MapAssertion[K, V]) IsNotEmpty() *MapAssertion[K, V] {
	a.t.Helper()
	if len(a.got) == 0 {
		a.t.Errorf("IsNotEmpty failed\n  got: empty map")
	}
	return a
}
