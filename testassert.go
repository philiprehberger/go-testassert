// Package testassert provides fluent, type-safe test assertions for Go.
//
// It offers a chainable API built with generics for asserting values in tests.
// Each assertion type is constructed with a That* function and exposes
// domain-specific methods that call t.Errorf on failure with readable messages.
package testassert

import (
	"cmp"
	"errors"
	"reflect"
	"regexp"
	"testing"
)

// Assertion provides fluent assertions for any type using reflect.DeepEqual.
type Assertion[T any] struct {
	t   testing.TB
	got T
}

// That creates a new [Assertion] for the given value.
func That[T any](t testing.TB, got T) *Assertion[T] {
	t.Helper()
	return &Assertion[T]{t: t, got: got}
}

// Equals asserts that got is deeply equal to want.
func (a *Assertion[T]) Equals(want T) *Assertion[T] {
	a.t.Helper()
	if !reflect.DeepEqual(a.got, want) {
		a.t.Errorf("Equals failed\n  got:  %v\n  want: %v", a.got, want)
	}
	return a
}

// NotEquals asserts that got is not deeply equal to want.
func (a *Assertion[T]) NotEquals(want T) *Assertion[T] {
	a.t.Helper()
	if reflect.DeepEqual(a.got, want) {
		a.t.Errorf("NotEquals failed\n  got: %v\n  should differ from: %v", a.got, want)
	}
	return a
}

// IsNil asserts that the value is nil.
// It handles both typed and untyped nils using reflection.
func (a *Assertion[T]) IsNil() *Assertion[T] {
	a.t.Helper()
	if !isNil(a.got) {
		a.t.Errorf("IsNil failed\n  got: %v\n  want: nil", a.got)
	}
	return a
}

// IsNotNil asserts that the value is not nil.
func (a *Assertion[T]) IsNotNil() *Assertion[T] {
	a.t.Helper()
	if isNil(a.got) {
		a.t.Errorf("IsNotNil failed\n  got: nil\n  want: non-nil value")
	}
	return a
}

// OrderedAssertion provides fluent assertions for ordered types that support
// comparison operators.
type OrderedAssertion[T cmp.Ordered] struct {
	t   testing.TB
	got T
}

// ThatOrdered creates a new [OrderedAssertion] for a value with an ordered type.
func ThatOrdered[T cmp.Ordered](t testing.TB, got T) *OrderedAssertion[T] {
	t.Helper()
	return &OrderedAssertion[T]{t: t, got: got}
}

// Equals asserts that got equals want.
func (a *OrderedAssertion[T]) Equals(want T) *OrderedAssertion[T] {
	a.t.Helper()
	if a.got != want {
		a.t.Errorf("Equals failed\n  got:  %v\n  want: %v", a.got, want)
	}
	return a
}

// NotEquals asserts that got does not equal want.
func (a *OrderedAssertion[T]) NotEquals(want T) *OrderedAssertion[T] {
	a.t.Helper()
	if a.got == want {
		a.t.Errorf("NotEquals failed\n  got: %v\n  should differ from: %v", a.got, want)
	}
	return a
}

// IsGreaterThan asserts that got > v.
func (a *OrderedAssertion[T]) IsGreaterThan(v T) *OrderedAssertion[T] {
	a.t.Helper()
	if !(a.got > v) {
		a.t.Errorf("IsGreaterThan failed\n  got:  %v\n  want: > %v", a.got, v)
	}
	return a
}

// IsLessThan asserts that got < v.
func (a *OrderedAssertion[T]) IsLessThan(v T) *OrderedAssertion[T] {
	a.t.Helper()
	if !(a.got < v) {
		a.t.Errorf("IsLessThan failed\n  got:  %v\n  want: < %v", a.got, v)
	}
	return a
}

// IsGreaterOrEqual asserts that got >= v.
func (a *OrderedAssertion[T]) IsGreaterOrEqual(v T) *OrderedAssertion[T] {
	a.t.Helper()
	if !(a.got >= v) {
		a.t.Errorf("IsGreaterOrEqual failed\n  got:  %v\n  want: >= %v", a.got, v)
	}
	return a
}

// IsLessOrEqual asserts that got <= v.
func (a *OrderedAssertion[T]) IsLessOrEqual(v T) *OrderedAssertion[T] {
	a.t.Helper()
	if !(a.got <= v) {
		a.t.Errorf("IsLessOrEqual failed\n  got:  %v\n  want: <= %v", a.got, v)
	}
	return a
}

// StringAssertion provides fluent assertions for string values.
type StringAssertion struct {
	t   testing.TB
	got string
}

// ThatString creates a new [StringAssertion] for a string value.
func ThatString(t testing.TB, got string) *StringAssertion {
	t.Helper()
	return &StringAssertion{t: t, got: got}
}

// Equals asserts that the string equals want.
func (a *StringAssertion) Equals(want string) *StringAssertion {
	a.t.Helper()
	if a.got != want {
		a.t.Errorf("Equals failed\n  got:  %q\n  want: %q", a.got, want)
	}
	return a
}

// Contains asserts that the string contains substr.
func (a *StringAssertion) Contains(substr string) *StringAssertion {
	a.t.Helper()
	if !contains(a.got, substr) {
		a.t.Errorf("Contains failed\n  got:    %q\n  substr: %q", a.got, substr)
	}
	return a
}

// HasPrefix asserts that the string starts with prefix.
func (a *StringAssertion) HasPrefix(prefix string) *StringAssertion {
	a.t.Helper()
	if !hasPrefix(a.got, prefix) {
		a.t.Errorf("HasPrefix failed\n  got:    %q\n  prefix: %q", a.got, prefix)
	}
	return a
}

