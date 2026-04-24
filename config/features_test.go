package config

import (
	"strings"
	"testing"
)

func TestEmbeddingConfig_Visit_Disabled(t *testing.T) {
	cfg := &EmbeddingConfig{Enabled: false}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsArg(b.args, "--embeddings") {
		t.Error("unexpected --embeddings for disabled config")
	}
}

func TestEmbeddingConfig_Visit_Enabled(t *testing.T) {
	cfg := &EmbeddingConfig{Enabled: true}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--embeddings") {
		t.Errorf("expected --embeddings, got %v", b.args)
	}
}

func TestEmbeddingConfig_Visit_Pooling(t *testing.T) {
	cases := []struct {
		pooling string
		want    string
	}{
		{"None", "none"},
		{"Mean", "mean"},
		{"Cls", "cls"},
		{"Last", "last"},
		{"Rank", "rank"},
	}
	for _, tc := range cases {
		t.Run(tc.pooling, func(t *testing.T) {
			p := tc.pooling
			cfg := &EmbeddingConfig{Enabled: true, Pooling: &p}
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--pooling", tc.want) {
				t.Errorf("expected --pooling %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestEmbeddingConfig_Visit_InvalidPooling(t *testing.T) {
	p := "Invalid"
	cfg := &EmbeddingConfig{Enabled: true, Pooling: &p}
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid pooling")
	}
}

func TestRerankConfig_Visit(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		b := newMockBuilder()
		if err := (&RerankConfig{Enabled: false}).Visit(b); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if containsArg(b.args, "--reranking") {
			t.Error("unexpected --reranking")
		}
	})
	t.Run("enabled", func(t *testing.T) {
		b := newMockBuilder()
		if err := (&RerankConfig{Enabled: true}).Visit(b); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !containsArg(b.args, "--reranking") {
			t.Errorf("expected --reranking, got %v", b.args)
		}
	})
}

func TestWebuiConfig_Visit_Disabled(t *testing.T) {
	b := newMockBuilder()
	if err := (&WebuiConfig{Enabled: false}).Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--no-webui") {
		t.Errorf("expected --no-webui, got %v", b.args)
	}
	if containsArg(b.args, "--webui") {
		t.Error("unexpected --webui for disabled config")
	}
}

func TestWebuiConfig_Visit_Enabled(t *testing.T) {
	b := newMockBuilder()
	if err := (&WebuiConfig{Enabled: true}).Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--webui") {
		t.Errorf("expected --webui, got %v", b.args)
	}
}

func TestWebuiConfig_Visit_ConfigFile(t *testing.T) {
	f := "/etc/webui.json"
	cfg := &WebuiConfig{Enabled: true, ConfigFile: &f}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--webui-config-file", "/etc/webui.json") {
		t.Errorf("expected --webui-config-file, got %v", b.args)
	}
}

func TestWebuiConfig_Visit_Config(t *testing.T) {
	v := any(map[string]string{"theme": "dark"})
	cfg := &WebuiConfig{Enabled: true, Config: &v}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--webui-config") {
		t.Errorf("expected --webui-config, got %v", b.args)
	}
	// verify it's a JSON string containing the key
	for i, a := range b.args {
		if a == "--webui-config" && i+1 < len(b.args) {
			if !strings.Contains(b.args[i+1], "theme") {
				t.Errorf("expected JSON with theme key, got %s", b.args[i+1])
			}
			break
		}
	}
}

func TestMetricsConfig_Visit(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&MetricsConfig{Enabled: false}).Visit(b)
		if containsArg(b.args, "--metrics") {
			t.Error("unexpected --metrics")
		}
	})
	t.Run("enabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&MetricsConfig{Enabled: true}).Visit(b)
		if !containsArg(b.args, "--metrics") {
			t.Errorf("expected --metrics, got %v", b.args)
		}
	})
}

func TestPropertiesConfig_Visit(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&PropertiesConfig{Enabled: false}).Visit(b)
		if containsArg(b.args, "--props") {
			t.Error("unexpected --props")
		}
	})
	t.Run("enabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&PropertiesConfig{Enabled: true}).Visit(b)
		if !containsArg(b.args, "--props") {
			t.Errorf("expected --props, got %v", b.args)
		}
	})
}

func TestJinjaConfig_Visit(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&JinjaConfig{Enabled: true}).Visit(b)
		if !containsArg(b.args, "--jinja") {
			t.Errorf("expected --jinja, got %v", b.args)
		}
	})
	t.Run("disabled", func(t *testing.T) {
		b := newMockBuilder()
		_ = (&JinjaConfig{Enabled: false}).Visit(b)
		if !containsArg(b.args, "--no-jinja") {
			t.Errorf("expected --no-jinja, got %v", b.args)
		}
	})
}

func TestFeaturesConfig_Visit_Default(t *testing.T) {
	cfg := defaultFeaturesConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// default: jinja enabled, webui disabled
	if !containsArg(b.args, "--jinja") {
		t.Errorf("expected --jinja from defaults, got %v", b.args)
	}
	if !containsArg(b.args, "--no-webui") {
		t.Errorf("expected --no-webui from defaults, got %v", b.args)
	}
}
