package codexcli

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

func TestCodexCLI_Compact_Success(t *testing.T) {
	mock := &mockExecutor{output: `{"type":"result","result":"# Compacted"}` + "\n"}
	a, err := New("test-codex", config.AgentConfig{Model: "codex"}, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Name() != "test-codex" {
		t.Errorf("expected 'test-codex', got %q", a.Name())
	}

	input := &agent.CompactInput{SourceName: "test", Content: "content", Template: "template", TotalChunks: 1}
	output, err := a.Compact(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Content != "# Compacted" {
		t.Errorf("expected '# Compacted', got %q", output.Content)
	}
	if !strings.Contains(mock.stdinReceived, "SYSTEM INSTRUCTIONS") {
		t.Error("expected system instructions in stdin")
	}
}

func TestCodexCLI_Compact_CLIError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("codex: not found")}
	a, _ := New("test", config.AgentConfig{}, mock)

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := a.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "codex CLI failed") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCodexCLI_DefaultModel(t *testing.T) {
	mock := &mockExecutor{output: "ok"}
	a, _ := New("test", config.AgentConfig{Model: ""}, mock)
	cli := a.(*CodexCLI)
	if cli.model != "codex" {
		t.Errorf("expected default model 'codex', got %q", cli.model)
	}
}

func TestExtractCodexResult_ResultEvent(t *testing.T) {
	output := `{"type":"message","role":"assistant","content":"partial"}
{"type":"result","result":"final answer"}`

	text, err := extractCodexResult(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "final answer" {
		t.Errorf("expected 'final answer', got %q", text)
	}
}

func TestExtractCodexResult_MessageFallback(t *testing.T) {
	output := `{"type":"message","role":"assistant","content":"assistant response"}`

	text, err := extractCodexResult(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "assistant response" {
		t.Errorf("expected 'assistant response', got %q", text)
	}
}

func TestExtractCodexResult_RawFallback(t *testing.T) {
	text, err := extractCodexResult("just plain text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "just plain text" {
		t.Errorf("expected raw fallback, got %q", text)
	}
}

func TestExtractCodexResult_Empty(t *testing.T) {
	_, err := extractCodexResult("")
	if err == nil {
		t.Fatal("expected error for empty output")
	}
}

func TestCodexCLI_Generate(t *testing.T) {
	mock := &mockExecutor{output: `{"type":"result","result":"gen output"}` + "\n"}
	a, _ := New("test", config.AgentConfig{}, mock)

	text, err := a.Generate(context.Background(), "sys", "usr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "gen output" {
		t.Errorf("expected 'gen output', got %q", text)
	}
}
