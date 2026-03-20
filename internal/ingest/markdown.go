package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dotbrains/distill/internal/config"
)

// MarkdownIngestor reads local markdown files.
type MarkdownIngestor struct{}

func (m *MarkdownIngestor) Ingest(src config.Source) ([]Chunk, error) {
	path := expandHome(src.Path)

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("accessing %s: %w", path, err)
	}

	if info.IsDir() {
		return m.ingestDir(path)
	}
	return m.ingestFile(path)
}

func (m *MarkdownIngestor) ingestFile(path string) ([]Chunk, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return nil, fmt.Errorf("file %s is empty", path)
	}
	return []Chunk{{Content: content}}, nil
}

func (m *MarkdownIngestor) ingestDir(dir string) ([]Chunk, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".md") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		return nil, fmt.Errorf("no .md files found in %s", dir)
	}

	var chunks []Chunk
	for i, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		content := strings.TrimSpace(string(data))
		if content == "" {
			continue
		}
		chunks = append(chunks, Chunk{
			Content:    content,
			ChapterNum: i + 1,
		})
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("all .md files in %s are empty", dir)
	}
	return chunks, nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
