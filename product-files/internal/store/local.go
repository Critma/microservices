package store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	maxFileSize int64
	basePath    string
}

func NewLocal(basePath string, maxSize int64) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

func (l *Local) Save(path string, content io.Reader) error {
	fp := filepath.Join(l.basePath, path)

	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("unable to get file info: %w", err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	_, err = io.CopyN(f, content, l.maxFileSize)
	if err != nil && err != io.EOF {
		return fmt.Errorf("unable to write to file: %w", err)
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fp := filepath.Join(l.basePath, path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	return f, nil
}
