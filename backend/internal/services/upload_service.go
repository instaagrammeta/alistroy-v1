package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type UploadService struct {
	dir        string
	publicBase string
	maxBytes   int64
}

func NewUploadService(dir, publicBase string, maxMB int64) (*UploadService, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &UploadService{
		dir:        dir,
		publicBase: strings.TrimRight(publicBase, "/"),
		maxBytes:   maxMB * 1024 * 1024,
	}, nil
}

var allowedExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true, ".ico": true,
}

// SaveFile validates and stores an uploaded file. Returns the public URL.
func (s *UploadService) SaveFile(file *multipart.FileHeader, subdir string) (string, error) {
	if file.Size <= 0 {
		return "", errors.New("empty file")
	}
	if file.Size > s.maxBytes {
		return "", fmt.Errorf("file too large (max %d MB)", s.maxBytes/(1024*1024))
	}
	ext := strings.ToLower(path.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", errors.New("unsupported file type")
	}
	subdir = sanitizeSubdir(subdir)

	dst := filepath.Join(s.dir, subdir)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return "", err
	}
	name, err := uniqueName(ext)
	if err != nil {
		return "", err
	}
	full := filepath.Join(dst, name)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	out, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	url := s.publicBase + "/" + path.Join(subdir, name)
	return url, nil
}

func sanitizeSubdir(s string) string {
	s = strings.Trim(s, "/")
	s = strings.ReplaceAll(s, "..", "")
	if s == "" {
		s = "misc"
	}
	return s
}

func uniqueName(ext string) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(b), ext), nil
}
