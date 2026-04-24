package config

import (
	"fmt"

	"github.com/sileader/llama-run/builder"
)

type LogConfig struct {
	Enabled   bool    `yaml:"enabled,omitempty"`
	File      *string `yaml:"file"`
	Level     string  `yaml:"level"`
	Timestamp bool    `yaml:"timestamp"`
	ColorMode string  `yaml:"color"`
}

func defaultLogConfig() LogConfig {
	return LogConfig{
		Enabled:   true,
		File:      nil,
		Level:     "Info",
		Timestamp: true,
		ColorMode: "Auto",
	}
}

func (c *LogConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Enabled {
		if c.File != nil {
			builder.AddArguments("--log-file", *c.File)
		}

		builder.AddArguments("--log-verbosity")
		switch c.Level {
		case "Debug":
			builder.AddArguments("4")
		case "Info":
			builder.AddArguments("3")
		case "Warn", "Warning":
			builder.AddArguments("2")
		case "Error":
			builder.AddArguments("1")
		case "Generic":
			builder.AddArguments("0")
		default:
			return fmt.Errorf("unknown log level: %s", c.Level)
		}

		if c.Timestamp {
			builder.AddArguments("--log-timestamps")
		}

		builder.AddArguments("--log-colors")
		switch c.ColorMode {
		case "Auto":
			builder.AddArguments("auto")
		case "On":
			builder.AddArguments("on")
		case "Off":
			builder.AddArguments("off")
		default:
			return fmt.Errorf("unknown color mode: %s", c.ColorMode)
		}
	} else {
		builder.AddArguments("--log-disable")
	}

	return nil
}
