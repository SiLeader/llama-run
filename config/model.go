package config

type ModelConfig struct {
	Docker      *string           `yaml:"docker"`
	HuggingFace *string           `yaml:"huggingFace"`
	Local       LocalModelsConfig `yaml:"local"`
}

type LocalModelsConfig struct {
	Models []AliasModelConfig `yaml:"models"`
}

type AliasModelConfig struct {
	Alias string `yaml:"alias"`
}

func defaultModelConfig() ModelConfig {
	return ModelConfig{
		Docker:      nil,
		HuggingFace: nil,
		Local: LocalModelsConfig{
			Models: []AliasModelConfig{},
		},
	}
}
