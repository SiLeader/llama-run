package config

import (
	"encoding/json"
	"fmt"

	"github.com/sileader/llama-run/builder"
)

type FeaturesConfig struct {
	Embedding  EmbeddingConfig  `yaml:"embedding"`
	Rerank     RerankConfig     `yaml:"rerank"`
	Webui      WebuiConfig      `yaml:"webui"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Properties PropertiesConfig `yaml:"properties"`
	Jinja      JinjaConfig      `yaml:"jinja"`
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

type RerankConfig struct {
	Enabled bool `yaml:"enabled"`
}

type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
}

type PropertiesConfig struct {
	Enabled bool `yaml:"enabled"`
}

type JinjaConfig struct {
	Enabled bool `yaml:"enabled"`
}

func defaultFeaturesConfig() FeaturesConfig {
	return FeaturesConfig{
		Embedding: EmbeddingConfig{
			Enabled: false,
			Pooling: nil,
		},
		Rerank: RerankConfig{
			Enabled: false,
		},
		Webui: WebuiConfig{
			Enabled:    false,
			Config:     nil,
			ConfigFile: nil,
		},
		Metrics: MetricsConfig{
			Enabled: false,
		},
		Properties: PropertiesConfig{
			Enabled: false,
		},
		Jinja: JinjaConfig{
			Enabled: true,
		},
	}
}

func (c *FeaturesConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	return visitAll(builder, &c.Embedding, &c.Rerank, &c.Webui, &c.Metrics, &c.Properties, &c.Jinja)
}

func (c *EmbeddingConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--embeddings")

		if c.Pooling != nil {
			builder.AddArguments("--pooling")
			switch *c.Pooling {
			case "None":
				builder.AddArguments("none")
			case "Mean":
				builder.AddArguments("mean")
			case "Cls":
				builder.AddArguments("cls")
			case "Last":
				builder.AddArguments("last")
			case "Rank":
				builder.AddArguments("rank")
			default:
				return fmt.Errorf("unknown pooling option '%s' was passed (allows 'None', 'Mean', 'Cls', 'Last', or 'Rank')", *c.Pooling)
			}
		}
	}
	return nil
}

func (c *RerankConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--reranking")
	}
	return nil
}

func (c *WebuiConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--webui")
		if c.Config != nil {
			bs, err := json.Marshal(*c.Config)
			if err != nil {
				return err
			}
			builder.AddArguments("--webui-config", string(bs))
		}
		if c.ConfigFile != nil {
			builder.AddArguments("--webui-config-file", *c.ConfigFile)
		}
	} else {
		builder.AddArguments("--no-webui")
	}
	return nil
}

func (c *MetricsConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--metrics")
	}
	return nil
}

func (c *PropertiesConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--props")
	}
	return nil
}

func (c *JinjaConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		builder.AddArguments("--jinja")
	} else {
		builder.AddArguments("--no-jinja")
	}
	return nil
}
