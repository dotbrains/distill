package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dotbrains/distill/internal/agent"
	"github.com/dotbrains/distill/internal/config"
)

func TestNew_MissingAPIKey(t *testing.T) {
	_, err := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_NONEXISTENT_KEY"})
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNew_Success(t *testing.T) {
	t.Setenv("DISTILL_TEST_ANTHROPIC_KEY", "sk-test")
	a, err := New("test", config.AgentConfig{
		APIKeyEnv: "DISTILL_TEST_ANTHROPIC_KEY",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "test" {
		t.Errorf("expected name 'test', got %q", a.Name())
	}
}

func TestClaude_Compact_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "sk-test" {
			t.Error("expected API key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Error("expected anthropic-version header")
		}

		resp := messagesResponse{
			Content: []contentBlock{
				{Type: "text", Text: "# Compacted output"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_TEST_KEY", "sk-test")
	a, _ := New("test", config.AgentConfig{
		APIKeyEnv: "DISTILL_TEST_KEY",
		Model:     "claude-sonnet-4-20250514",
	})
	claude := a.(*Claude)
	claude.SetBaseURL(server.URL)
	claude.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "content", Template: "template", TotalChunks: 1}
	output, err := claude.Compact(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Content != "# Compacted output" {
		t.Errorf("expected compacted content, got %q", output.Content)
	}
}

func TestClaude_Compact_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer server.Close()

	t.Setenv("DISTILL_TEST_KEY2", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_KEY2", Model: "model"})
	claude := a.(*Claude)
	claude.SetBaseURL(server.URL)
	claude.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := claude.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("expected status 429 in error, got: %v", err)
	}
}

func TestClaude_Compact_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := messagesResponse{Content: []contentBlock{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_TEST_KEY3", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_KEY3", Model: "model"})
	claude := a.(*Claude)
	claude.SetBaseURL(server.URL)
	claude.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := claude.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error for empty response")
	}
}

func TestClaude_Generate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := messagesResponse{Content: []contentBlock{{Type: "text", Text: "generated"}}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_TEST_KEY4", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_KEY4", Model: "model"})
	claude := a.(*Claude)
	claude.SetBaseURL(server.URL)
	claude.SetClient(server.Client())

	text, err := claude.Generate(context.Background(), "sys", "usr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "generated" {
		t.Errorf("expected 'generated', got %q", text)
	}
}

func TestNew_DefaultMaxTokens(t *testing.T) {
	t.Setenv("DISTILL_TEST_KEY5", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_KEY5", Model: "model", MaxTokens: 0})
	claude := a.(*Claude)
	if claude.maxTokens != 8192 {
		t.Errorf("expected default max tokens 8192, got %d", claude.maxTokens)
	}
}
