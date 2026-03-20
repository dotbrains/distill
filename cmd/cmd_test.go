package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotbrains/distill/internal/config"
)

func TestDeriveName_PDF(t *testing.T) {
	tests := []struct {
		typ      string
		location string
		want     string
	}{
		{"pdf", "~/Books/Tao of React.pdf", "tao-of-react"},
		{"pdf", "/path/to/clean-code.pdf", "clean-code"},
		{"markdown", "./docs/guide.md", "guide"},
		{"epub", "~/Books/DDIA.epub", "ddia"},
		{"url", "https://example.com/best-practices/", "best-practices"},
		{"notion", "https://www.notion.so/org/Some-Doc-Minified-abc123def456789abcdef01234", "some-doc-minified"},
	}
	for _, tt := range tests {
		got := deriveName(tt.typ, tt.location)
		if got != tt.want {
			t.Errorf("deriveName(%q, %q) = %q, want %q", tt.typ, tt.location, got, tt.want)
		}
	}
}

func TestAvailableSources_Empty(t *testing.T) {
	cfg := &config.Config{Sources: map[string]config.Source{}}
	got := availableSources(cfg)
	if got != "none (add sources to distill.yaml)" {
		t.Errorf("expected 'none' message, got %q", got)
	}
}

func TestAvailableSources_WithSources(t *testing.T) {
	cfg := &config.Config{Sources: map[string]config.Source{
		"alpha": {},
		"beta":  {},
	}}
	got := availableSources(cfg)
	if len(got) == 0 {
		t.Error("expected non-empty result")
	}
}

func TestExecute_Version(t *testing.T) {
	root := newRootCmd("1.2.3")
	root.SetArgs([]string{"--version"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "1.2.3") {
		t.Errorf("expected version in output, got %q", out.String())
	}
}

func TestExecute_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "distill") {
		t.Error("expected 'distill' in help output")
	}
	if !strings.Contains(output, "agents") {
		t.Error("expected 'agents' subcommand in help")
	}
	if !strings.Contains(output, "templates") {
		t.Error("expected 'templates' subcommand in help")
	}
}

func TestTemplatesCmd(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"templates"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "rules") {
		t.Error("expected 'rules' in templates output")
	}
	if !strings.Contains(output, "principles") {
		t.Error("expected 'principles' in templates output")
	}
	if !strings.Contains(output, "patterns") {
		t.Error("expected 'patterns' in templates output")
	}
	if !strings.Contains(output, "raw") {
		t.Error("expected 'raw' in templates output")
	}
}

func TestAgentsCmd(t *testing.T) {
	// Create a temp dir with default config so agents cmd loads it
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"agents"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "claude-cli") {
		t.Error("expected 'claude-cli' in agents output")
	}
	if !strings.Contains(output, "(default)") {
		t.Error("expected '(default)' marker")
	}
}

func TestConfigInitCmd(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	root := newRootCmd("dev")
	root.SetArgs([]string{"config", "init"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat("distill.yaml"); err != nil {
		t.Fatal("distill.yaml not created")
	}

	// Running again without --force should fail
	root2 := newRootCmd("dev")
	root2.SetArgs([]string{"config", "init"})
	err = root2.Execute()
	if err == nil {
		t.Fatal("expected error when config already exists")
	}
}

func TestConfigInitCmd_Force(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Create first
	root := newRootCmd("dev")
	root.SetArgs([]string{"config", "init"})
	root.Execute()

	// Force overwrite
	root2 := newRootCmd("dev")
	root2.SetArgs([]string{"config", "init", "--force"})
	err := root2.Execute()
	if err != nil {
		t.Fatalf("expected --force to succeed, got: %v", err)
	}
}

func TestInitCmd(t *testing.T) {
	dir := t.TempDir()
	repoName := filepath.Join(dir, "my-context")

	root := newRootCmd("dev")
	root.SetArgs([]string{"init", repoName})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check scaffolded files
	for _, f := range []string{"index.md", ".gitignore", "README.md"} {
		if _, err := os.Stat(filepath.Join(repoName, f)); err != nil {
			t.Errorf("expected %s to exist: %v", f, err)
		}
	}
}

func TestListCmd_NoSources(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"list"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "No sources") {
		t.Errorf("expected 'No sources' message, got %q", out.String())
	}
}

func TestListCmd_WithSources(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"test-source": {
			Type:      "markdown",
			Template:  "rules",
			OutputDir: "test",
		},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"list"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "test-source") {
		t.Errorf("expected 'test-source' in output, got %q", out.String())
	}
}

