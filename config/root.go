package config

import (
	"github.com/sileader/llama-run/builder"
	"github.com/sileader/llama-run/downloader"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	LlamaServer builder.LlamaServerConfig `yaml:"llamaServer"`
	Downloader  downloader.Config         `yaml:"downloader"`
	Server      ServerConfig              `yaml:"server"`
	Features    FeaturesConfig            `yaml:"features"`
	Log         LogConfig                 `yaml:"log"`
	Chat        ChatConfig                `yaml:"chat"`
	Reasoning   ReasoningConfig           `yaml:"reasoning"`
	Device      DeviceConfig              `yaml:"device"`
	Model       ModelConfig               `yaml:"model"`
	Sampling    SamplingConfig            `yaml:"sampling"`
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

func (c *Config) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	return visitAll(builder, &c.Server, &c.Features, &c.Log, &c.Chat, &c.Reasoning, &c.Device, &c.Model, &c.Sampling)
}
