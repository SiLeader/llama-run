package config

import (
	"strings"
	"testing"
)

func TestReasoningConfig_Visit_Modes(t *testing.T) {
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
			cfg := defaultReasoningConfig()
			cfg.Mode = tc.mode
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--reasoning", tc.want) {
				t.Errorf("expected --reasoning %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestReasoningConfig_Visit_InvalidMode(t *testing.T) {
	cfg := defaultReasoningConfig()
	cfg.Mode = "Maybe"
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestReasoningConfig_Visit_Formats(t *testing.T) {
	cases := []struct {
		format string
		want   string
	}{
		{"None", "none"},
		{"Deepseek", "deepseek"},
		{"DeepseekLegacy", "deepseek-legacy"},
	}
	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			cfg := defaultReasoningConfig()
			cfg.Format = tc.format
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--reasoning-format", tc.want) {
				t.Errorf("expected --reasoning-format %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestReasoningConfig_Visit_InvalidFormat(t *testing.T) {
	cfg := defaultReasoningConfig()
	cfg.Format = "Custom"
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestReasoningConfig_Visit_Budgets(t *testing.T) {
	cases := []struct {
		name   string
		budget IntOrString
		want   string
	}{
		{"Unrestricted", NewIntOrStringForString("Unrestricted"), "-1"},
		{"Immediate", NewIntOrStringForString("Immediate"), "0"},
		{"numeric", NewIntOrStringForInt(2048), "2048"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultReasoningConfig()
			cfg.Budget = tc.budget
			b := newMockBuilder()
			if err := cfg.Visit(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSequence(b.args, "--reasoning-budget", tc.want) {
				t.Errorf("expected --reasoning-budget %s, got %v", tc.want, b.args)
			}
		})
	}
}

func TestReasoningConfig_Visit_InvalidBudget(t *testing.T) {
	cfg := defaultReasoningConfig()
	cfg.Budget = NewIntOrStringForString("Unlimited")
	b := newMockBuilder()
	if err := cfg.Visit(b); err == nil {
		t.Error("expected error for invalid budget")
	}
}

func TestReasoningConfig_Visit_BudgetMessage(t *testing.T) {
	cfg := defaultReasoningConfig()
	cfg.BudgetMessage = "Think carefully"
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--reasoning-budget-message", "Think carefully") {
		t.Errorf("expected --reasoning-budget-message, got %v", b.args)
	}
}

func TestReasoningConfig_Visit_EmptyBudgetMessage(t *testing.T) {
	cfg := defaultReasoningConfig()
	cfg.BudgetMessage = ""
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsArg(b.args, "--reasoning-budget-message") {
		t.Error("unexpected --reasoning-budget-message for empty message")
	}
}

func TestChatConfig_Visit_Template(t *testing.T) {
	tmpl := "chatml"
	cfg := &ChatConfig{Template: &tmpl}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--chat-template", "chatml") {
		t.Errorf("expected --chat-template chatml, got %v", b.args)
	}
}

func TestChatConfig_Visit_TemplateFile(t *testing.T) {
	f := "/etc/llama/tmpl.jinja"
	cfg := &ChatConfig{TemplateFile: &f}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSequence(b.args, "--chat-template-file", "/etc/llama/tmpl.jinja") {
		t.Errorf("expected --chat-template-file, got %v", b.args)
	}
}

func TestChatConfig_Visit_TemplateArguments(t *testing.T) {
	cfg := &ChatConfig{
		TemplateArguments: map[string]string{"enable_thinking": "true"},
	}
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsArg(b.args, "--chat-template-kwargs") {
		t.Errorf("expected --chat-template-kwargs, got %v", b.args)
	}
	// verify JSON contains the key
	for i, a := range b.args {
		if a == "--chat-template-kwargs" && i+1 < len(b.args) {
			if !strings.Contains(b.args[i+1], "enable_thinking") {
				t.Errorf("expected JSON with enable_thinking, got %s", b.args[i+1])
			}
			break
		}
	}
}

func TestChatConfig_Visit_NoArguments(t *testing.T) {
	cfg := defaultChatConfig()
	b := newMockBuilder()
	if err := cfg.Visit(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsArg(b.args, "--chat-template") {
		t.Error("unexpected --chat-template for nil template")
	}
	if containsArg(b.args, "--chat-template-kwargs") {
		t.Error("unexpected --chat-template-kwargs for empty arguments")
	}
}
