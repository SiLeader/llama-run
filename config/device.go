package config

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sileader/llama-run/builder"
)

type DeviceConfig struct {
	Cpu    CpuConfig    `yaml:"cpu"`
	Memory MemoryConfig `yaml:"memory"`
	Gpu    GpuConfig    `yaml:"gpu"`
}

type CpuConfig struct {
	Threads IntOrString `yaml:"threads"`
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

func (c *DeviceConfig) Visit(builder builder.ApplicationBuilder) error {
	if err := c.Cpu.Visit(builder); err != nil {
		return err
	}
	if err := c.Memory.Visit(builder); err != nil {
		return err
	}
	if err := c.Gpu.Visit(builder); err != nil {
		return err
	}
	return nil
}

func (c *CpuConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Threads.IsNumber() {
		builder.AddArguments("--threads", fmt.Sprintf("%d", *c.Threads.IntVal))
	} else if c.Threads.IsStringAndEquals("Auto") {
		threads := getCpuThreads()
		builder.AddArguments("--threads", fmt.Sprintf("%d", threads))
	}

	return nil
}

func (c *MemoryConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Mmap {
		builder.AddArguments("--mmap")
	} else {
		builder.AddArguments("--no-mmap")
	}
	return nil
}

func (c *GpuConfig) Visit(builder builder.ApplicationBuilder) error {
	if c == nil {
		return nil
	}

	if c.Layers.IsNumber() {
		builder.AddArguments("--n-gpu-layers", fmt.Sprintf("%d", *c.Layers.IntVal))
	} else if c.Layers.IsStringAndEquals("Auto") {
		builder.AddArguments("--n-gpu-layers", "auto")
	} else if c.Layers.IsStringAndEquals("All") {
		builder.AddArguments("--n-gpu-layers", "all")
	}

	builder.AddArguments("--main-gpu", fmt.Sprintf("%d", c.MainIndex))

	return nil
}

func getCpuThreads() (threads int) {
	threads = runtime.NumCPU()

	file, err := os.Open("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if sl := strings.Split(line, " "); len(sl) == 2 {
			if sl[0] == "max" {
				continue
			}
			allowedUs, err := strconv.ParseInt(sl[0], 10, 64)
			if err != nil {
				slog.Warn("failed to parse CPU allowed micro secs", "error", err)
				return
			}
			unitUs, err := strconv.ParseInt(sl[1], 10, 64)
			if err != nil {
				slog.Warn("failed to parse CPU unit micro secs", "error", err)
				return
			}

			threads = int(max(allowedUs/unitUs, 1))

			return
		}
	}
	return
}