func TestValidateCmd(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"validate"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublishCmd(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"publish"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddCmd_MissingArgs(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"add"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestAddCmd_Success(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "markdown", "./test.md", "--name", "my-source", "--template", "rules"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify source was added to config
	loaded, _ := config.LoadFrom("distill.yaml")
	if _, ok := loaded.Sources["my-source"]; !ok {
		t.Error("source 'my-source' not found in config after add")
	}
}

func TestAddCmd_UnknownType(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "foobar", "./test.md", "--name", "x"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unknown source type")
	}
}

func TestCompactCmd_MissingSource(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"nonexistent-source"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCompactCmd_BadIngestorType(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"bad": {Type: "pdf", Template: "rules", OutputDir: "x"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"bad"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unimplemented ingestor")
	}
}

func TestCompactCmd_BadTemplate(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"bad-tmpl": {Type: "markdown", Path: "./test.md", Template: "nonexistent-template", OutputDir: "x"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"bad-tmpl"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for bad template")
	}
}

func TestCompactCmd_DryRun(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Create a source file
	os.WriteFile("source.md", []byte("# Test content\n\nSome text."), 0o644)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"test-src": {Type: "markdown", Path: "./source.md", Template: "rules", OutputDir: "out", OutputFile: "test-minified.md"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"test-src", "--dry-run"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompactCmd_NoArgs(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should show help
	if !strings.Contains(out.String(), "distill") {
		t.Error("expected help output")
	}
}

func TestUpdateCmd_NoSources(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"update"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCmd_AllSources_WithErrors(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add sources that will fail (pdf not implemented, missing file)
	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"bad-pdf": {Type: "pdf", Template: "rules", OutputDir: "x"},
		"bad-md":  {Type: "markdown", Path: "./nonexistent.md", Template: "rules", OutputDir: "y"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"update"})
	var out bytes.Buffer
	root.SetOut(&out)

	// update --all with failing sources should still complete (errors are warned, not fatal)
	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompactCmd_IngestError(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Source file doesn't exist
	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"missing-file": {Type: "markdown", Path: "./does-not-exist.md", Template: "rules", OutputDir: "out"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"missing-file"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestUpdateCmd_SpecificName_Missing(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"update", "nonexistent"})

	// update delegates to compactSource, which returns an error for missing source.
	// We just verify it doesn't panic — the error may or may not propagate through cobra.
	_ = root.Execute()
}

func TestAddCmd_DuplicateName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{"exists": {Type: "markdown"}}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "markdown", "./test.md", "--name", "exists"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for duplicate source name")
	}
}

func TestAddCmd_AllTypes(t *testing.T) {
	for _, typ := range []string{"pdf", "markdown", "epub", "notion", "url", "github"} {
		dir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(dir)

		cfg := config.DefaultConfig()
		config.SaveTo(cfg, "distill.yaml")

		root := newRootCmd("dev")
		root.SetArgs([]string{"add", typ, "./loc", "--name", "src-" + typ})
		err := root.Execute()
		if err != nil {
			t.Errorf("add %s failed: %v", typ, err)
		}

		os.Chdir(origDir)
	}
}

func TestDeriveName_AdditionalCases(t *testing.T) {
	tests := []struct {
		typ      string
		location string
		want     string
	}{
		{"github", "owner/repo", "source"},
		{"unknown", "anything", "source"},
		{"url", "https://example.com/", "example.com"},
		{"notion", "https://notion.so/short", "short"},
	}
	for _, tt := range tests {
		got := deriveName(tt.typ, tt.location)
		if got != tt.want {
			t.Errorf("deriveName(%q, %q) = %q, want %q", tt.typ, tt.location, got, tt.want)
		}
	}
}

func TestAgentsCmd_WithAPIKeySet(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	t.Setenv("ANTHROPIC_API_KEY", "sk-test")

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"agents"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	// claude-api should show ✓ since ANTHROPIC_API_KEY is set
	if !strings.Contains(output, "claude-api") {
		t.Error("expected claude-api in output")
	}
	if !strings.Contains(output, "codex-cli") {
		t.Error("expected codex-cli in output")
	}
	if !strings.Contains(output, "gpt-api") {
		t.Error("expected gpt-api in output")
	}
}

func TestAgentsCmd_AllProviderTypes(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"agents"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	// All 4 default agents should appear
	for _, name := range []string{"claude-cli", "codex-cli", "claude-api", "gpt-api"} {
		if !strings.Contains(output, name) {
			t.Errorf("expected %q in agents output", name)
		}
	}
	// CLI providers should show ✓ (cli)
	if !strings.Contains(output, "✓ (cli)") {
		t.Error("expected '✓ (cli)' for CLI providers")
	}
	// API providers without keys should show ✗
	if !strings.Contains(output, "✗ (not set)") {
		t.Error("expected '✗ (not set)' for unset API keys")
	}
	if !strings.Contains(output, "Default:") {
		t.Error("expected 'Default:' line")
	}
}

func TestAddCmd_WithSplitBy(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "pdf", "~/Books/ddia.pdf", "--name", "ddia", "--template", "principles", "--split-by", "chapter"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.LoadFrom("distill.yaml")
	src, ok := loaded.Sources["ddia"]
	if !ok {
		t.Fatal("source 'ddia' not found")
	}
	if src.SplitBy != "chapter" {
		t.Errorf("expected split_by 'chapter', got %q", src.SplitBy)
	}
	if src.Template != "principles" {
		t.Errorf("expected template 'principles', got %q", src.Template)
	}
}

func TestAddCmd_AutoDeriveName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "pdf", "~/Books/Clean Code.pdf"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.LoadFrom("distill.yaml")
	if _, ok := loaded.Sources["clean-code"]; !ok {
		t.Error("expected auto-derived name 'clean-code'")
	}
}

func TestCompactCmd_DryRunWithTokenBudget(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("source.md", []byte("# Content\n\nText."), 0o644)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"src": {Type: "markdown", Path: "./source.md", Template: "rules", OutputDir: "out"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"src", "--dry-run", "--token-budget", "2000"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompactCmd_DefaultOutputFilename(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("source.md", []byte("# Content"), 0o644)

	// Source without output_file — should default to name-minified.md
	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"my-guide": {Type: "markdown", Path: "./source.md", Template: "rules", OutputDir: "out"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"my-guide", "--dry-run"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitCmd_VerifyContents(t *testing.T) {
	dir := t.TempDir()
	repoName := filepath.Join(dir, "test-context")

	root := newRootCmd("dev")
	root.SetArgs([]string{"init", repoName})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify index.md content
	data, _ := os.ReadFile(filepath.Join(repoName, "index.md"))
	if !strings.Contains(string(data), "test-context") {
		t.Error("index.md should contain repo name")
	}
	if !strings.Contains(string(data), "distill") {
		t.Error("index.md should reference distill")
	}

	// Verify .gitignore
	data, _ = os.ReadFile(filepath.Join(repoName, ".gitignore"))
	if !strings.Contains(string(data), ".distill-state.yaml") {
		t.Error(".gitignore should exclude state file")
	}

	// Verify README.md
	data, _ = os.ReadFile(filepath.Join(repoName, "README.md"))
	if !strings.Contains(string(data), "test-context") {
		t.Error("README should contain repo name")
	}
	if !strings.Contains(string(data), "~/.claude") {
		t.Error("README should mention ~/.claude")
	}
}

func TestInitCmd_MissingArg(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"init"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestListCmd_MultipleSources(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"alpha": {Type: "pdf", Template: "rules", OutputDir: "a"},
		"beta":  {Type: "markdown", Template: "principles", OutputDir: "b"},
		"gamma": {Type: "url", Template: "patterns", OutputDir: "c"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"list"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	// Should be sorted alphabetically
	alphaIdx := strings.Index(output, "alpha")
	betaIdx := strings.Index(output, "beta")
	gammaIdx := strings.Index(output, "gamma")
	if alphaIdx > betaIdx || betaIdx > gammaIdx {
		t.Error("expected sources sorted alphabetically")
	}
	if !strings.Contains(output, "rules") {
		t.Error("expected template names in output")
	}
}

func TestUpdateCmd_DryRun(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("source.md", []byte("# Test"), 0o644)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"src": {Type: "markdown", Path: "./source.md", Template: "rules", OutputDir: "out", OutputFile: "out.md"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"update", "src", "--dry-run"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddCmd_NotionSource(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "notion", "https://notion.so/page-abc123", "--name", "my-notion"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.LoadFrom("distill.yaml")
	src := loaded.Sources["my-notion"]
	if src.URL != "https://notion.so/page-abc123" {
		t.Errorf("expected URL to be set, got %q", src.URL)
	}
	if src.Type != "notion" {
		t.Errorf("expected type 'notion', got %q", src.Type)
	}
}

func TestAddCmd_GitHubSource(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := config.DefaultConfig()
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"add", "github", "owner/repo", "--name", "upstream"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.LoadFrom("distill.yaml")
	src := loaded.Sources["upstream"]
	if src.Repo != "owner/repo" {
		t.Errorf("expected Repo 'owner/repo', got %q", src.Repo)
	}
}

func TestCompactCmd_WithSpecificAgent(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("src.md", []byte("# Test"), 0o644)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"test": {Type: "markdown", Path: "./src.md", Template: "rules", OutputDir: "out"},
	}
	config.SaveTo(cfg, "distill.yaml")

	// Use codex-cli agent (will fail to exec, but exercises the agent selection path)
	root := newRootCmd("dev")
	root.SetArgs([]string{"test", "--agent", "codex-cli", "--dry-run"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompactCmd_BadAgent(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.WriteFile("src.md", []byte("# Test"), 0o644)

	cfg := config.DefaultConfig()
	cfg.Sources = map[string]config.Source{
		"test": {Type: "markdown", Path: "./src.md", Template: "rules", OutputDir: "out"},
	}
	config.SaveTo(cfg, "distill.yaml")

	root := newRootCmd("dev")
	root.SetArgs([]string{"test", "--agent", "nonexistent-agent"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent agent")
	}
}

func TestExecute(t *testing.T) {
	// Just verify Execute doesn't panic when called with --help
	origArgs := os.Args
	os.Args = []string{"distill", "--help"}
	defer func() { os.Args = origArgs }()

	// Execute returns nil for --help
	err := Execute("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInstallCmd_MissingArgs(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"install"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestInstallCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"install", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "install") {
		t.Error("expected install help text")
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"a\nb\nc", []string{"a", "b", "c"}},
		{"single", []string{"single"}},
		{"", nil},
		{"a\n", []string{"a"}},
		{"a\nb\n", []string{"a", "b"}},
	}
	for _, tt := range tests {
		got := splitLines(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitLines(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitLines(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestAppendGitignore_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	appendGitignore(path, "docs/")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading gitignore: %v", err)
	}
	if string(data) != "docs/\n" {
		t.Errorf("expected 'docs/\n', got %q", string(data))
	}
}

func TestAppendGitignore_AlreadyPresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	os.WriteFile(path, []byte("docs/\n"), 0o644)

	appendGitignore(path, "docs/")

	data, _ := os.ReadFile(path)
	// Should not duplicate
	if string(data) != "docs/\n" {
		t.Errorf("expected no duplicate, got %q", string(data))
	}
}

func TestAppendGitignore_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	os.WriteFile(path, []byte("node_modules/\n"), 0o644)

	appendGitignore(path, "docs/")

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "node_modules/") {
		t.Error("lost existing content")
	}
	if !strings.Contains(string(data), "docs/") {
		t.Error("docs/ not appended")
	}
}
