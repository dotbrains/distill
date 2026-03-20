package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/exec"
)

func newInstallCmd() *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "install <repo-url>",
		Short: "Clone a context repo into ~/.claude/docs/ for agent consumption",
		Long: `Install a shared context repository so AI agents can discover and load
the compacted documents. Clones the repo into the target directory
(default: ~/.claude/docs/).

If the target directory already exists as a git repo, pulls the latest
changes instead of re-cloning.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoURL := args[0]

			if target == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("unable to determine home directory: %w", err)
				}
				target = filepath.Join(home, ".claude", "docs")
			}

			executor := exec.NewRealExecutor()
			ctx := cmd.Context()

			// Check if target already exists and is a git repo
			gitDir := filepath.Join(target, ".git")
			if _, err := os.Stat(gitDir); err == nil {
				// Already cloned — pull latest
				fmt.Printf("→ %s already exists, pulling latest...\n", target)
				out, err := executor.Run(ctx, "git", "-C", target, "pull")
				if err != nil {
					return fmt.Errorf("git pull failed: %w", err)
				}
				fmt.Print(out)
				fmt.Println("✓ Context repo updated.")
				return nil
			}

			// Ensure parent directory exists
			parent := filepath.Dir(target)
			if err := os.MkdirAll(parent, 0o755); err != nil {
				return fmt.Errorf("creating parent directory: %w", err)
			}

			// Add docs/ to ~/.claude/.gitignore if ~/.claude is a git repo
			claudeGitignore := filepath.Join(parent, ".gitignore")
			if _, err := os.Stat(filepath.Join(parent, ".git")); err == nil {
				appendGitignore(claudeGitignore, "docs/")
			}

			// Clone
			fmt.Printf("→ Cloning %s into %s...\n", repoURL, target)
			out, err := executor.Run(ctx, "git", "clone", repoURL, target)
			if err != nil {
				return fmt.Errorf("git clone failed: %w", err)
			}
			fmt.Print(out)
			fmt.Printf("✓ Context repo installed at %s\n", target)
			fmt.Println("  Agents can now load documents from this directory.")
			return nil
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "target directory (default: ~/.claude/docs/)")

	return cmd
}

// appendGitignore adds an entry to a .gitignore file if it's not already present.
func appendGitignore(path, entry string) {
	data, err := os.ReadFile(path)
	if err == nil {
		content := string(data)
		for _, line := range filepath.SplitList(content) {
			if line == entry {
				return // already present
			}
		}
		// Check line by line
		lines := splitLines(content)
		for _, line := range lines {
			if line == entry {
				return
			}
		}
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return // non-fatal
	}
	defer f.Close()
	_, _ = f.WriteString(entry + "\n")
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
