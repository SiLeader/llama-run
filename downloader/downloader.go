package downloader

import (
	"context"
	"fmt"

	"github.com/sileader/llama-run/downloader/huggingface"
	"github.com/sileader/llama-run/downloader/s3"
)

type Type int

const (
	TypeS3 Type = iota
	TypeHuggingFace
)

type Config struct {
	S3 *s3.Config `yaml:"s3"`
}

type Downloader interface {
	Download(ctx context.Context, destinationFile string, model string) error
}

type Builder struct {
	S3 *s3.Config
}

func NewBuilder(config Config) Builder {
	return Builder{
		S3: config.S3,
	}
}

func (b *Builder) Create(dlType Type) (Downloader, error) {
	switch dlType {
	case TypeS3:
		if b.S3 != nil {
			return s3.NewFromConfig(*b.S3)
		}

		return nil, fmt.Errorf("s3 downloader not configured")
	case TypeHuggingFace:
		return huggingface.NewDownloader(), nil
	}
	return nil, fmt.Errorf("unknown downloader type: %v", dlType)
}
