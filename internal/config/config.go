package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the top-level distill configuration.
type Config struct {
	DefaultAgent string                 `yaml:"default_agent"`
	Agents       map[string]AgentConfig `yaml:"agents"`
	Output       OutputConfig           `yaml:"output"`
	Sources      map[string]Source       `yaml:"sources"`
}

// AgentConfig defines a single AI agent provider.
type AgentConfig struct {
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`
	MaxTokens int    `yaml:"max_tokens,omitempty"`
}

// OutputConfig controls where compacted output is written.
type OutputConfig struct {
	Dir              string `yaml:"dir"`
	GenerateIndexes  bool   `yaml:"generate_indexes"`
	TokenBudget      int    `yaml:"token_budget"`
	CustomTemplates  string `yaml:"custom_templates_dir,omitempty"`
}

// Source defines a single source to be compacted.
type Source struct {
	Type          string `yaml:"type"`                    // pdf, markdown, notion, url, epub, github
	Path          string `yaml:"path,omitempty"`          // local file/dir path
	URL           string `yaml:"url,omitempty"`           // remote URL (notion, url, github)
	Repo          string `yaml:"repo,omitempty"`          // github repo (owner/repo)
	Ref           string `yaml:"ref,omitempty"`           // github ref
	Template      string `yaml:"template"`                // rules, principles, patterns, raw, or custom name
	OutputDir     string `yaml:"output_dir"`              // subdirectory under output.dir
	OutputFile    string `yaml:"output_file,omitempty"`   // single output filename
	OutputPattern string `yaml:"output_pattern,omitempty"` // pattern for multi-file output
	SplitBy       string `yaml:"split_by,omitempty"`      // "chapter" for per-chapter splitting
	ChapterPattern string `yaml:"chapter_pattern,omitempty"` // regex for chapter detection
}

// CLIProviders are providers that use local CLI binaries and don't need API keys.
var CLIProviders = map[string]bool{
	"claude-cli": true,
	"codex-cli":  true,
}

// DefaultConfig returns the built-in default configuration.
func DefaultConfig() *Config {
	return &Config{
		DefaultAgent: "claude-cli",
		Agents: map[string]AgentConfig{
			"claude-cli": {
				Provider: "claude-cli",
				Model:    "sonnet",
			},
			"codex-cli": {
				Provider: "codex-cli",
				Model:    "codex",
			},
			"claude-api": {
				Provider:  "anthropic",
				Model:     "claude-sonnet-4-20250514",
				APIKeyEnv: "ANTHROPIC_API_KEY",
				MaxTokens: 8192,
			},
			"gpt-api": {
				Provider:  "openai",
				Model:     "gpt-4o",
				APIKeyEnv: "OPENAI_API_KEY",
				MaxTokens: 8192,
			},
		},
		Output: OutputConfig{
			Dir:             "./output",
			GenerateIndexes: true,
			TokenBudget:     4000,
		},
		Sources: map[string]Source{},
	}
}

// GlobalConfigDir returns the global config directory path.
func GlobalConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "distill"), nil
}

// GlobalConfigPath returns the full path to the global config file.
func GlobalConfigPath() (string, error) {
	dir, err := GlobalConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// ProjectConfigPath returns "distill.yaml" in the current directory.
func ProjectConfigPath() string {
	return "distill.yaml"
}

// Load reads the config, checking project-level first, then global, then defaults.
func Load() (*Config, error) {
	// Try project-level first
	projectPath := ProjectConfigPath()
	if _, err := os.Stat(projectPath); err == nil {
		return LoadFrom(projectPath)
	}

	// Try global
	globalPath, err := GlobalConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}
	if _, err := os.Stat(globalPath); err == nil {
		return LoadFrom(globalPath)
	}

	return DefaultConfig(), nil
}

// LoadFrom reads the config from a specific path.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, nil
}

// SaveTo writes the config to a specific path.
func SaveTo(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
