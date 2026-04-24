package config

import (
	"fmt"
	"strings"

	"github.com/sileader/llama-run/builder"
)

type SamplingConfig struct {
	Samplers         []string      `yaml:"samplers,omitempty"`
	Seed             IntOrString   `yaml:"seed,omitempty"`
	Temperature      *float64      `yaml:"temperature,omitempty"`
	TopK             *int          `yaml:"topK,omitempty"`
	TopP             *float64      `yaml:"topP,omitempty"`
	MinP             *float64      `yaml:"minP,omitempty"`
	RepeatLastN      IntOrString   `yaml:"repeatLastN,omitempty"`
	RepeatPenalty    FloatOrString `yaml:"repeatPenalty,omitempty"`
	PresencePenalty  FloatOrString `yaml:"presencePenalty,omitempty"`
	FrequencyPenalty FloatOrString `yaml:"frequencyPenalty,omitempty"`
}

func defaultSamplingConfig() SamplingConfig {
	return SamplingConfig{
		Samplers:         nil,
		Seed:             NewIntOrStringForString("Random"),
		Temperature:      nil,
		TopK:             nil,
		TopP:             nil,
		MinP:             nil,
		RepeatLastN:      NewIntOrStringForInt(64),
		RepeatPenalty:    NewFloatOrStringForString("Disabled"),
		PresencePenalty:  NewFloatOrStringForString("Disabled"),
		FrequencyPenalty: NewFloatOrStringForString("Disabled"),
	}
}

func (c *SamplingConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}
	if len(c.Samplers) > 0 {
		samplers := strings.Join(c.Samplers, ";")
		builder.AddArguments("--samplers", samplers)
	}
	if c.Seed.IsNumber() {
		builder.AddArguments("--seed", fmt.Sprintf("%d", *c.Seed.IntVal))
	} else if !c.Seed.IsStringAndEquals("Random") {
		return fmt.Errorf("sampling config: seed must be number or 'Random'")
	}

	if c.Temperature != nil {
		builder.AddArguments("--temperature", fmt.Sprintf("%0.2f", *c.Temperature))
	}
	if c.TopK != nil {
		builder.AddArguments("--top-k", fmt.Sprintf("%d", *c.TopK))
	}
	if c.TopP != nil {
		builder.AddArguments("--top-p", fmt.Sprintf("%0.2f", *c.TopP))
	}
	if c.MinP != nil {
		builder.AddArguments("--min-p", fmt.Sprintf("%0.2f", *c.MinP))
	}

	if c.RepeatLastN.IsNumber() {
		builder.AddArguments("--repeat-last-n", fmt.Sprintf("%d", *c.RepeatLastN.IntVal))
	} else if c.RepeatLastN.IsStringAndEquals("Disabled") {
		builder.AddArguments("--repeat-last-n", "0")
	} else if c.RepeatLastN.IsStringAndEquals("Context") {
		builder.AddArguments("--repeat-last-n", "-1")
	} else {
		return fmt.Errorf("sampling config: repeatLastN must be number, 'Disabled', or 'Context'")
	}

	if c.RepeatPenalty.IsNumber() {
		builder.AddArguments("--repeat-penalty", fmt.Sprintf("%0.2f", *c.RepeatPenalty.FloatVal))
	} else if !c.RepeatPenalty.IsStringAndEquals("Disabled") {
		return fmt.Errorf("sampling config: repeatPenalty must be number or 'Disabled'")
	}

	if c.PresencePenalty.IsNumber() {
		builder.AddArguments("--presence-penalty", fmt.Sprintf("%0.2f", *c.PresencePenalty.FloatVal))
	} else if !c.PresencePenalty.IsStringAndEquals("Disabled") {
		return fmt.Errorf("sampling config: presencePenalty must be number or 'Disabled'")
	}

	if c.FrequencyPenalty.IsNumber() {
		builder.AddArguments("--frequency-penalty", fmt.Sprintf("%0.2f", *c.FrequencyPenalty.FloatVal))
	} else if !c.FrequencyPenalty.IsStringAndEquals("Disabled") {
		return fmt.Errorf("sampling config: frequencyPenalty must be number or 'Disabled'")
	}

	return nil
}
