package agent

import (
	"context"
	"testing"

	"github.com/dotbrains/distill/internal/config"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 5, "hello..."},
		{"", 5, ""},
		{"abc", 0, "..."},
	}
	for _, tt := range tests {
		got := Truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

// mockAgent implements Agent for testing the registry.
type mockAgent struct{ name string }

func (m *mockAgent) Name() string { return m.name }
func (m *mockAgent) Compact(ctx context.Context, input *CompactInput) (*CompactOutput, error) {
	return &CompactOutput{Content: "mocked"}, nil
}
func (m *mockAgent) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return "mocked", nil
}

func TestRegistry_NewAgent_Known(t *testing.T) {
	RegisterProvider("test-provider", func(name string, cfg config.AgentConfig) (Agent, error) {
		return &mockAgent{name: name}, nil
	})
	defer func() { delete(providers, "test-provider") }()

	a, err := NewAgent("my-agent", config.AgentConfig{Provider: "test-provider"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "my-agent" {
		t.Errorf("expected name 'my-agent', got %q", a.Name())
	}
}

func TestRegistry_NewAgent_Unknown(t *testing.T) {
	_, err := NewAgent("x", config.AgentConfig{Provider: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestRegistry_NewAgentFromConfig(t *testing.T) {
	RegisterProvider("test-provider2", func(name string, cfg config.AgentConfig) (Agent, error) {
		return &mockAgent{name: name}, nil
	})
	defer func() { delete(providers, "test-provider2") }()

	cfg := &config.Config{
		DefaultAgent: "default-agent",
		Agents: map[string]config.AgentConfig{
			"default-agent": {Provider: "test-provider2"},
			"other-agent":   {Provider: "test-provider2"},
		},
	}

	// Use default
	a, err := NewAgentFromConfig("", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "default-agent" {
		t.Errorf("expected 'default-agent', got %q", a.Name())
	}

	// Use named
	a, err = NewAgentFromConfig("other-agent", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "other-agent" {
		t.Errorf("expected 'other-agent', got %q", a.Name())
	}

	// Missing agent
	_, err = NewAgentFromConfig("missing", cfg)
	if err == nil {
		t.Fatal("expected error for missing agent")
	}
}
