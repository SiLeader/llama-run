package config

type ReasoningConfig struct {
	/// Reasoning mode. "Auto", "On", "Off"
	Mode string `yaml:"mode"`

	/// Reasoning format. "Auto", "Deepseek", "DeepseekLegacy"
	Format string `yaml:"format"`

	/// Reasoning budget. "Unrestricted", "Immediate", or a number of tokens
	Budget IntOrString `yaml:"budget"`

	BudgetMessage string `yaml:"budgetMessage"`
}

type ChatConfig struct {
	Template     *string `yaml:"template"`
	TemplateFile *string `yaml:"templateFile"`
}

func defaultReasoningConfig() ReasoningConfig {
	return ReasoningConfig{
		Mode:          "Auto",
		Format:        "Auto",
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
