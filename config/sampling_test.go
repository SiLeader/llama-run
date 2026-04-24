package config

import "testing"

func TestSamplingConfig_Visit_Default(t *testing.T) {
	cfg := defaultSamplingConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// seed "Random" -> no --seed flag
	if containsArg(b.args, "--seed") {
		t.Error("unexpected --seed for Random seed")
	}
	// repeatLastN 64 -> --repeat-last-n 64
	if !containsSequence(b.args, "--repeat-last-n", "64") {
		t.Errorf("expected --repeat-last-n 64, got %v", b.args)
	}
	// penalties "Disabled" -> no flags
	if containsArg(b.args, "--repeat-penalty") {
		t.Error("unexpected --repeat-penalty for Disabled")
	}
}

func TestSamplingConfig_Visit_NumericSeed(t *testing.T) {
	cfg := defaultSamplingConfig()
	cfg.Seed = NewIntOrStringForInt(12345)
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--seed", "12345") {
		t.Errorf("expected --seed 12345, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_InvalidSeed(t *testing.T) {
	cfg := defaultSamplingConfig()
	cfg.Seed = NewIntOrStringForString("Fixed")
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid seed value")
	}
}

func TestSamplingConfig_Visit_Temperature(t *testing.T) {
	cfg := defaultSamplingConfig()
	temp := 0.7
	cfg.Temperature = &temp
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--temperature", "0.70") {
		t.Errorf("expected --temperature 0.70, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_TopK(t *testing.T) {
	cfg := defaultSamplingConfig()
	k := 40
	cfg.TopK = &k
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--top-k", "40") {
		t.Errorf("expected --top-k 40, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_TopP(t *testing.T) {
	cfg := defaultSamplingConfig()
	p := 0.95
	cfg.TopP = &p
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--top-p", "0.95") {
		t.Errorf("expected --top-p 0.95, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_MinP(t *testing.T) {
	cfg := defaultSamplingConfig()
	p := 0.05
	cfg.MinP = &p
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--min-p", "0.05") {
		t.Errorf("expected --min-p 0.05, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_Samplers(t *testing.T) {
	cfg := defaultSamplingConfig()
	cfg.Samplers = []string{"top_k", "top_p", "temp"}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--samplers", "top_k;top_p;temp") {
		t.Errorf("expected --samplers top_k;top_p;temp, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_RepeatLastN(t *testing.T) {
	cases := []struct {
		name  string
		val   IntOrString
		want  string
		isErr bool
	}{
		{"number", NewIntOrStringForInt(128), "128", false},
		{"Disabled", NewIntOrStringForString("Disabled"), "0", false},
		{"Context", NewIntOrStringForString("Context"), "-1", false},
		{"invalid", NewIntOrStringForString("Invalid"), "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultSamplingConfig()
			cfg.RepeatLastN = tc.val
			b := newMockBuilder()
			err := cfg.Visit(b)
			if tc.isErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--repeat-last-n", tc.want) {
				t.Errorf("expected --repeat-last-n %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestSamplingConfig_Visit_RepeatPenalty(t *testing.T) {
	cases := []struct {
		name  string
		val   FloatOrString
		want  string
		isErr bool
	}{
		{"number", FloatOrString{FloatVal: func() *float64 { v := 1.1; return &v }()}, "1.10", false},
		{"Disabled", NewFloatOrStringForString("Disabled"), "", false},
		{"invalid", NewFloatOrStringForString("Invalid"), "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultSamplingConfig()
			cfg.RepeatPenalty = tc.val
			b := newMockBuilder()
			err := cfg.Visit(b)
			if tc.isErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.want != "" && !containsSequence(b.args, "--repeat-penalty", tc.want) {
				t.Errorf("expected --repeat-penalty %s, got %v", tc.want, b.args)
			}
			if tc.want == "" && containsArg(b.args, "--repeat-penalty") {
				t.Error("unexpected --repeat-penalty for Disabled")
			}
		})
	}
}

func TestSamplingConfig_Visit_PresencePenalty(t *testing.T) {
	cfg := defaultSamplingConfig()
	v := 0.5
	cfg.PresencePenalty = FloatOrString{FloatVal: &v}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--presence-penalty", "0.50") {
		t.Errorf("expected --presence-penalty 0.50, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_InvalidPresencePenalty(t *testing.T) {
	cfg := defaultSamplingConfig()
	cfg.PresencePenalty = NewFloatOrStringForString("Invalid")
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid presencePenalty")
	}
}

func TestSamplingConfig_Visit_FrequencyPenalty(t *testing.T) {
	cfg := defaultSamplingConfig()
	v := 0.3
	cfg.FrequencyPenalty = FloatOrString{FloatVal: &v}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--frequency-penalty", "0.30") {
		t.Errorf("expected --frequency-penalty 0.30, got %v", b.args)
	}
}

func TestSamplingConfig_Visit_InvalidFrequencyPenalty(t *testing.T) {
	cfg := defaultSamplingConfig()
	cfg.FrequencyPenalty = NewFloatOrStringForString("Bad")
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid frequencyPenalty")
	}
}
