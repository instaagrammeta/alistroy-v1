package services

import (
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
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".svg": true, ".ico": true, ".mp4": true, ".webm": true, ".mov": true,
}

// SaveFile validates and stores an uploaded file. Returns the public URL + size.
func (s *UploadService) SaveFile(file *multipart.FileHeader, subdir string) (string, int64, string, error) {
	if file.Size <= 0 {
		return "", 0, "", errors.New("empty file")
	}
	if file.Size > s.maxBytes {
		return "", 0, "", fmt.Errorf("file too large (max %d MB)", s.maxBytes/(1024*1024))
	}
	ext := strings.ToLower(path.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", 0, "", errors.New("unsupported file type")
	}
	subdir = sanitizeSubdir(subdir)
	dst := filepath.Join(s.dir, subdir)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return "", 0, "", err
	}
	name, err := uniqueName(ext)
	if err != nil {
		return "", 0, "", err
	}
	full := filepath.Join(dst, name)

	src, err := file.Open()
	if err != nil {
		return "", 0, "", err
	}
	defer src.Close()
	out, err := os.Create(full)
	if err != nil {
		return "", 0, "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return "", 0, "", err
	}
	url := s.publicBase + "/" + path.Join(subdir, name)
	return url, file.Size, file.Header.Get("Content-Type"), nil
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
	h, err := randomHex(8)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), h, ext), nil
}
