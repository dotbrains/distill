package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DefaultAgent != "claude-cli" {
		t.Errorf("expected default agent 'claude-cli', got %q", cfg.DefaultAgent)
	}
	if len(cfg.Agents) != 4 {
		t.Errorf("expected 4 agents, got %d", len(cfg.Agents))
	}
	if _, ok := cfg.Agents["claude-cli"]; !ok {
		t.Error("expected claude-cli agent in defaults")
	}
	if _, ok := cfg.Agents["codex-cli"]; !ok {
		t.Error("expected codex-cli agent in defaults")
	}
	if _, ok := cfg.Agents["claude-api"]; !ok {
		t.Error("expected claude-api agent in defaults")
	}
	if _, ok := cfg.Agents["gpt-api"]; !ok {
		t.Error("expected gpt-api agent in defaults")
	}
	if cfg.Output.Dir != "./output" {
		t.Errorf("expected output dir './output', got %q", cfg.Output.Dir)
	}
	if cfg.Output.TokenBudget != 4000 {
		t.Errorf("expected token budget 4000, got %d", cfg.Output.TokenBudget)
	}
	if !cfg.Output.GenerateIndexes {
		t.Error("expected generate_indexes to be true")
	}
	if cfg.Sources == nil {
		t.Error("expected Sources to be initialized")
	}
}

func TestSaveToAndLoadFrom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "distill.yaml")

	cfg := DefaultConfig()
	cfg.Sources = map[string]Source{
		"test-source": {
			Type:      "markdown",
			Path:      "./test.md",
			Template:  "rules",
			OutputDir: "test",
		},
	}

	if err := SaveTo(cfg, path); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if loaded.DefaultAgent != "claude-cli" {
		t.Errorf("loaded default agent mismatch: %q", loaded.DefaultAgent)
	}
	src, ok := loaded.Sources["test-source"]
	if !ok {
		t.Fatal("source 'test-source' not found after load")
	}
	if src.Type != "markdown" {
		t.Errorf("expected type 'markdown', got %q", src.Type)
	}
	if src.Path != "./test.md" {
		t.Errorf("expected path './test.md', got %q", src.Path)
	}
}

func TestLoadFrom_NotExist(t *testing.T) {
	cfg, err := LoadFrom("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg.DefaultAgent != "claude-cli" {
		t.Error("expected defaults when file not found")
	}
}

func TestLoadFrom_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte("default_agent: [invalid\n  broken: {{"), 0o644)

	_, err := LoadFrom(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSaveTo_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "config.yaml")

	cfg := DefaultConfig()
	if err := SaveTo(cfg, path); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("nested config file not created: %v", err)
	}
}

func TestProjectConfigPath(t *testing.T) {
	if ProjectConfigPath() != "distill.yaml" {
		t.Errorf("expected 'distill.yaml', got %q", ProjectConfigPath())
	}
}

func TestGlobalConfigPath(t *testing.T) {
	path, err := GlobalConfigPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected config.yaml, got %q", filepath.Base(path))
	}
	if filepath.Base(filepath.Dir(path)) != "distill" {
		t.Errorf("expected distill dir, got %q", filepath.Base(filepath.Dir(path)))
	}
}

func TestLoad_ProjectLevel(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Write a project config with a custom source
	cfg := DefaultConfig()
	cfg.Sources = map[string]Source{"proj": {Type: "markdown", Template: "rules", OutputDir: "x"}}
	SaveTo(cfg, "distill.yaml")

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if _, ok := loaded.Sources["proj"]; !ok {
		t.Error("expected project-level source to be loaded")
	}
}

func TestLoad_FallsBackToDefaults(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// No config file exists
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.DefaultAgent != "claude-cli" {
		t.Error("expected defaults")
	}
}

func TestGlobalConfigDir(t *testing.T) {
	dir, err := GlobalConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(dir) != "distill" {
		t.Errorf("expected 'distill', got %q", filepath.Base(dir))
	}
}

func TestSaveTo_InvalidPath(t *testing.T) {
	err := SaveTo(DefaultConfig(), "/dev/null/impossible/config.yaml")
	if err == nil {
		t.Error("expected error saving to invalid path")
	}
}

func TestCLIProviders(t *testing.T) {
	if !CLIProviders["claude-cli"] {
		t.Error("expected claude-cli to be a CLI provider")
	}
	if !CLIProviders["codex-cli"] {
		t.Error("expected codex-cli to be a CLI provider")
	}
	if CLIProviders["anthropic"] {
		t.Error("anthropic should not be a CLI provider")
	}
}
