package testassert

import (
	"testing"
)

func TestThatJSON_Equals(t *testing.T) {
	got := `{"name":"alice","age":30}`
	want := `{
		"age": 30,
		"name": "alice"
	}`
	ThatJSON(t, got).Equals(want)

	mt := newMockT()
	ThatJSON(mt, `{"name":"alice"}`).Equals(`{"name":"bob"}`)
	if !mt.failed {
		t.Error("expected failure for different JSON values")
	}
}

func TestThatJSON_Equals_InvalidJSON(t *testing.T) {
	mt := newMockT()
	ThatJSON(mt, `not json`).Equals(`{}`)
	if !mt.failed {
		t.Error("expected failure for invalid JSON")
	}

	mt2 := newMockT()
	ThatJSON(mt2, `{}`).Equals(`not json`)
	if !mt2.failed {
		t.Error("expected failure for invalid want JSON")
	}
}

func TestThatJSON_Contains(t *testing.T) {
	got := `{"name":"alice","age":30,"active":true}`
	ThatJSON(t, got).Contains("name", "alice")
	ThatJSON(t, got).Contains("age", 30)
	ThatJSON(t, got).Contains("active", true)

	mt := newMockT()
	ThatJSON(mt, got).Contains("name", "bob")
	if !mt.failed {
		t.Error("expected failure for wrong value")
	}

	mt2 := newMockT()
	ThatJSON(mt2, got).Contains("missing", "value")
	if !mt2.failed {
		t.Error("expected failure for missing key")
	}
}

func TestThatJSON_HasKey(t *testing.T) {
	got := `{"name":"alice","age":30}`
	ThatJSON(t, got).HasKey("name").HasKey("age")

	mt := newMockT()
	ThatJSON(mt, got).HasKey("missing")
	if !mt.failed {
		t.Error("expected failure for missing key")
	}
}

func TestThatJSON_Chaining(t *testing.T) {
	got := `{"name":"alice","age":30}`
	ThatJSON(t, got).
		HasKey("name").
		HasKey("age").
		Contains("name", "alice").
		Equals(`{"age":30,"name":"alice"}`)
}
