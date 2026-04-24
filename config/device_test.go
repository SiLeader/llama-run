package config

import (
	"fmt"
	"testing"
)

func TestCpuConfig_Visit_NumericThreads(t *testing.T) {
	cfg := &CpuConfig{Threads: NewIntOrStringForInt(8)}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--threads", "8") {
		t.Errorf("expected --threads 8, got %v", b.args)
	}
}

func TestCpuConfig_Visit_AutoThreads(t *testing.T) {
	cfg := &CpuConfig{Threads: NewIntOrStringForString("Auto")}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Auto resolves to runtime.NumCPU() (or cgroup limit); just verify the flag is present
	if !containsArg(b.args, "--threads") {
		t.Errorf("expected --threads flag, got %v", b.args)
	}
}

func TestCpuConfig_Visit_UnknownString(t *testing.T) {
	// unknown string -> flag is silently omitted (no error)
	cfg := &CpuConfig{Threads: NewIntOrStringForString("Many")}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsArg(b.args, "--threads") {
		t.Errorf("expected no --threads for unknown string, got %v", b.args)
	}
}

func TestMemoryConfig_Visit_Mmap(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		b := newMockBuilder()
		if err := (&MemoryConfig{Mmap: true}).Visit(b); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !containsArg(b.args, "--mmap") {
			t.Errorf("expected --mmap, got %v", b.args)
		}
		if containsArg(b.args, "--no-mmap") {
			t.Error("unexpected --no-mmap")
		}
	})
	t.Run("disabled", func(t *testing.T) {
		b := newMockBuilder()
		if err := (&MemoryConfig{Mmap: false}).Visit(b); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !containsArg(b.args, "--no-mmap") {
			t.Errorf("expected --no-mmap, got %v", b.args)
		}
	})
}

func TestGpuConfig_Visit_Layers(t *testing.T) {
	cases := []struct {
		name    string
		layers  IntOrString
		wantArg string
	}{
		{"number", NewIntOrStringForInt(32), "32"},
		{"Auto", NewIntOrStringForString("Auto"), "auto"},
		{"All", NewIntOrStringForString("All"), "all"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &GpuConfig{Layers: tc.layers, MainIndex: 0}
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--n-gpu-layers", tc.wantArg) {
				t.Errorf("expected --n-gpu-layers %s, got %v", tc.wantArg, b.args)
			}
		})
	}
}

func TestGpuConfig_Visit_MainIndex(t *testing.T) {
	cfg := &GpuConfig{Layers: NewIntOrStringForString("Auto"), MainIndex: 2}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--main-gpu", fmt.Sprintf("%d", 2)) {
		t.Errorf("expected --main-gpu 2, got %v", b.args)
	}
}

func TestDeviceConfig_Visit_Default(t *testing.T) {
	cfg := defaultDeviceConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// default memory: mmap=true
	if !containsArg(b.args, "--mmap") {
		t.Errorf("expected --mmap from defaults, got %v", b.args)
	}
	// default gpu layers: Auto
	if !containsSequence(b.args, "--n-gpu-layers", "auto") {
		t.Errorf("expected --n-gpu-layers auto from defaults, got %v", b.args)
	}
	// default main-gpu: 0
	if !containsSequence(b.args, "--main-gpu", "0") {
		t.Errorf("expected --main-gpu 0 from defaults, got %v", b.args)
	}
}
