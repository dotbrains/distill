package openai

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
}

func TestNew_Success(t *testing.T) {
	t.Setenv("DISTILL_TEST_OAI_KEY", "sk-test")
	a, err := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_TEST_OAI_KEY", Model: "gpt-4o"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "test" {
		t.Errorf("expected 'test', got %q", a.Name())
	}
}

func TestGPT_Compact_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Error("expected Bearer token")
		}
		resp := chatResponse{Choices: []chatChoice{{Message: chatMessage{Content: "# Compacted"}}}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_OAI_KEY", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_OAI_KEY", Model: "gpt-4o"})
	gpt := a.(*GPT)
	gpt.SetBaseURL(server.URL)
	gpt.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "content", Template: "template", TotalChunks: 1}
	output, err := gpt.Compact(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Content != "# Compacted" {
		t.Errorf("expected '# Compacted', got %q", output.Content)
	}
}

func TestGPT_Compact_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal"}`))
	}))
	defer server.Close()

	t.Setenv("DISTILL_OAI_KEY2", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_OAI_KEY2", Model: "gpt-4o"})
	gpt := a.(*GPT)
	gpt.SetBaseURL(server.URL)
	gpt.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := gpt.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected 500 in error, got: %v", err)
	}
}

func TestGPT_Compact_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{Choices: []chatChoice{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_OAI_KEY3", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_OAI_KEY3", Model: "gpt-4o"})
	gpt := a.(*GPT)
	gpt.SetBaseURL(server.URL)
	gpt.SetClient(server.Client())

	input := &agent.CompactInput{SourceName: "test", Content: "x", Template: "y", TotalChunks: 1}
	_, err := gpt.Compact(context.Background(), input)
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
}

func TestGPT_Generate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{Choices: []chatChoice{{Message: chatMessage{Content: "gen"}}}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Setenv("DISTILL_OAI_KEY4", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_OAI_KEY4", Model: "gpt-4o"})
	gpt := a.(*GPT)
	gpt.SetBaseURL(server.URL)
	gpt.SetClient(server.Client())

	text, err := gpt.Generate(context.Background(), "sys", "usr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "gen" {
		t.Errorf("expected 'gen', got %q", text)
	}
}

func TestNew_DefaultMaxTokens(t *testing.T) {
	t.Setenv("DISTILL_OAI_KEY5", "sk-test")
	a, _ := New("test", config.AgentConfig{APIKeyEnv: "DISTILL_OAI_KEY5", Model: "gpt-4o", MaxTokens: 0})
	gpt := a.(*GPT)
	if gpt.maxTokens != 8192 {
		t.Errorf("expected 8192, got %d", gpt.maxTokens)
	}
}
