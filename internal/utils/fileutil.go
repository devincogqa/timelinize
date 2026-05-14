package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks whether a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir creates a directory (and parents) if it doesn't already exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExtension returns the lowercase file extension without the leading dot.
func FileExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext[1:])
}

// HumanReadableSize converts a byte count into a human-readable string.
func HumanReadableSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ListFiles returns all file paths (non-recursive) in the given directory.
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}
