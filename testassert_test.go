package testassert

import (
	"errors"
	"fmt"
	"testing"
)

// mockT captures test failures without failing the real test.
type mockT struct {
	testing.TB
	failed   bool
	messages []string
}

func newMockT() *mockT { return &mockT{} }

func (m *mockT) Helper() {}

func (m *mockT) Errorf(format string, args ...any) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

// --- That (generic Assertion) ---

func TestThat_Equals(t *testing.T) {
	That(t, 42).Equals(42)
	That(t, "hello").Equals("hello")
	That(t, []int{1, 2, 3}).Equals([]int{1, 2, 3})

	mt := newMockT()
	That(mt, 1).Equals(2)
	if !mt.failed {
		t.Error("expected failure when values differ")
	}
}

func TestThat_NotEquals(t *testing.T) {
	That(t, 42).NotEquals(99)
	That(t, "hello").NotEquals("world")

	mt := newMockT()
	That(mt, 42).NotEquals(42)
	if !mt.failed {
		t.Error("expected failure when values are equal")
	}
}

func TestThat_IsNil(t *testing.T) {
	var p *int
	That[*int](t, nil).IsNil()
	That(t, p).IsNil()

	mt := newMockT()
	v := 5
	That(mt, &v).IsNil()
	if !mt.failed {
		t.Error("expected failure for non-nil pointer")
	}
}

func TestThat_IsNotNil(t *testing.T) {
	v := 5
	That(t, &v).IsNotNil()

	mt := newMockT()
	var p *int
	That(mt, p).IsNotNil()
	if !mt.failed {
		t.Error("expected failure for nil pointer")
	}
}

// --- OrderedAssertion ---

