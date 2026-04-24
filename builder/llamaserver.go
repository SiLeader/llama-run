package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sileader/llama-run/downloader"
	"golang.org/x/sync/errgroup"
)

type LlamaServerApplicationBuilder struct {
	cmd       *exec.Cmd
	eg        *errgroup.Group
	ctx       context.Context
	directory directoryConfig
	dlb       downloader.Builder
}

type LlamaServerConfig struct {
	Executable string          `yaml:"executable"`
	Arguments  []string        `yaml:"arguments"`
	Directory  directoryConfig `yaml:"directory"`
}

type directoryConfig struct {
	Model  string `yaml:"model"`
	Config string `yaml:"config"`
}

func NewLlamaServerApplicationBuilder(ctx context.Context, config LlamaServerConfig, dlb downloader.Builder) *LlamaServerApplicationBuilder {
	cmd := exec.Command(config.Executable, config.Arguments...)
	cmd.Env = os.Environ()

	eg, ctx := errgroup.WithContext(ctx)
	return &LlamaServerApplicationBuilder{
		cmd:       cmd,
		eg:        eg,
		ctx:       ctx,
		directory: config.Directory,
		dlb:       dlb,
	}
}

func (b *LlamaServerApplicationBuilder) Exec() error {
	if err := b.eg.Wait(); err != nil {
		return err
	}
	return syscall.Exec(b.cmd.Path, b.cmd.Args, b.cmd.Env)
}

func (b *LlamaServerApplicationBuilder) AddArguments(args ...string) {
	b.cmd.Args = append(b.cmd.Args, args...)
}

func (b *LlamaServerApplicationBuilder) AddEnvironmentVariable(name, value string) {
	b.cmd.Env = append(b.cmd.Env, fmt.Sprintf("%s=%s", name, value))
}

func (b *LlamaServerApplicationBuilder) Go(task func(ctx context.Context) error) {
	b.eg.Go(func() error {
		return task(b.ctx)
	})
}

func (b *LlamaServerApplicationBuilder) GetDownloader(dlType downloader.Type) (downloader.Downloader, error) {
	return b.dlb.Create(dlType)
}

func (b *LlamaServerApplicationBuilder) GetModelDirectory() string {
	return b.directory.Model
}

func (b *LlamaServerApplicationBuilder) GetConfigDirectory() string {
	return b.directory.Config
}
