package config

type ServerConfig struct {
	Host       string     `yaml:"host"`
	Port       int        `yaml:"port"`
	ReusePort  bool       `yaml:"reusePort"`
	ApiPrefix  *string    `yaml:"apiPrefix"`
	StaticPath *string    `yaml:"staticPath"`
	ApiKey     []string   `yaml:"apiKey,omitempty"`
	ApiKeyFile *string    `yaml:"apiKeyFile,omitempty"`
	Tls        *TlsConfig `yaml:"tls,omitempty"`
}

type TlsConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func defaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:       "0.0.0.0",
		Port:       8080,
		ReusePort:  false,
		ApiPrefix:  nil,
		StaticPath: nil,
		ApiKey:     nil,
		ApiKeyFile: nil,
		Tls:        nil,
	}
}
