package config

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