func TestThatOrdered_Equals(t *testing.T) {
	ThatOrdered(t, 10).Equals(10)
	ThatOrdered(t, "abc").Equals("abc")

	mt := newMockT()
	ThatOrdered(mt, 10).Equals(20)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatOrdered_NotEquals(t *testing.T) {
	ThatOrdered(t, 10).NotEquals(20)

	mt := newMockT()
	ThatOrdered(mt, 10).NotEquals(10)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatOrdered_GreaterThan(t *testing.T) {
	ThatOrdered(t, 10).IsGreaterThan(5)
	ThatOrdered(t, 3.14).IsGreaterThan(2.71)

	mt := newMockT()
	ThatOrdered(mt, 5).IsGreaterThan(10)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatOrdered_LessThan(t *testing.T) {
	ThatOrdered(t, 5).IsLessThan(10)

	mt := newMockT()
	ThatOrdered(mt, 10).IsLessThan(5)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatOrdered_GreaterOrEqual(t *testing.T) {
	ThatOrdered(t, 10).IsGreaterOrEqual(10)
	ThatOrdered(t, 10).IsGreaterOrEqual(5)

	mt := newMockT()
	ThatOrdered(mt, 5).IsGreaterOrEqual(10)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatOrdered_LessOrEqual(t *testing.T) {
	ThatOrdered(t, 10).IsLessOrEqual(10)
	ThatOrdered(t, 5).IsLessOrEqual(10)

	mt := newMockT()
	ThatOrdered(mt, 10).IsLessOrEqual(5)
	if !mt.failed {
		t.Error("expected failure")
	}
}

// --- StringAssertion ---

func TestThatString_Equals(t *testing.T) {
	ThatString(t, "hello").Equals("hello")

	mt := newMockT()
	ThatString(mt, "hello").Equals("world")
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_Contains(t *testing.T) {
	ThatString(t, "hello world").Contains("world")

	mt := newMockT()
	ThatString(mt, "hello").Contains("xyz")
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_HasPrefix(t *testing.T) {
	ThatString(t, "hello world").HasPrefix("hello")

	mt := newMockT()
	ThatString(mt, "hello").HasPrefix("world")
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_HasSuffix(t *testing.T) {
	ThatString(t, "hello world").HasSuffix("world")

	mt := newMockT()
	ThatString(mt, "hello").HasSuffix("world")
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_IsEmpty(t *testing.T) {
	ThatString(t, "").IsEmpty()

	mt := newMockT()
	ThatString(mt, "not empty").IsEmpty()
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_HasLen(t *testing.T) {
	ThatString(t, "hello").HasLen(5)

	mt := newMockT()
	ThatString(mt, "hello").HasLen(3)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatString_Matches(t *testing.T) {
	ThatString(t, "hello123").Matches(`^hello\d+$`)

	mt := newMockT()
	ThatString(mt, "hello").Matches(`^\d+$`)
	if !mt.failed {
		t.Error("expected failure")
	}
}

// --- SliceAssertion ---

func TestThatSlice_HasLen(t *testing.T) {
	ThatSlice(t, []int{1, 2, 3}).HasLen(3)

	mt := newMockT()
	ThatSlice(mt, []int{1, 2}).HasLen(5)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatSlice_IsEmpty(t *testing.T) {
	ThatSlice(t, []int{}).IsEmpty()

	mt := newMockT()
	ThatSlice(mt, []int{1}).IsEmpty()
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatSlice_IsNotEmpty(t *testing.T) {
	ThatSlice(t, []int{1}).IsNotEmpty()

	mt := newMockT()
	ThatSlice(mt, []int{}).IsNotEmpty()
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatSlice_Contains(t *testing.T) {
	ThatSlice(t, []string{"a", "b", "c"}).Contains("b")

	mt := newMockT()
	ThatSlice(mt, []string{"a", "b"}).Contains("z")
	if !mt.failed {
		t.Error("expected failure")
	}
}

// --- ErrorAssertion ---

func TestThatError_IsNil(t *testing.T) {
	ThatError(t, nil).IsNil()

	mt := newMockT()
	ThatError(mt, errors.New("oops")).IsNil()
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatError_IsNotNil(t *testing.T) {
	ThatError(t, errors.New("oops")).IsNotNil()

	mt := newMockT()
	ThatError(mt, nil).IsNotNil()
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatError_Is(t *testing.T) {
	sentinel := errors.New("not found")
	wrapped := fmt.Errorf("wrap: %w", sentinel)
	ThatError(t, wrapped).Is(sentinel)

	mt := newMockT()
	ThatError(mt, errors.New("other")).Is(sentinel)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatError_Contains(t *testing.T) {
	ThatError(t, errors.New("file not found")).Contains("not found")

	mt := newMockT()
	ThatError(mt, errors.New("timeout")).Contains("not found")
	if !mt.failed {
		t.Error("expected failure")
	}
}

type pathError struct {
	Path string
}

func (e *pathError) Error() string { return "path error: " + e.Path }

func TestThatError_As(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &pathError{Path: "/tmp"})
	var target *pathError
	ThatError(t, err).As(&target)
	if target == nil || target.Path != "/tmp" {
		t.Errorf("expected target to be set, got %v", target)
	}

	mt := newMockT()
	var target2 *pathError
	ThatError(mt, errors.New("plain")).As(&target2)
	if !mt.failed {
		t.Error("expected failure")
	}
}

func TestThatError_Contains_NilError(t *testing.T) {
	mt := newMockT()
	ThatError(mt, nil).Contains("anything")
	if !mt.failed {
		t.Error("expected failure for nil error")
	}
}

// --- WithMessage ---

func TestThat_WithMessage(t *testing.T) {
	That(t, 42).WithMessage("answer check").Equals(42)

	mt := newMockT()
	That(mt, 1).WithMessage("custom prefix").Equals(2)
	if !mt.failed {
		t.Error("expected failure")
	}
	if len(mt.messages) == 0 {
		t.Fatal("expected error message")
	}
	if !findSubstring(mt.messages[0], "custom prefix") {
		t.Errorf("expected custom prefix in message, got: %s", mt.messages[0])
	}
}

func TestThat_WithMessage_Chaining(t *testing.T) {
	mt := newMockT()
	That(mt, "hello").WithMessage("value check").NotEquals("hello")
	if !mt.failed {
		t.Error("expected failure")
	}
	if len(mt.messages) == 0 {
		t.Fatal("expected error message")
	}
	if !findSubstring(mt.messages[0], "value check") {
		t.Errorf("expected custom prefix in message, got: %s", mt.messages[0])
	}
}

// --- Within (NumericAssertion) ---

func TestThatNumeric_Within(t *testing.T) {
	ThatNumeric(t, 10).Within(10, 0)
	ThatNumeric(t, 10).Within(12, 3)
	ThatNumeric(t, 10).Within(8, 3)
	ThatNumeric(t, 3.14).Within(3.0, 0.2)

	mt := newMockT()
	ThatNumeric(mt, 10).Within(20, 3)
	if !mt.failed {
		t.Error("expected failure when value is outside tolerance")
	}
}

func TestThatNumeric_Within_Boundary(t *testing.T) {
	ThatNumeric(t, 10).Within(7, 3)
	ThatNumeric(t, 10).Within(13, 3)

	mt := newMockT()
	ThatNumeric(mt, 10).Within(6, 3)
	if !mt.failed {
		t.Error("expected failure when value is just outside tolerance")
	}
}

// --- Panics / NotPanics ---

func TestPanics(t *testing.T) {
	Panics(t, func() {
		panic("boom")
	})

	mt := newMockT()
	Panics(mt, func() {
		// does not panic
	})
	if !mt.failed {
		t.Error("expected failure when function does not panic")
	}
}

func TestNotPanics(t *testing.T) {
	NotPanics(t, func() {
		// does not panic
	})

	mt := newMockT()
	NotPanics(mt, func() {
		panic("boom")
	})
	if !mt.failed {
		t.Error("expected failure when function panics")
	}
}

func TestPanics_WithNilPanic(t *testing.T) {
	Panics(t, func() {
		panic(nil)
	})
}

// --- Chaining ---

func TestChaining(t *testing.T) {
	ThatString(t, "hello world").
		Contains("hello").
		Contains("world").
		HasPrefix("hello").
		HasSuffix("world").
		HasLen(11)

	ThatOrdered(t, 10).
		IsGreaterThan(5).
		IsLessThan(20).
		IsGreaterOrEqual(10).
		IsLessOrEqual(10)

	ThatSlice(t, []int{1, 2, 3}).
		HasLen(3).
		IsNotEmpty().
		Contains(2)
}
