package exec

import (
	"context"
	"runtime"
	"strings"
	"testing"
)

func TestNewRealExecutor(t *testing.T) {
	e := NewRealExecutor()
	if e == nil {
		t.Fatal("expected non-nil executor")
	}
}

func TestRealExecutor_Run(t *testing.T) {
	e := NewRealExecutor()
	out, err := e.Run(context.Background(), "echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

func TestRealExecutor_Run_Error(t *testing.T) {
	e := NewRealExecutor()
	_, err := e.Run(context.Background(), "nonexistent-command-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}

func TestRealExecutor_RunWithStdin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	e := NewRealExecutor()
	out, err := e.RunWithStdin(context.Background(), "hello from stdin", "cat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "hello from stdin" {
		t.Errorf("expected 'hello from stdin', got %q", out)
	}
}
