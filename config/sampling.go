package config

type SamplingConfig struct {
	Samplers         []string    `yaml:"samplers,omitempty"`
	Seed             IntOrString `yaml:"seed,omitempty"`
	Temperature      *float64    `yaml:"temperature,omitempty"`
	TopK             *int        `yaml:"topK,omitempty"`
	TopP             *float64    `yaml:"topP,omitempty"`
	MinP             *float64    `yaml:"minP,omitempty"`
	RepeatLastN      IntOrString `yaml:"repeatLastN,omitempty"`
	RepeatPenalty    *float64    `yaml:"repeatPenalty,omitempty"`
	FrequencyPenalty *float64    `yaml:"frequencyPenalty,omitempty"`
}

func defaultSamplingConfig() SamplingConfig {
	return SamplingConfig{
		Samplers:         []string{"penalties", "dry", "top_n_sigma", "top_k", "typ_p", "top_p", "min_p", "xtc", "temperature"},
		Seed:             NewIntOrStringForString("Random"),
		Temperature:      nil,
		TopK:             nil,
		TopP:             nil,
		MinP:             nil,
		RepeatLastN:      NewIntOrStringForInt(64),
		RepeatPenalty:    nil,
		FrequencyPenalty: nil,
	}
}