// HasSuffix asserts that the string ends with suffix.
func (a *StringAssertion) HasSuffix(suffix string) *StringAssertion {
	a.t.Helper()
	if !hasSuffix(a.got, suffix) {
		a.t.Errorf("HasSuffix failed\n  got:    %q\n  suffix: %q", a.got, suffix)
	}
	return a
}

// IsEmpty asserts that the string is empty.
func (a *StringAssertion) IsEmpty() *StringAssertion {
	a.t.Helper()
	if a.got != "" {
		a.t.Errorf("IsEmpty failed\n  got: %q\n  want: empty string", a.got)
	}
	return a
}

// HasLen asserts that the string has length n.
func (a *StringAssertion) HasLen(n int) *StringAssertion {
	a.t.Helper()
	if len(a.got) != n {
		a.t.Errorf("HasLen failed\n  got length: %d\n  want length: %d\n  value: %q", len(a.got), n, a.got)
	}
	return a
}

// Matches asserts that the string matches the regular expression pattern.
func (a *StringAssertion) Matches(pattern string) *StringAssertion {
	a.t.Helper()
	matched, err := regexp.MatchString(pattern, a.got)
	if err != nil {
		a.t.Errorf("Matches failed: invalid pattern %q: %v", pattern, err)
		return a
	}
	if !matched {
		a.t.Errorf("Matches failed\n  got:     %q\n  pattern: %q", a.got, pattern)
	}
	return a
}

// SliceAssertion provides fluent assertions for slice values.
type SliceAssertion[T any] struct {
	t   testing.TB
	got []T
}

// ThatSlice creates a new [SliceAssertion] for a slice value.
func ThatSlice[T any](t testing.TB, got []T) *SliceAssertion[T] {
	t.Helper()
	return &SliceAssertion[T]{t: t, got: got}
}

// HasLen asserts that the slice has length n.
func (a *SliceAssertion[T]) HasLen(n int) *SliceAssertion[T] {
	a.t.Helper()
	if len(a.got) != n {
		a.t.Errorf("HasLen failed\n  got length: %d\n  want length: %d\n  value: %v", len(a.got), n, a.got)
	}
	return a
}

// IsEmpty asserts that the slice is empty.
func (a *SliceAssertion[T]) IsEmpty() *SliceAssertion[T] {
	a.t.Helper()
	if len(a.got) != 0 {
		a.t.Errorf("IsEmpty failed\n  got length: %d\n  value: %v", len(a.got), a.got)
	}
	return a
}

// IsNotEmpty asserts that the slice is not empty.
func (a *SliceAssertion[T]) IsNotEmpty() *SliceAssertion[T] {
	a.t.Helper()
	if len(a.got) == 0 {
		a.t.Errorf("IsNotEmpty failed\n  got: empty slice")
	}
	return a
}

// Contains asserts that the slice contains elem using reflect.DeepEqual.
func (a *SliceAssertion[T]) Contains(elem T) *SliceAssertion[T] {
	a.t.Helper()
	for _, v := range a.got {
		if reflect.DeepEqual(v, elem) {
			return a
		}
	}
	a.t.Errorf("Contains failed\n  slice: %v\n  element: %v", a.got, elem)
	return a
}

// ErrorAssertion provides fluent assertions for error values.
type ErrorAssertion struct {
	t   testing.TB
	got error
}

// ThatError creates a new [ErrorAssertion] for an error value.
func ThatError(t testing.TB, got error) *ErrorAssertion {
	t.Helper()
	return &ErrorAssertion{t: t, got: got}
}

// IsNil asserts that the error is nil.
func (a *ErrorAssertion) IsNil() *ErrorAssertion {
	a.t.Helper()
	if a.got != nil {
		a.t.Errorf("IsNil failed\n  got: %v\n  want: nil", a.got)
	}
	return a
}

// IsNotNil asserts that the error is not nil.
func (a *ErrorAssertion) IsNotNil() *ErrorAssertion {
	a.t.Helper()
	if a.got == nil {
		a.t.Errorf("IsNotNil failed\n  got: nil\n  want: non-nil error")
	}
	return a
}

// Is asserts that errors.Is(got, target) returns true.
func (a *ErrorAssertion) Is(target error) *ErrorAssertion {
	a.t.Helper()
	if !errorsIs(a.got, target) {
		a.t.Errorf("Is failed\n  got:    %v\n  target: %v", a.got, target)
	}
	return a
}

// Contains asserts that the error message contains msg.
func (a *ErrorAssertion) Contains(msg string) *ErrorAssertion {
	a.t.Helper()
	if a.got == nil {
		a.t.Errorf("Contains failed\n  got: nil error\n  want message containing: %q", msg)
		return a
	}
	if !contains(a.got.Error(), msg) {
		a.t.Errorf("Contains failed\n  got:  %q\n  want: substring %q", a.got.Error(), msg)
	}
	return a
}

// As asserts that errors.As(got, target) returns true.
// The target must be a pointer to the desired error type.
func (a *ErrorAssertion) As(target any) *ErrorAssertion {
	a.t.Helper()
	if !errorsAs(a.got, target) {
		a.t.Errorf("As failed\n  got:  %v (%T)\n  want: assignable to %T", a.got, a.got, target)
	}
	return a
}

// isNil checks whether v is nil, handling both typed and untyped nils.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// contains reports whether s contains substr.
// Wrapper to avoid importing strings.
func contains(s, substr string) bool {
	return len(substr) == 0 || findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// hasPrefix reports whether s starts with prefix.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// hasSuffix reports whether s ends with suffix.
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// errorsIs delegates to errors.Is from the standard library.
func errorsIs(err, target error) bool {
	return errors.Is(err, target)
}

// errorsAs delegates to errors.As from the standard library.
func errorsAs(err error, target any) bool {
	return errors.As(err, target)
}
