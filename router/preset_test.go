package router

import (
	"strings"
	"testing"
)

func TestConfig_AddModel_Valid(t *testing.T) {
	c := &Config{Models: map[string]Info{}}
	info := Info{Context: intPtr(4096)}
	if err := c.AddModel("my-model", info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := c.Models["my-model"]; !ok {
		t.Error("model not added")
	}
}

func TestConfig_AddModel_InvalidAlias(t *testing.T) {
	cases := []string{
		"invalid alias",
		"alias with spaces",
		"alias/slash",
		"alias.dot",
		"",
	}
	for _, alias := range cases {
		t.Run(alias, func(t *testing.T) {
			c := &Config{Models: map[string]Info{}}
			if err := c.AddModel(alias, Info{}); err == nil {
				t.Errorf("expected error for alias %q", alias)
			}
		})
	}
}

func TestConfig_AddModel_ValidAliasChars(t *testing.T) {
	validAliases := []string{
		"model",
		"model-v2",
		"model_v2",
		"Model123",
		"a",
	}
	for _, alias := range validAliases {
		t.Run(alias, func(t *testing.T) {
			c := &Config{Models: map[string]Info{}}
			if err := c.AddModel(alias, Info{}); err != nil {
				t.Errorf("unexpected error for alias %q: %v", alias, err)
			}
		})
	}
}

func TestConfig_AddModel_NilMap(t *testing.T) {
	c := &Config{}
	if err := c.AddModel("model", Info{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(c.Models))
	}
}

func TestInfo_String_Empty(t *testing.T) {
	info := Info{}
	if info.String() != "" {
		t.Errorf("expected empty string, got %q", info.String())
	}
}

func TestInfo_String_Fields(t *testing.T) {
	tmpl := "chatml"
	gpuLayers := 32
	jinja := true
	ctx := 4096
	model := "/models/llama.gguf"
	info := Info{
		ChatTemplate: &tmpl,
		GpuLayers:    &gpuLayers,
		Jinja:        &jinja,
		Context:      &ctx,
		Model:        &model,
	}
	s := info.String()
	if !strings.Contains(s, "chat-template = chatml") {
		t.Errorf("expected chat-template in output, got %q", s)
	}
	if !strings.Contains(s, "gpu-layers = 32") {
		t.Errorf("expected gpu-layers in output, got %q", s)
	}
	if !strings.Contains(s, "jinja = true") {
		t.Errorf("expected jinja in output, got %q", s)
	}
	if !strings.Contains(s, "c = 4096") {
		t.Errorf("expected c in output, got %q", s)
	}
	if !strings.Contains(s, "model = /models/llama.gguf") {
		t.Errorf("expected model in output, got %q", s)
	}
}

func TestConfig_String_Version(t *testing.T) {
	c := Config{Models: map[string]Info{}}
	s := c.String()
	if !strings.HasPrefix(s, "version = 1\n") {
		t.Errorf("expected version header, got %q", s)
	}
}

func TestConfig_String_Default(t *testing.T) {
	ctx := 8192
	c := Config{
		Default: &Info{Context: &ctx},
		Models:  map[string]Info{},
	}
	s := c.String()
	if !strings.Contains(s, "[*]") {
		t.Errorf("expected [*] section for default, got %q", s)
	}
	if !strings.Contains(s, "c = 8192") {
		t.Errorf("expected c = 8192, got %q", s)
	}
}

func TestConfig_String_Models(t *testing.T) {
	model := "/models/test.gguf"
	c := Config{
		Models: map[string]Info{
			"test": {Model: &model},
		},
	}
	s := c.String()
	if !strings.Contains(s, "[test]") {
		t.Errorf("expected [test] section, got %q", s)
	}
	if !strings.Contains(s, "model = /models/test.gguf") {
		t.Errorf("expected model path, got %q", s)
	}
}

func intPtr(v int) *int {
	return &v
}
