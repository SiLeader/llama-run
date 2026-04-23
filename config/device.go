package config

type DeviceConfig struct {
	Cpu    CpuConfig    `yaml:"cpu"`
	Memory MemoryConfig `yaml:"memory"`
	Gpu    GpuConfig    `yaml:"gpu"`
}

type CpuConfig struct {
	Threads IntOrString `yaml:"threads"`
	Batch   *int        `yaml:"batch"`
}

type MemoryConfig struct {
	Mmap bool `yaml:"mmap"`
}

type GpuConfig struct {
	Layers    IntOrString `yaml:"layers"`
	MainIndex int         `yaml:"mainIndex"`
}

func defaultDeviceConfig() DeviceConfig {
	return DeviceConfig{
		Cpu: CpuConfig{
			Threads: NewIntOrStringForString("Auto"),
			Batch:   nil,
		},
		Memory: MemoryConfig{
			Mmap: true,
		},
		Gpu: GpuConfig{
			Layers:    NewIntOrStringForString("Auto"),
			MainIndex: 0,
		},
	}
}
