package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sileader/llama-run/builder"
	"github.com/sileader/llama-run/downloader"
	"github.com/sileader/llama-run/router"
)

type ModelConfig struct {
	Alias       *string             `yaml:"alias"`
	Aliases     []string            `yaml:"aliases"`
	Docker      *string             `yaml:"docker"`
	HuggingFace *string             `yaml:"huggingFace"`
	Router      *RouterModelsConfig `yaml:"router"`
}

type RouterModelsConfig struct {
	Default *RouterModelsDefaultConfig `yaml:"default"`
	Models  []AliasModelConfig         `yaml:"models"`
}

type RouterModelsDefaultConfig struct {
	Context   *int `yaml:"context"`
	GpuLayers *int `yaml:"gpuLayers"`
}

type AliasModelConfig struct {
	Alias       string  `yaml:"alias"`
	Context     *int    `yaml:"context"`
	Path        *string `yaml:"path"`
	HuggingFace *string `yaml:"huggingFace"`
	S3          *string `yaml:"s3"`
}

func defaultModelConfig() ModelConfig {
	return ModelConfig{
		Docker:      nil,
		HuggingFace: nil,
		Router:      nil,
	}
}

func (c *ModelConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Docker != nil {
		builder.AddArguments("--docker-repo", *c.Docker)
	}
	if c.HuggingFace != nil {
		builder.AddArguments("--hf-repo", *c.HuggingFace)
	}

	if c.Alias != nil || len(c.Aliases) > 0 {
		aliases := c.Aliases
		if c.Alias != nil {
			aliases = append(c.Aliases, *c.Alias)
		}
		alias := strings.Join(aliases, ",")
		builder.AddArguments("--alias", alias)
	}

	if c.Router != nil {
		if err := c.Router.Visit(builder); err != nil {
			return err
		}
	}

	return nil
}

func (c *RouterModelsConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}
	preset := router.Config{
		Default: nil,
		Models:  map[string]router.Info{},
	}

	if c.Default != nil {
		preset.Default = &router.Info{
			ChatTemplate: nil,
			GpuLayers:    c.Default.GpuLayers,
			Jinja:        nil,
			Context:      c.Default.Context,
			Model:        nil,
		}
	}
	modelDir := builder.GetModelDirectory()
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return err
	}

	for _, model := range c.Models {
		modelPath := path.Join(modelDir, fmt.Sprintf("%s.gguf", model.Alias))
		info := router.Info{
			ChatTemplate: nil,
			GpuLayers:    nil,
			Jinja:        nil,
			Context:      model.Context,
			Model:        &modelPath,
		}
		if err := preset.AddModel(model.Alias, info); err != nil {
			return err
		}

		builder.Go(func(ctx context.Context) error {
			var dlType downloader.Type
			var m string
			if model.S3 != nil {
				dlType = downloader.TypeS3
				m = *model.S3
			} else if model.HuggingFace != nil {
				dlType = downloader.TypeHuggingFace
				m = *model.HuggingFace
			} else {
				return fmt.Errorf("unknown downloader type: %v", model)
			}
			dlr, err := builder.GetDownloader(dlType)
			if err != nil {
				return err
			}

			return dlr.Download(ctx, modelPath, m)
		})
	}

	presetIni := path.Join(builder.GetConfigDirectory(), "preset.ini")
	builder.Go(func(ctx context.Context) error {
		return os.WriteFile(presetIni, []byte(preset.String()), 0644)
	})

	builder.AddArguments("--models-preset", presetIni)

	return nil
}
