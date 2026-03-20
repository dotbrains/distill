package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_Builtin(t *testing.T) {
	for _, tt := range BuiltinTemplates {
		content, err := Load(tt.Name, "")
		if err != nil {
			t.Errorf("Load(%q) failed: %v", tt.Name, err)
			continue
		}
		if content == "" {
			t.Errorf("Load(%q) returned empty content", tt.Name)
		}
	}
}

func TestLoad_BuiltinRulesContent(t *testing.T) {
	content, err := Load("rules", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "imperative") {
		t.Error("rules template should mention 'imperative'")
	}
	if !strings.Contains(content, "DIRECTIVES") {
		t.Error("rules template should contain DIRECTIVES")
	}
}

func TestLoad_Unknown(t *testing.T) {
	_, err := Load("nonexistent-template", "")
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
	if !strings.Contains(err.Error(), "unknown template") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoad_Custom(t *testing.T) {
	dir := t.TempDir()
	content := "Custom template content here."
	os.WriteFile(filepath.Join(dir, "my-template.md"), []byte(content), 0o644)

	loaded, err := Load("my-template", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != content {
		t.Errorf("expected %q, got %q", content, loaded)
	}
}

func TestLoad_CustomWithFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\nname: my-template\ndescription: test\n---\nActual content."
	os.WriteFile(filepath.Join(dir, "fm-template.md"), []byte(content), 0o644)

	loaded, err := Load("fm-template", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != "Actual content." {
		t.Errorf("expected frontmatter stripped, got %q", loaded)
	}
}

func TestLoad_CustomOverridesBuiltin(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "rules.md"), []byte("custom rules override"), 0o644)

	loaded, err := Load("rules", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != "custom rules override" {
		t.Errorf("custom should override builtin, got %q", loaded)
	}
}

func TestIsBuiltin(t *testing.T) {
	if !IsBuiltin("rules") {
		t.Error("'rules' should be builtin")
	}
	if !IsBuiltin("principles") {
		t.Error("'principles' should be builtin")
	}
	if !IsBuiltin("patterns") {
		t.Error("'patterns' should be builtin")
	}
	if !IsBuiltin("raw") {
		t.Error("'raw' should be builtin")
	}
	if IsBuiltin("nonexistent") {
		t.Error("'nonexistent' should not be builtin")
	}
}

func TestStripFrontmatter_NoFrontmatter(t *testing.T) {
	input := "Just regular content."
	if got := stripFrontmatter(input); got != input {
		t.Errorf("expected no change, got %q", got)
	}
}

func TestStripFrontmatter_WithFrontmatter(t *testing.T) {
	input := "---\nname: test\n---\nBody content."
	got := stripFrontmatter(input)
	if got != "Body content." {
		t.Errorf("expected 'Body content.', got %q", got)
	}
}

func TestStripFrontmatter_UnclosedFrontmatter(t *testing.T) {
	input := "---\nname: test\nno closing delimiter"
	got := stripFrontmatter(input)
	if got != input {
		t.Errorf("expected original content for unclosed frontmatter, got %q", got)
	}
}
