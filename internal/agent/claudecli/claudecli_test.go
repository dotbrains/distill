package claudecli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/dotbrains/distill/internal/agent"
	"github.com/dotbrains/distill/internal/config"
)

type mockExecutor struct {
	stdinReceived string
	output        string
	err           error
}

func (m *mockExecutor) Run(ctx context.Context, name string, args ...string) (string, error) {
	return m.output, m.err
}

func (m *mockExecutor) RunWithStdin(ctx context.Context, stdin string, name string, args ...string) (string, error) {
	m.stdinReceived = stdin
	return m.output, m.err
}

func TestClaudeCLI_Compact_Success(t *testing.T) {
	cliOutput := `{"type":"result","result":"# Compacted Output\n\n**1.1 Rule** - Do this.","is_error":false}`
	mock := &mockExecutor{output: cliOutput}
	a, err := New("test-claude", config.AgentConfig{Model: "sonnet"}, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Name() != "test-claude" {
		t.Errorf("expected name 'test-claude', got %q", a.Name())
	}

	input := &agent.CompactInput{
		SourceName:  "tao-of-react",
		ChunkIndex:  0,
		TotalChunks: 1,
		Content:     "some source content",
		Template:    "You are a compactor.",
		TokenBudget: 4000,
	}

	output, err := a.Compact(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output.Content, "Compacted Output") {
		t.Errorf("expected compacted content, got %q", output.Content)
	}
	if mock.stdinReceived == "" {
		t.Error("expected stdin to be sent")
	}
}

func TestClaudeCLI_Compact_CLIError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("claude: command not found")}
	a, _ := New("test", config.AgentConfig{}, mock)

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := a.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "claude CLI failed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeCLI_Compact_IsError(t *testing.T) {
	mock := &mockExecutor{output: `{"type":"result","result":"rate limited","is_error":true}`}
	a, _ := New("test", config.AgentConfig{}, mock)

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := a.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limited") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeCLI_Compact_RawFallback(t *testing.T) {
	mock := &mockExecutor{output: "# Raw markdown output"}
	a, _ := New("test", config.AgentConfig{}, mock)

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	output, err := a.Compact(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Content != "# Raw markdown output" {
		t.Errorf("expected raw fallback, got %q", output.Content)
	}
}

func TestClaudeCLI_Generate(t *testing.T) {
	mock := &mockExecutor{output: `{"type":"result","result":"generated text","is_error":false}`}
	a, _ := New("test", config.AgentConfig{}, mock)

	text, err := a.Generate(context.Background(), "system", "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "generated text" {
		t.Errorf("expected 'generated text', got %q", text)
	}
}

func TestClaudeCLI_DefaultModel(t *testing.T) {
	mock := &mockExecutor{output: "ok"}
	a, _ := New("test", config.AgentConfig{Model: ""}, mock)
	cli := a.(*ClaudeCLI)
	if cli.model != "sonnet" {
		t.Errorf("expected default model 'sonnet', got %q", cli.model)
	}
}
