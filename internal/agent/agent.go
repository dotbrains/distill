package agent

import "context"

// Agent is the interface all AI compaction providers must implement.
type Agent interface {
	// Name returns the agent's configured name.
	Name() string

	// Compact sends source content to the AI and returns compacted markdown.
	Compact(ctx context.Context, input *CompactInput) (*CompactOutput, error)

	// Generate sends a system+user prompt and returns raw text.
	// Used by index generation and other non-compaction features.
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// CompactInput contains everything the agent needs to compact a chunk.
type CompactInput struct {
	SourceName  string         // e.g., "tao-of-react"
	ChunkIndex  int            // 0-indexed chunk number
	TotalChunks int            // total chunks for this source
	Content     string         // raw text of this chunk
	Template    string         // system prompt from the template
	TokenBudget int            // target max tokens for output
	Metadata    SourceMetadata // source type, chapter info, etc.
}

// CompactOutput is the result from an AI compaction.
type CompactOutput struct {
	Content    string // compacted markdown
	TokenCount int    // estimated token count of output
}

// SourceMetadata provides context about the source being compacted.
type SourceMetadata struct {
	Type         string // pdf, notion, markdown, url, epub, github
	ChapterNum   int    // for split_by: chapter
	ChapterTitle string
	SectionCtx   string // surrounding section context for continuity
}

// Truncate returns s truncated to maxLen with an ellipsis if needed.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
