package state

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestHashContent(t *testing.T) {
	h1 := HashContent("hello")
	h2 := HashContent("hello")
	h3 := HashContent("world")

	if h1 != h2 {
		t.Error("same input should produce same hash")
	}
	if h1 == h3 {
		t.Error("different input should produce different hash")
	}
	if !strings.HasPrefix(h1, "sha256:") {
		t.Errorf("expected sha256: prefix, got %q", h1)
	}
}

func TestState_IsDirty_NewSource(t *testing.T) {
	s := &State{Sources: map[string]SourceState{}}
	if !s.IsDirty("new-source", "hash", "tmpl", "agent") {
		t.Error("new source should be dirty")
	}
}

func TestState_IsDirty_Unchanged(t *testing.T) {
	s := &State{Sources: map[string]SourceState{
		"src": {ContentHash: "h1", TemplateHash: "t1", Agent: "a1"},
	}}
	if s.IsDirty("src", "h1", "t1", "a1") {
		t.Error("unchanged source should not be dirty")
	}
}

func TestState_IsDirty_ContentChanged(t *testing.T) {
	s := &State{Sources: map[string]SourceState{
		"src": {ContentHash: "h1", TemplateHash: "t1", Agent: "a1"},
	}}
	if !s.IsDirty("src", "h2", "t1", "a1") {
		t.Error("content change should be dirty")
	}
}

func TestState_IsDirty_TemplateChanged(t *testing.T) {
	s := &State{Sources: map[string]SourceState{
		"src": {ContentHash: "h1", TemplateHash: "t1", Agent: "a1"},
	}}
	if !s.IsDirty("src", "h1", "t2", "a1") {
		t.Error("template change should be dirty")
	}
}

func TestState_IsDirty_AgentChanged(t *testing.T) {
	s := &State{Sources: map[string]SourceState{
		"src": {ContentHash: "h1", TemplateHash: "t1", Agent: "a1"},
	}}
	if !s.IsDirty("src", "h1", "t1", "a2") {
		t.Error("agent change should be dirty")
	}
}

func TestState_Update(t *testing.T) {
	s := &State{Sources: map[string]SourceState{}}
	s.Update("src", SourceState{
		ContentHash:   "h1",
		TemplateHash:  "t1",
		Agent:         "claude",
		LastCompacted: time.Now(),
		OutputTokens:  2847,
		OutputFiles:   []string{"tao/file.md"},
	})

	entry, ok := s.Sources["src"]
	if !ok {
		t.Fatal("source not found after update")
	}
	if entry.ContentHash != "h1" {
		t.Errorf("expected 'h1', got %q", entry.ContentHash)
	}
	if entry.OutputTokens != 2847 {
		t.Errorf("expected 2847, got %d", entry.OutputTokens)
	}
}

func TestState_Update_NilMap(t *testing.T) {
	s := &State{}
	s.Update("src", SourceState{ContentHash: "h"})
	if _, ok := s.Sources["src"]; !ok {
		t.Error("update should initialize nil map")
	}
}

func TestState_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	s := &State{Sources: map[string]SourceState{
		"test": {
			ContentHash:   "hash123",
			TemplateHash:  "tmpl456",
			Agent:         "claude-cli",
			LastCompacted: time.Date(2025, 3, 20, 14, 0, 0, 0, time.UTC),
			OutputTokens:  2847,
			OutputFiles:   []string{"output/test-minified.md"},
		},
	}}

	if err := s.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	entry, ok := loaded.Sources["test"]
	if !ok {
		t.Fatal("source 'test' not found after load")
	}
	if entry.ContentHash != "hash123" {
		t.Errorf("expected 'hash123', got %q", entry.ContentHash)
	}
	if entry.Agent != "claude-cli" {
		t.Errorf("expected 'claude-cli', got %q", entry.Agent)
	}
}

func TestState_Save_InvalidDir(t *testing.T) {
	// Try to save in a read-only / non-writable location
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Make directory read-only
	os.Chmod(dir, 0o444)
	defer os.Chmod(dir, 0o755)

	s := &State{Sources: map[string]SourceState{
		"test": {ContentHash: "h"},
	}}
	err := s.Save()
	if err == nil {
		t.Error("expected error saving to read-only directory")
	}
}

func TestLoad_NoFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	s, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Sources) != 0 {
		t.Error("expected empty sources for missing state file")
	}
}
