package testassert

import (
	"testing"
)

// --- MapAssertion ---

func TestThatMap_HasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	ThatMap(t, m).HasKey("a")
	ThatMap(t, m).HasKey("b")

	mt := newMockT()
	ThatMap(mt, m).HasKey("z")
	if !mt.failed {
		t.Error("expected failure for missing key")
	}
}

func TestThatMap_HasLen(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	ThatMap(t, m).HasLen(3)

	mt := newMockT()
	ThatMap(mt, m).HasLen(5)
	if !mt.failed {
		t.Error("expected failure for wrong length")
	}
}

func TestThatMap_IsEmpty(t *testing.T) {
	ThatMap(t, map[string]int{}).IsEmpty()

	mt := newMockT()
	ThatMap(mt, map[string]int{"a": 1}).IsEmpty()
	if !mt.failed {
		t.Error("expected failure for non-empty map")
	}
}

func TestThatMap_IsNotEmpty(t *testing.T) {
	ThatMap(t, map[string]int{"a": 1}).IsNotEmpty()

	mt := newMockT()
	ThatMap(mt, map[string]int{}).IsNotEmpty()
	if !mt.failed {
		t.Error("expected failure for empty map")
	}
}

func TestThatMap_Chaining(t *testing.T) {
	m := map[string]int{"x": 10, "y": 20}
	ThatMap(t, m).
		HasLen(2).
		IsNotEmpty().
		HasKey("x").
		HasKey("y")
}

func TestThatMap_IntKeys(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	ThatMap(t, m).HasKey(1).HasLen(2).IsNotEmpty()

	mt := newMockT()
	ThatMap(mt, m).HasKey(99)
	if !mt.failed {
		t.Error("expected failure for missing int key")
	}
}
