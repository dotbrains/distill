package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotbrains/distill/internal/config"
)

func TestNew_Markdown(t *testing.T) {
	i, err := New("markdown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i == nil {
		t.Fatal("expected non-nil ingestor")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNew_UnimplementedTypes(t *testing.T) {
	for _, typ := range []string{"pdf", "notion", "url", "epub", "github"} {
		_, err := New(typ)
		if err == nil {
			t.Errorf("expected error for unimplemented type %q", typ)
		}
		if !strings.Contains(err.Error(), "not yet implemented") {
			t.Errorf("expected 'not yet implemented' for %q, got: %v", typ, err)
		}
	}
}

func TestMarkdownIngestor_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "guide.md")
	os.WriteFile(path, []byte("# Guide\n\nSome content here."), 0o644)

	i := &MarkdownIngestor{}
	chunks, err := i.Ingest(config.Source{Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if !strings.Contains(chunks[0].Content, "# Guide") {
		t.Errorf("expected content to contain '# Guide', got %q", chunks[0].Content)
	}
}

func TestMarkdownIngestor_Directory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "01-intro.md"), []byte("# Intro"), 0o644)
	os.WriteFile(filepath.Join(dir, "02-body.md"), []byte("# Body"), 0o644)
	os.WriteFile(filepath.Join(dir, "03-conclusion.md"), []byte("# Conclusion"), 0o644)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not markdown"), 0o644) // should be ignored

	i := &MarkdownIngestor{}
	chunks, err := i.Ingest(config.Source{Path: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	// Should be sorted alphabetically
	if !strings.Contains(chunks[0].Content, "Intro") {
		t.Error("first chunk should be intro")
	}
	if chunks[0].ChapterNum != 1 {
		t.Errorf("expected chapter 1, got %d", chunks[0].ChapterNum)
	}
	if chunks[2].ChapterNum != 3 {
		t.Errorf("expected chapter 3, got %d", chunks[2].ChapterNum)
	}
}

func TestMarkdownIngestor_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.md")
	os.WriteFile(path, []byte(""), 0o644)

	i := &MarkdownIngestor{}
	_, err := i.Ingest(config.Source{Path: path})
	if err == nil {
		t.Fatal("expected error for empty file")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMarkdownIngestor_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	i := &MarkdownIngestor{}
	_, err := i.Ingest(config.Source{Path: dir})
	if err == nil {
		t.Fatal("expected error for empty directory")
	}
}

func TestMarkdownIngestor_MissingPath(t *testing.T) {
	i := &MarkdownIngestor{}
	_, err := i.Ingest(config.Source{Path: "/nonexistent/path"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestMarkdownIngestor_DirectorySkipsEmptyFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("content"), 0o644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte("  \n  "), 0o644) // whitespace only

	i := &MarkdownIngestor{}
	chunks, err := i.Ingest(config.Source{Path: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk (empty skipped), got %d", len(chunks))
	}
}

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	got := expandHome("~/test.md")
	if got != filepath.Join(home, "test.md") {
		t.Errorf("expected %q, got %q", filepath.Join(home, "test.md"), got)
	}

	// Non-home path unchanged
	if expandHome("/absolute/path") != "/absolute/path" {
		t.Error("absolute path should be unchanged")
	}
	if expandHome("relative/path") != "relative/path" {
		t.Error("relative path should be unchanged")
	}
}
