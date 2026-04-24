package huggingface

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
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
	return c.downloadWithVerify(ctx, repo, info.RFilename, destPath, info.SHA256)
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
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
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

	h := sha256.New()
	// ファイルへの書き込みとハッシュ計算を同時に行う
	writer := io.MultiWriter(f, h)
	if _, err := io.Copy(writer, resp.Body); err != nil {
		return err
	}
	slog.DebugContext(ctx, "Download complete")

	sum := hex.EncodeToString(h.Sum(nil))
	if expectedSHA256 != "" && sum != expectedSHA256 {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA256, sum)
	}

	slog.DebugContext(ctx, "Moving to final location", "dest", destPath)
	isOpen = false
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, destPath)
}
