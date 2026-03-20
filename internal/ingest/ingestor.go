package ingest

import (
	"fmt"

	"github.com/dotbrains/distill/internal/config"
)

// Chunk represents a logical section of ingested content.
type Chunk struct {
	Content      string
	ChapterNum   int
	ChapterTitle string
}

// Ingestor reads source material and returns raw text chunks.
type Ingestor interface {
	Ingest(src config.Source) ([]Chunk, error)
}

// New returns the appropriate ingestor for the given source type.
func New(sourceType string) (Ingestor, error) {
	switch sourceType {
	case "markdown":
		return &MarkdownIngestor{}, nil
	case "pdf":
		return nil, fmt.Errorf("pdf ingestor not yet implemented")
	case "notion":
		return nil, fmt.Errorf("notion ingestor not yet implemented")
	case "url":
		return nil, fmt.Errorf("url ingestor not yet implemented")
	case "epub":
		return nil, fmt.Errorf("epub ingestor not yet implemented")
	case "github":
		return nil, fmt.Errorf("github ingestor not yet implemented")
	default:
		return nil, fmt.Errorf("unknown source type %q", sourceType)
	}
}
