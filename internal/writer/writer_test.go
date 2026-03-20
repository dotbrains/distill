package writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteOutput(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteOutput(dir, "tao", "tao-of-react-minified.md", "# Compacted\n\n**1.1 Rule**")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "tao", "tao-of-react-minified.md")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if string(data) != "# Compacted\n\n**1.1 Rule**" {
		t.Errorf("content mismatch: %q", string(data))
	}
}

func TestWriteOutput_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteOutput(dir, "deep/nested/dir", "file.md", "content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestWriteOutput_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	WriteOutput(dir, "sub", "file.md", "original")
	path, err := WriteOutput(dir, "sub", "file.md", "updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "updated" {
		t.Errorf("expected overwritten content, got %q", string(data))
	}
}

func TestWriteOutput_InvalidDir(t *testing.T) {
	// Writing to /dev/null/impossible should fail
	_, err := WriteOutput("/dev/null", "sub", "file.md", "content")
	if err == nil {
		t.Error("expected error writing to invalid directory")
	}
}

func TestWriteIndex_InvalidDir(t *testing.T) {
	_, err := WriteIndex("/dev/null", "sub", "content")
	if err == nil {
		t.Error("expected error writing to invalid directory")
	}
}

func TestWriteIndex(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteIndex(dir, "tao", "# Index\n\nFiles listed here.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "tao", "index.md")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading index: %v", err)
	}
	if string(data) != "# Index\n\nFiles listed here." {
		t.Errorf("content mismatch: %q", string(data))
	}
}
