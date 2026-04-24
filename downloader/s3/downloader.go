package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Downloader struct {
	client *s3.Client
}

type Config struct {
	Region       string  `yaml:"region"`
	Endpoint     *string `yaml:"endpoint"`
	AccessKeyEnv *string `yaml:"accessKeyEnv"`
	SecretKeyEnv *string `yaml:"secretKeyEnv"`
}

func NewFromConfig(config Config) (*Downloader, error) {
	if config.Endpoint == nil {
		return newDownloader(config.Region)
	}
	var accessKey string
	var secretKey string
	if config.AccessKeyEnv != nil {
		accessKey = os.Getenv(*config.AccessKeyEnv)
	} else {
		accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if config.SecretKeyEnv != nil {
		secretKey = os.Getenv(*config.SecretKeyEnv)
	} else {
		secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	return newDownloaderForS3Compatible(*config.Endpoint, config.Region, accessKey, secretKey)
}

func newDownloader(region string) (*Downloader, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	dlr := Downloader{
		client: s3.NewFromConfig(cfg),
	}
	return &dlr, nil
}

func newDownloaderForS3Compatible(
	endpoint,
	region,
	accessKey,
	secretKey string,
) (*Downloader, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	dlr := Downloader{
		client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)

			// S3互換サーバーでは多くの場合 path-style が必要:
			//   http://host/bucket/key
			// AWS S3 の virtual-host style:
			//   https://bucket.s3.region.amazonaws.com/key
			o.UsePathStyle = true
		}),
	}
	return &dlr, nil
}

func parseModel(model string) (string, string, error) {
	u, err := url.Parse(model)
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "s3" {
		return "", "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")
	return bucket, key, nil
}

func (d *Downloader) Download(ctx context.Context, destPath string, model string) error {
	bucket, key, err := parseModel(model)
	if err != nil {
		return err
	}

	result, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket:       &bucket,
		Key:          &key,
		ChecksumMode: types.ChecksumModeEnabled,
	})
	if err != nil {
		return err
	}

	defer result.Body.Close()
	tmp := destPath + ".llamarunpartialdownload"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	isOpen := true
	defer func() {
		if isOpen {
			file.Close()
		}
		os.Remove(tmp)
	}()
	_, err = io.Copy(file, result.Body)
	if err != nil {
		return err
	}
	isOpen = false
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, destPath)
}
