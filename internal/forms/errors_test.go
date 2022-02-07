package forms

import "testing"

func TestErrors_Add(t *testing.T) {
	e := errors{}
	e.Add("a", "test")

	if e["a"][0] != "test" {
		t.Error("expected to find `a` in errors, but not found")
	}
}

func TestErrors_Get(t *testing.T) {
	e := errors{}
	e["a"] = append(e["a"], "test")

	if e.Get("a") != "test" {
		t.Error("expected to find `a` in errors but not found")
	}

	if e.Get("b") != "" {
		t.Error("expected to not find `b` in errors but found somehow")
	}
}
