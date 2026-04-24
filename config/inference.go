package config

import (
	"encoding/json"
	"fmt"

	"github.com/sileader/llama-run/builder"
)

type ReasoningConfig struct {
	/// Reasoning mode. "Auto", "On", "Off"
	Mode string `yaml:"mode"`

	/// Reasoning format. "None", "Deepseek", "DeepseekLegacy"
	Format string `yaml:"format"`

	/// Reasoning budget. "Unrestricted", "Immediate", or a number of tokens
	Budget IntOrString `yaml:"budget"`

	BudgetMessage string `yaml:"budgetMessage"`
}

type ChatConfig struct {
	Template          *string           `yaml:"template"`
	TemplateFile      *string           `yaml:"templateFile"`
	TemplateArguments map[string]string `yaml:"templateArguments"`
}

func defaultReasoningConfig() ReasoningConfig {
	return ReasoningConfig{
		Mode:          "Auto",
		Format:        "None",
		Budget:        NewIntOrStringForString("Unrestricted"),
		BudgetMessage: "",
	}
}

func defaultChatConfig() ChatConfig {
	return ChatConfig{
		Template:     nil,
		TemplateFile: nil,
	}
}

func (c *ReasoningConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	builder.AddArguments("--reasoning")
	switch c.Mode {
	case "Auto":
		builder.AddArguments("auto")
	case "On":
		builder.AddArguments("on")
	case "Off":
		builder.AddArguments("off")
	default:
		return fmt.Errorf("unknown reasoning mode: %s (allows 'Auto', 'On', or 'Off')", c.Mode)
	}

	builder.AddArguments("--reasoning-format")
	switch c.Format {
	case "None":
		builder.AddArguments("none")
	case "Deepseek":
		builder.AddArguments("deepseek")
	case "DeepseekLegacy":
		builder.AddArguments("deepseek-legacy")
	default:
		return fmt.Errorf("unknown reasoning format: %s (allows 'None', 'Deepseek', or 'DeepseekLegacy')", c.Format)
	}

	if c.Budget.IsNumber() {
		builder.AddArguments("--reasoning-budget", fmt.Sprintf("%d", *c.Budget.IntVal))
	} else if c.Budget.IsStringAndEquals("Unrestricted") {
		builder.AddArguments("--reasoning-budget", "-1")
	} else if c.Budget.IsStringAndEquals("Immediate") {
		builder.AddArguments("--reasoning-budget", "0")
	} else {
		return fmt.Errorf("unknown reasoning budget: %s (allows number, 'Unrestricted', or 'Immediate')", *c.Budget.StrVal)
	}

	if len(c.BudgetMessage) > 0 {
		builder.AddArguments("--reasoning-budget-message", c.BudgetMessage)
	}

	return nil
}

func (c *ChatConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Template != nil {
		builder.AddArguments("--chat-template", *c.Template)
	}

	if c.TemplateFile != nil {
		builder.AddArguments("--chat-template-file", *c.TemplateFile)
	}

	if len(c.TemplateArguments) > 0 {
		bs, err := json.Marshal(c.TemplateArguments)
		if err != nil {
			return err
		}
		builder.AddArguments("--chat-template-kwargs", string(bs))
	}

	return nil
}
