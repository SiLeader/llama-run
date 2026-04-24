package config

import (
	"context"
	"fmt"

	"github.com/sileader/llama-run/downloader"
)

type mockBuilder struct {
	args []string
	envs []string
}

func newMockBuilder() *mockBuilder {
	return &mockBuilder{}
}

func (m *mockBuilder) AddArguments(args ...string) {
	m.args = append(m.args, args...)
}

func (m *mockBuilder) AddEnvironmentVariable(name, value string) {
	m.envs = append(m.envs, fmt.Sprintf("%s=%s", name, value))
}

func (m *mockBuilder) Go(func(ctx context.Context) error) {}

func (m *mockBuilder) GetDownloader(downloader.Type) (downloader.Downloader, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockBuilder) GetModelDirectory() string {
	return "/models"
}

func (m *mockBuilder) GetConfigDirectory() string {
	return "/config"
}

// containsSequence reports whether seq appears as a contiguous subsequence in args.
func containsSequence(args []string, seq ...string) bool {
	if len(seq) == 0 {
		return true
	}
outer:
	for i := 0; i <= len(args)-len(seq); i++ {
		for j, s := range seq {
			if args[i+j] != s {
				continue outer
			}
		}
		return true
	}
	return false
}

func containsArg(args []string, arg string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}
