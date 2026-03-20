package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed prompts/*.md
var builtinFS embed.FS

// BuiltinTemplates lists the names of all built-in templates.
var BuiltinTemplates = []struct {
	Name        string
	Description string
}{
	{"rules", "Numbered imperative rules grouped by section"},
	{"principles", "Chapter-based core principles with loading guidance"},
	{"patterns", "Named patterns with rationale and code examples"},
	{"raw", "Minimal compaction — condense prose, keep structure, no reformatting"},
}

// Load returns the system prompt for a template by name.
// Checks custom templates directory first, then built-in.
func Load(name string, customDir string) (string, error) {
	// Try custom directory first
	if customDir != "" {
		path := filepath.Join(customDir, name+".md")
		data, err := os.ReadFile(path)
		if err == nil {
			return stripFrontmatter(string(data)), nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("reading custom template %s: %w", path, err)
		}
	}

	// Try built-in
	data, err := builtinFS.ReadFile("prompts/" + name + ".md")
	if err != nil {
		return "", fmt.Errorf("unknown template %q (available: %s)", name, availableNames())
	}
	return string(data), nil
}

// IsBuiltin returns true if the name matches a built-in template.
func IsBuiltin(name string) bool {
	for _, t := range BuiltinTemplates {
		if t.Name == name {
			return true
		}
	}
	return false
}

func availableNames() string {
	names := make([]string, len(BuiltinTemplates))
	for i, t := range BuiltinTemplates {
		names[i] = t.Name
	}
	return strings.Join(names, ", ")
}

// stripFrontmatter removes YAML frontmatter (--- ... ---) from the top of a file.
func stripFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---") {
		return content
	}
	lines := strings.Split(content, "\n")
	count := 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			count++
			if count == 2 {
				rest := strings.Join(lines[i+1:], "\n")
				return strings.TrimLeft(rest, "\n")
			}
		}
	}
	return content
}
