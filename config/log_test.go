package config

import "testing"

func TestLogConfig_Visit_Disabled(t *testing.T) {
	cfg := &LogConfig{Enabled: false}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--log-disable") {
		t.Errorf("expected --log-disable, got %v", b.args)
	}
	if containsArg(b.args, "--log-verbosity") {
		t.Error("unexpected --log-verbosity when log is disabled")
	}
}

func TestLogConfig_Visit_Levels(t *testing.T) {
	cases := []struct {
		level string
		want  string
	}{
		{"Debug", "4"},
		{"Info", "3"},
		{"Warn", "2"},
		{"Warning", "2"},
		{"Error", "1"},
		{"Generic", "0"},
	}
	for _, tc := range cases {
		t.Run(tc.level, func(t *testing.T) {
			cfg := &LogConfig{Enabled: true, Level: tc.level, ColorMode: "Auto"}
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--log-verbosity", tc.want) {
				t.Errorf("expected --log-verbosity %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestLogConfig_Visit_InvalidLevel(t *testing.T) {
	cfg := &LogConfig{Enabled: true, Level: "Verbose", ColorMode: "Auto"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid log level")
	}
}

func TestLogConfig_Visit_Timestamp(t *testing.T) {
	cfg := &LogConfig{Enabled: true, Level: "Info", Timestamp: true, ColorMode: "Auto"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--log-timestamps") {
		t.Errorf("expected --log-timestamps, got %v", b.args)
	}
}

func TestLogConfig_Visit_NoTimestamp(t *testing.T) {
	cfg := &LogConfig{Enabled: true, Level: "Info", Timestamp: false, ColorMode: "Auto"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsArg(b.args, "--log-timestamps") {
		t.Error("unexpected --log-timestamps")
	}
}

func TestLogConfig_Visit_ColorModes(t *testing.T) {
	cases := []struct {
		mode string
		want string
	}{
		{"Auto", "auto"},
		{"On", "on"},
		{"Off", "off"},
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			cfg := &LogConfig{Enabled: true, Level: "Info", ColorMode: tc.mode}
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--log-colors", tc.want) {
				t.Errorf("expected --log-colors %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestLogConfig_Visit_InvalidColorMode(t *testing.T) {
	cfg := &LogConfig{Enabled: true, Level: "Info", ColorMode: "Rainbow"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid color mode")
	}
}

func TestLogConfig_Visit_File(t *testing.T) {
	f := "/var/log/llama.log"
	cfg := &LogConfig{Enabled: true, Level: "Info", File: &f, ColorMode: "Auto"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--log-file", "/var/log/llama.log") {
		t.Errorf("expected --log-file, got %v", b.args)
	}
}

func TestLogConfig_Visit_Default(t *testing.T) {
	cfg := defaultLogConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--log-verbosity", "3") {
		t.Errorf("expected --log-verbosity 3 (Info default), got %v", b.args)
	}
	if !containsSequence(b.args, "--log-colors", "auto") {
		t.Errorf("expected --log-colors auto (Auto default), got %v", b.args)
	}
}
