package builder

import (
	"context"

	"github.com/sileader/llama-run/downloader"
)

type ApplicationBuilder interface {
	AddArguments(arguments ...string)
	AddEnvironmentVariable(name, value string)
	Go(func(ctx context.Context) error)
	GetDownloader(dlType downloader.Type) (downloader.Downloader, error)
	GetModelDirectory() string
	GetConfigDirectory() string
}
