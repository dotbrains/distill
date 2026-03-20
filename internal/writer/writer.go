package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteOutput writes compacted content to the output directory.
func WriteOutput(baseDir, outputDir, filename, content string) (string, error) {
	dir := filepath.Join(baseDir, outputDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("writing output: %w", err)
	}
	return path, nil
}

// WriteIndex writes an index.md file to the given directory.
func WriteIndex(baseDir, outputDir, content string) (string, error) {
	dir := filepath.Join(baseDir, outputDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating index directory: %w", err)
	}

	path := filepath.Join(dir, "index.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("writing index: %w", err)
	}
	return path, nil
}
