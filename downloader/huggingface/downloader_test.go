package huggingface

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sileader/llama-run/downloader/checksum"
)

func TestParseModel_Valid(t *testing.T) {
	repo, quantize, err := parseModel("org/repo:Q4_K_M")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo != "org/repo" {
		t.Errorf("expected repo org/repo, got %s", repo)
	}
	if quantize != "Q4_K_M" {
		t.Errorf("expected quantize Q4_K_M, got %s", quantize)
	}
}

func TestParseModel_MissingColon(t *testing.T) {
	if _, _, err := parseModel("org/repo"); err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParseModel_MissingSlash(t *testing.T) {
	if _, _, err := parseModel("repo:Q4_K_M"); err == nil {
		t.Error("expected error for missing slash in repo name")
	}
}

func TestParseModel_TooManyColons(t *testing.T) {
	if _, _, err := parseModel("org/repo:q4:extra"); err == nil {
		t.Error("expected error for too many colons")
	}
}

func TestDownloadWithVerify_Success(t *testing.T) {
	content := []byte("fake model content")
	hash := sha256.Sum256(content)
	expectedSHA := hex.EncodeToString(hash[:])

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	dir := t.TempDir()
	dest := filepath.Join(dir, "model.gguf")

	if err := dl.downloadWithVerify(context.Background(), "org/repo", "model.gguf", dest, expectedSHA); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestDownloadWithVerify_ChecksumMismatch(t *testing.T) {
	content := []byte("fake model content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	dir := t.TempDir()
	dest := filepath.Join(dir, "model.gguf")

	err := dl.downloadWithVerify(context.Background(), "org/repo", "model.gguf", dest, "deadbeef")
	if err == nil {
		t.Error("expected checksum mismatch error")
	}

	// temp file should be cleaned up
	if _, statErr := os.Stat(dest + ".llamarunpartialdownload"); statErr == nil {
		t.Error("expected temp file to be cleaned up")
	}
}

func TestDownloadWithVerify_NoChecksum(t *testing.T) {
	content := []byte("model data")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	sha256sumBytes := sha256.Sum256(content)
	sha256sum := hex.EncodeToString(sha256sumBytes[:])

	dl := newDownloaderWithBaseURL("", srv.URL)
	dir := t.TempDir()
	dest := filepath.Join(dir, "model.gguf")

	// empty expectedSHA -> skip checksum verification
	if err := dl.downloadWithVerify(context.Background(), "org/repo", "model.gguf", dest, sha256sum); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("dest file not created: %v", err)
	}
}

func TestDownloadWithVerify_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	dir := t.TempDir()
	dest := filepath.Join(dir, "model.gguf")

	if err := dl.downloadWithVerify(context.Background(), "org/repo", "model.gguf", dest, ""); err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestFetchFileInfo_Found(t *testing.T) {
	info := modelInfo{
		Siblings: []fileInfo{
			{RFilename: "model-Q4_K_M.gguf", Lfs: fileLfsInfo{SHA256: "abc123"}},
			{RFilename: "model-Q8_0.gguf", Lfs: fileLfsInfo{SHA256: "def456"}},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(info)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	fi, err := dl.fetchFileInfo(context.Background(), "org/repo", "Q4_K_M")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.RFilename != "model-Q4_K_M.gguf" {
		t.Errorf("expected model-Q4_K_M.gguf, got %s", fi.RFilename)
	}
	if fi.Lfs.SHA256 != "abc123" {
		t.Errorf("expected abc123, got %s", fi.Lfs.SHA256)
	}
}

func TestFetchFileInfo_NotFound(t *testing.T) {
	info := modelInfo{
		Siblings: []fileInfo{
			{RFilename: "model-Q8_0.gguf", Lfs: fileLfsInfo{SHA256: "def456"}},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(info)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	if _, err := dl.fetchFileInfo(context.Background(), "org/repo", "Q4_K_M"); err == nil {
		t.Error("expected error when quantization not found")
	}
}

func TestFetchFileInfo_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("", srv.URL)
	if _, err := dl.fetchFileInfo(context.Background(), "org/repo", "Q4_K_M"); err == nil {
		t.Error("expected error for non-200 response")
	}
}

func TestFetchFileInfo_BearerToken(t *testing.T) {
	var receivedAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		info := modelInfo{Siblings: []fileInfo{{RFilename: "model-q4.gguf"}}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(info)
	}))
	defer srv.Close()

	dl := newDownloaderWithBaseURL("mytoken", srv.URL)
	_, _ = dl.fetchFileInfo(context.Background(), "org/repo", "q4")
	if receivedAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %q", receivedAuth)
	}
}

func TestChecksumFile(t *testing.T) {
	content := []byte("hello world")
	f, err := os.CreateTemp(t.TempDir(), "checksum-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	f.Close()

	sum, err := checksum.ChecksumFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := sha256.Sum256(content)
	want := hex.EncodeToString(expected[:])
	if sum != want {
		t.Errorf("expected %s, got %s", want, sum)
	}
}
