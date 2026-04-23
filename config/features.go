package config

type FeaturesConfig struct {
	Embedding  EmbeddingConfig
	Rerank     SimpleEnabledConfig
	Webui      WebuiConfig
	Metrics    SimpleEnabledConfig
	Properties SimpleEnabledConfig
	Jinja      SimpleEnabledConfig
}

type EmbeddingConfig struct {
	Enabled bool    `yaml:"enabled"`
	Pooling *string `yaml:"pooling"`
}

type WebuiConfig struct {
	Enabled    bool    `yaml:"enabled"`
	Config     *any    `yaml:"config,omitempty"`
	ConfigFile *string `yaml:"configFile,omitempty"`
}

type SimpleEnabledConfig struct {
	Enabled bool `yaml:"enabled"`
}

func defaultFeaturesConfig() FeaturesConfig {
	return FeaturesConfig{
		Embedding: EmbeddingConfig{
			Enabled: false,
			Pooling: nil,
		},
		Rerank: SimpleEnabledConfig{
			Enabled: false,
		},
		Webui: WebuiConfig{
			Enabled:    false,
			Config:     nil,
			ConfigFile: nil,
		},
		Metrics: SimpleEnabledConfig{
			Enabled: false,
		},
		Properties: SimpleEnabledConfig{
			Enabled: false,
		},
		Jinja: SimpleEnabledConfig{
			Enabled: true,
		},
	}
}
