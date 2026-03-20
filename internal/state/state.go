package state

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const StateFile = ".distill-state.yaml"

// State tracks compaction state for all sources.
type State struct {
	Sources map[string]SourceState `yaml:"sources"`
}

// SourceState tracks a single source's compaction state.
type SourceState struct {
	ContentHash   string    `yaml:"content_hash"`
	TemplateHash  string    `yaml:"template_hash"`
	Agent         string    `yaml:"agent"`
	LastCompacted time.Time `yaml:"last_compacted"`
	OutputTokens  int       `yaml:"output_tokens"`
	OutputFiles   []string  `yaml:"output_files"`
}

// Load reads the state file from disk.
func Load() (*State, error) {
	data, err := os.ReadFile(StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Sources: map[string]SourceState{}}, nil
		}
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	s := &State{Sources: map[string]SourceState{}}
	if err := yaml.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}
	return s, nil
}

// Save writes the state file to disk.
func (s *State) Save() error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	if err := os.WriteFile(StateFile, data, 0o644); err != nil {
		return fmt.Errorf("writing state file: %w", err)
	}
	return nil
}

// IsDirty returns true if the source needs re-compaction.
func (s *State) IsDirty(name, contentHash, templateHash, agentID string) bool {
	prev, ok := s.Sources[name]
	if !ok {
		return true
	}
	return prev.ContentHash != contentHash ||
		prev.TemplateHash != templateHash ||
		prev.Agent != agentID
}

// Update records the compaction result for a source.
func (s *State) Update(name string, entry SourceState) {
	if s.Sources == nil {
		s.Sources = map[string]SourceState{}
	}
	s.Sources[name] = entry
}

// HashContent returns a SHA-256 hash of the given content.
func HashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("sha256:%x", h)
}
