package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

func ChecksumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

type Sha256FileWriter struct {
	hash   hash.Hash
	writer io.Writer
}

func NewSha256FileWriter(writer io.Writer) *Sha256FileWriter {
	h := sha256.New()
	return &Sha256FileWriter{
		hash:   h,
		writer: io.MultiWriter(writer, h),
	}
}

func (w *Sha256FileWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return n, err
	}
	w.hash.Write(p[:n])
	return n, nil
}

func (w *Sha256FileWriter) CheckDigest(sum string) error {
	calculatedSum := hex.EncodeToString(w.hash.Sum(nil))
	if calculatedSum != sum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", sum, calculatedSum)
	}
	return nil
}
