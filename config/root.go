package config

import (
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Features  FeaturesConfig  `yaml:"features"`
	Log       LogConfig       `yaml:"log"`
	Chat      ChatConfig      `yaml:"chat"`
	Reasoning ReasoningConfig `yaml:"reasoning"`
	Device    DeviceConfig    `yaml:"device"`
	Model     ModelConfig     `yaml:"model"`
	Sampling  SamplingConfig  `yaml:"sampling"`
}

func defaultConfig() Config {
	return Config{
		Server:    defaultServerConfig(),
		Features:  defaultFeaturesConfig(),
		Log:       defaultLogConfig(),
		Chat:      defaultChatConfig(),
		Reasoning: defaultReasoningConfig(),
		Device:    defaultDeviceConfig(),
		Model:     defaultModelConfig(),
		Sampling:  defaultSamplingConfig(),
	}
}

func UnmarshalConfig(data []byte) (*Config, error) {
	config := defaultConfig()
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
