package config

import (
	"testing"
)

func TestServerConfig_Visit_Defaults(t *testing.T) {
	cfg := defaultServerConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsSequence(b.args, "--host", "0.0.0.0") {
		t.Errorf("expected --host 0.0.0.0, got %v", b.args)
	}
	if !containsSequence(b.args, "--port", "8080") {
		t.Errorf("expected --port 8080, got %v", b.args)
	}
	if containsArg(b.args, "--reuse-port") {
		t.Error("unexpected --reuse-port")
	}
	if containsArg(b.args, "--api-prefix") {
		t.Error("unexpected --api-prefix")
	}
	if containsArg(b.args, "--api-key") {
		t.Error("unexpected --api-key")
	}
}

func TestServerConfig_Visit_ReusePort(t *testing.T) {
	cfg := defaultServerConfig()
	cfg.ReusePort = true
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--reuse-port") {
		t.Errorf("expected --reuse-port, got %v", b.args)
	}
}

func TestServerConfig_Visit_ApiPrefix(t *testing.T) {
	cfg := defaultServerConfig()
	p := "/v1"
	cfg.ApiPrefix = &p
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--api-prefix", "/v1") {
		t.Errorf("expected --api-prefix /v1, got %v", b.args)
	}
}

func TestServerConfig_Visit_ApiKey_Single(t *testing.T) {
	cfg := defaultServerConfig()
	cfg.UnsafeApiKey = []string{"secret"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--api-key", "secret") {
		t.Errorf("expected --api-key secret, got %v", b.args)
	}
}

func TestServerConfig_Visit_ApiKey_Multiple(t *testing.T) {
	cfg := defaultServerConfig()
	cfg.UnsafeApiKey = []string{"key1", "key2"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--api-key", "key1,key2") {
		t.Errorf("expected --api-key key1,key2, got %v", b.args)
	}
}

func TestServerConfig_Visit_ApiKeyFile(t *testing.T) {
	cfg := defaultServerConfig()
	f := "/etc/llama/keys"
	cfg.ApiKeyFile = &f
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--api-key-file", "/etc/llama/keys") {
		t.Errorf("expected --api-key-file, got %v", b.args)
	}
}

func TestServerConfig_Visit_StaticPath(t *testing.T) {
	cfg := defaultServerConfig()
	p := "/var/www"
	cfg.StaticPath = &p
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--path", "/var/www") {
		t.Errorf("expected --path /var/www, got %v", b.args)
	}
}

func TestTlsConfig_Visit(t *testing.T) {
	cfg := &TlsConfig{
		CertFile: "/certs/cert.pem",
		KeyFile:  "/certs/key.pem",
	}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--ssl-cert-file", "/certs/cert.pem") {
		t.Errorf("expected --ssl-cert-file, got %v", b.args)
	}
	if !containsSequence(b.args, "--ssl-key-file", "/certs/key.pem") {
		t.Errorf("expected --ssl-key-file, got %v", b.args)
	}
}

func TestTlsConfig_Visit_Nil(t *testing.T) {
	var cfg *TlsConfig
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.args) != 0 {
		t.Errorf("expected no args for nil TlsConfig, got %v", b.args)
	}
}
