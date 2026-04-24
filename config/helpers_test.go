package config

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestIntOrString_IsNumber(t *testing.T) {
	vi := NewIntOrStringForInt(42)
	if !vi.IsNumber() {
		t.Error("expected IsNumber true for int value")
	}
	vs := NewIntOrStringForString("Auto")
	if vs.IsNumber() {
		t.Error("expected IsNumber false for string value")
	}
}

func TestIntOrString_IsStringAndEquals(t *testing.T) {
	v := NewIntOrStringForString("Auto")
	if !v.IsStringAndEquals("Auto") {
		t.Error("expected match for 'Auto'")
	}
	if v.IsStringAndEquals("Other") {
		t.Error("expected no match for 'Other'")
	}
	vi := NewIntOrStringForInt(42)
	if vi.IsStringAndEquals("42") {
		t.Error("expected no match when value is int")
	}
}

func TestIntOrString_UnmarshalYAML(t *testing.T) {
	type wrapper struct {
		Val IntOrString `yaml:"val"`
	}

	t.Run("int", func(t *testing.T) {
		var w wrapper
		if err := yaml.Unmarshal([]byte("val: 42"), &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !w.Val.IsNumber() || *w.Val.IntVal != 42 {
			t.Errorf("expected int 42, got %+v", w.Val)
		}
	})

	t.Run("string", func(t *testing.T) {
		var w wrapper
		if err := yaml.Unmarshal([]byte("val: Auto"), &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Val.IsNumber() || !w.Val.IsStringAndEquals("Auto") {
			t.Errorf("expected string 'Auto', got %+v", w.Val)
		}
	})

	t.Run("invalid_type", func(t *testing.T) {
		var w wrapper
		err := yaml.Unmarshal([]byte("val:\n  - a\n  - b\n"), &w)
		if err == nil {
			t.Error("expected error for sequence node")
		}
	})
}

func TestFloatOrString_IsNumber(t *testing.T) {
	f := 1.5
	v := FloatOrString{FloatVal: &f}
	if !v.IsNumber() {
		t.Error("expected IsNumber true")
	}
	vs := NewFloatOrStringForString("Disabled")
	if vs.IsNumber() {
		t.Error("expected IsNumber false for string")
	}
}

func TestFloatOrString_IsStringAndEquals(t *testing.T) {
	v := NewFloatOrStringForString("Disabled")
	if !v.IsStringAndEquals("Disabled") {
		t.Error("expected match")
	}
	if v.IsStringAndEquals("Other") {
		t.Error("expected no match")
	}
}

func TestFloatOrString_UnmarshalYAML(t *testing.T) {
	type wrapper struct {
		Val FloatOrString `yaml:"val"`
	}

	t.Run("float", func(t *testing.T) {
		var w wrapper
		if err := yaml.Unmarshal([]byte("val: 1.5"), &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !w.Val.IsNumber() || *w.Val.FloatVal != 1.5 {
			t.Errorf("expected float 1.5, got %+v", w.Val)
		}
	})

	t.Run("string", func(t *testing.T) {
		var w wrapper
		if err := yaml.Unmarshal([]byte("val: Disabled"), &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Val.IsNumber() || !w.Val.IsStringAndEquals("Disabled") {
			t.Errorf("expected string 'Disabled', got %+v", w.Val)
		}
	})

	t.Run("invalid_type", func(t *testing.T) {
		var w wrapper
		err := yaml.Unmarshal([]byte("val:\n  - a\n"), &w)
		if err == nil {
			t.Error("expected error for sequence node")
		}
	})
}
