package utils

import (
	"io"
	"os"
	"path/filepath"
)

// EnsureDir creates a directory and all parent directories if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks whether a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CopyFile copies a file from src to dst. If dst already exists, it will be
// overwritten.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	// BUG: srcFile is never closed — this leaks a file descriptor.
	// The defer srcFile.Close() is missing here.

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// ReadFileString reads the entire contents of a file as a string.
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ListFiles returns all files (not directories) in the given directory.
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}
