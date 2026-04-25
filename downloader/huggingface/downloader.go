package huggingface

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sileader/llama-run/downloader/checksum"
)

type Downloader struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewDownloader() *Downloader {
	token := os.Getenv("HF_TOKEN")
	return newDownloaderWithBaseURL(token, "https://huggingface.co")
}

func newDownloaderWithBaseURL(token, baseURL string) *Downloader {
	return &Downloader{
		token:   token,
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
	}
}

func (c *Downloader) Download(ctx context.Context, destPath string, model string) error {
	repo, quantize, err := parseModel(model)
	if err != nil {
		return err
	}
	info, err := c.fetchFileInfo(ctx, repo, quantize)
	if err != nil {
		return err
	}
	return c.downloadWithVerify(ctx, repo, info.RFilename, destPath, info.Lfs.SHA256)
}

func parseModel(model string) (string, string, error) {
	spl := strings.Split(model, ":")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("invalid model: %s", model)
	}
	s2 := strings.Split(spl[0], "/")
	if len(s2) != 2 {
		return "", "", fmt.Errorf("invalid model: %s", model)
	}
	return spl[0], spl[1], nil
}

func (c *Downloader) downloadWithVerify(ctx context.Context, repo, filename, destPath, expectedSHA256 string) error {
	slog.InfoContext(ctx, "Downloading model from Hugging Face", "repo", repo, "filename", filename, "dest", destPath)
	if len(expectedSHA256) == 0 {
		slog.WarnContext(ctx, "No expected SHA256 provided, skipping verification", "repo", repo, "filename", filename)
	}
	// キャッシュ済みの場合もチェックサムを検証
	if len(expectedSHA256) > 0 {
		if _, err := os.Stat(destPath); err == nil {
			slog.DebugContext(ctx, "Cache hit", "repo", repo, "filename", filename)
			sum, err := checksum.ChecksumFile(destPath)
			if err != nil {
				return err
			}
			if sum == expectedSHA256 {
				slog.InfoContext(ctx, "Checksum match, skipping download", "repo", repo, "filename", filename)
				return nil // 正常
			}
			// 壊れているので再ダウンロード
			slog.InfoContext(ctx, "Checksum mismatch, redownloading", "repo", repo, "filename", filename)
			os.Remove(destPath)
		}
	}

	url := fmt.Sprintf("%s/%s/resolve/main/%s", c.baseURL, repo, filename)
	slog.DebugContext(ctx, "Downloading", "url", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	tmp := destPath + ".llamarunpartialdownload"
	slog.DebugContext(ctx, "Writing to", "tmp", tmp)
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	isOpen := true
	defer func() {
		if isOpen {
			f.Close()
		}
		os.Remove(tmp) // エラー時の掃除（Rename後は空振り）
	}()

	writer := checksum.NewSha256FileWriter(f)
	// ファイルへの書き込みとハッシュ計算を同時に行う
	if _, err := io.Copy(writer, resp.Body); err != nil {
		return err
	}
	slog.DebugContext(ctx, "Download complete")

	if len(expectedSHA256) > 0 {
		if err := writer.CheckDigest(expectedSHA256); err != nil {
			return err
		}
	}

	slog.DebugContext(ctx, "Moving to final location", "dest", destPath)
	isOpen = false
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, destPath)
}
