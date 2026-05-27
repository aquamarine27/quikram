package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, 0755)
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Upload(key string, reader io.Reader, contentType string) error {
	path := filepath.Join(s.basePath, key)
	os.MkdirAll(filepath.Dir(path), 0755)

	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, reader)
	return err
}

func (s *LocalStorage) Download(key string) (io.ReadCloser, error) {
	path := filepath.Join(s.basePath, key)
	return os.Open(path)
}

func (s *LocalStorage) Delete(key string) error {
	path := filepath.Join(s.basePath, key)
	return os.Remove(path)
}
