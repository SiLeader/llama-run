package config

import (
	"fmt"
	"strings"

	"github.com/sileader/llama-run/builder"
)

type ServerConfig struct {
	Host         string     `yaml:"host"`
	Port         int        `yaml:"port"`
	ReusePort    bool       `yaml:"reusePort"`
	ApiPrefix    *string    `yaml:"apiPrefix"`
	StaticPath   *string    `yaml:"staticPath"`
	UnsafeApiKey []string   `yaml:"unsafeApiKey,omitempty"`
	ApiKeyFile   *string    `yaml:"apiKeyFile,omitempty"`
	Tls          *TlsConfig `yaml:"tls,omitempty"`
}

type TlsConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func defaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:         "0.0.0.0",
		Port:         8080,
		ReusePort:    false,
		ApiPrefix:    nil,
		StaticPath:   nil,
		UnsafeApiKey: nil,
		ApiKeyFile:   nil,
		Tls:          nil,
	}
}

func (c *ServerConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}
	builder.AddArguments("--host", c.Host)
	builder.AddArguments("--port", fmt.Sprintf("%d", c.Port))
	if c.ReusePort {
		builder.AddArguments("--reuse-port")
	}
	if c.ApiPrefix != nil {
		builder.AddArguments("--api-prefix", *c.ApiPrefix)
	}
	if len(c.UnsafeApiKey) > 0 {
		apiKey := strings.Join(c.UnsafeApiKey, ",")
		builder.AddArguments("--api-key", apiKey)
	}
	if c.ApiKeyFile != nil {
		builder.AddArguments("--api-key-file", *c.ApiKeyFile)
	}
	if c.StaticPath != nil {
		builder.AddArguments("--path", *c.StaticPath)
	}

	return nil
}

func (c *TlsConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}
	builder.AddArguments("--ssl-key-file", c.KeyFile)
	builder.AddArguments("--ssl-cert-file", c.CertFile)

	return nil
}
