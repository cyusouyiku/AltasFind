package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrPatternRequired = errors.New("search pattern is required")
	ErrRootNotFound    = errors.New("root path does not exist or is not a directory")
	ErrInvalidWorkers  = errors.New("workers must be greater than zero")
)

// FileReadError wraps a file path with its underlying read error.
type FileReadError struct {
	Path string
	Err  error
}

func (e *FileReadError) Error() string {
	return fmt.Sprintf("read %s: %v", e.Path, e.Err)
}

func (e *FileReadError) Unwrap() error {
	return e.Err
}
