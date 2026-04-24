package router

import (
	"fmt"
	"regexp"
)

type Config struct {
	Default *Info
	Models  map[string]Info
}

type Info struct {
	ChatTemplate *string
	GpuLayers    *int
	Jinja        *bool
	Context      *int
	Model        *string
}

func (c Config) String() string {
	s := "version = 1\n"

	if c.Default != nil {
		s += fmt.Sprintf("[*]\n%s\n", c.Default)
	}
	for alias, model := range c.Models {
		s += fmt.Sprintf("[%s]\n%s\n", alias, model)
	}
	return s
}

func (i Info) String() string {
	s := ""
	if i.ChatTemplate != nil {
		s += fmt.Sprintf("chat-template = %s\n", *i.ChatTemplate)
	}
	if i.GpuLayers != nil {
		s += fmt.Sprintf("gpu-layers = %d\n", *i.GpuLayers)
	}
	if i.Jinja != nil {
		s += fmt.Sprintf("jinja = %t\n", *i.Jinja)
	}
	if i.Context != nil {
		s += fmt.Sprintf("c = %d\n", *i.Context)
	}
	if i.Model != nil {
		s += fmt.Sprintf("model = %s\n", *i.Model)
	}
	return s
}

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateAlias(alias string) error {
	if !aliasRegex.MatchString(alias) {
		return fmt.Errorf("invalid alias: %s", alias)
	}
	return nil
}

func (c *Config) AddModel(alias string, model Info) error {
	if err := validateAlias(alias); err != nil {
		return err
	}

	if c.Models == nil {
		c.Models = map[string]Info{}
	}
	c.Models[alias] = model
	return nil
}
