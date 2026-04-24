package config

import "testing"

func TestUnmarshalConfig_Empty(t *testing.T) {
	cfg, err := UnmarshalConfig([]byte("{}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// verify defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if !cfg.Features.Jinja.Enabled {
		t.Error("expected jinja enabled by default")
	}
	if cfg.Log.Level != "Info" {
		t.Errorf("expected default log level Info, got %s", cfg.Log.Level)
	}
	if cfg.Log.ColorMode != "Auto" {
		t.Errorf("expected default color Auto, got %s", cfg.Log.ColorMode)
	}
	if cfg.Reasoning.Mode != "Auto" {
		t.Errorf("expected default reasoning mode Auto, got %s", cfg.Reasoning.Mode)
	}
	if cfg.Reasoning.Format != "None" {
		t.Errorf("expected default reasoning format None, got %s", cfg.Reasoning.Format)
	}
	if !cfg.Device.Memory.Mmap {
		t.Error("expected default mmap=true")
	}
}

func TestUnmarshalConfig_Override(t *testing.T) {
	yaml := `
server:
  host: 127.0.0.1
  port: 9090
log:
  level: Debug
`
	cfg, err := UnmarshalConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected 9090, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "Debug" {
		t.Errorf("expected Debug, got %s", cfg.Log.Level)
	}
	// non-overridden defaults still apply
	if cfg.Log.ColorMode != "Auto" {
		t.Errorf("expected default color Auto, got %s", cfg.Log.ColorMode)
	}
}

func TestUnmarshalConfig_InvalidYAML(t *testing.T) {
	_, err := UnmarshalConfig([]byte("server: [invalid"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestConfig_Visit_NilSafe(t *testing.T) {
	var cfg *Config
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error for nil config: %v", err)
	}
}

func TestConfig_Visit_ProducesArgs(t *testing.T) {
	cfg, err := UnmarshalConfig([]byte("{}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error from Visit: %v", err)
	}
	// should have at least --host from server config
	if !containsArg(b.args, "--host") {
		t.Errorf("expected --host in args, got %v", b.args)
	}
	// should have --log-verbosity from log config
	if !containsArg(b.args, "--log-verbosity") {
		t.Errorf("expected --log-verbosity in args, got %v", b.args)
	}
}

func TestUnmarshalConfig_ApiKeys(t *testing.T) {
	yaml := `
server:
  apiKey:
    - key1
    - key2
`
	cfg, err := UnmarshalConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Server.UnsafeApiKey) != 2 {
		t.Errorf("expected 2 api keys, got %d", len(cfg.Server.UnsafeApiKey))
	}
}

func TestUnmarshalConfig_TLS(t *testing.T) {
	yaml := `
server:
  tls:
    certFile: /certs/cert.pem
    keyFile: /certs/key.pem
`
	cfg, err := UnmarshalConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Tls == nil {
		t.Fatal("expected TLS config")
	}
	if cfg.Server.Tls.CertFile != "/certs/cert.pem" {
		t.Errorf("unexpected cert file: %s", cfg.Server.Tls.CertFile)
	}
}
