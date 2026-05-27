package storage

import (
	"io"
)

type FileStorage interface {
	Upload(key string, reader io.Reader, contentType string) error
	Download(key string) (io.ReadCloser, error)
	Delete(key string) error
}
