package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/config"
)

func newAddCmd() *cobra.Command {
	var (
		name      string
		tmpl      string
		outputDir string
		splitBy   string
	)

	cmd := &cobra.Command{
		Use:   "add <type> <location>",
		Short: "Register a new source for tracking",
		Long:  "Add a source (pdf, markdown, notion, url, epub, github) to distill.yaml.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceType := args[0]
			location := args[1]

			// Derive name from location if not provided
			if name == "" {
				name = deriveName(sourceType, location)
			}
			if tmpl == "" {
				tmpl = "rules"
			}
			if outputDir == "" {
				outputDir = name
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if _, exists := cfg.Sources[name]; exists {
				return fmt.Errorf("source %q already exists in config", name)
			}

			src := config.Source{
				Type:      sourceType,
				Template:  tmpl,
				OutputDir: outputDir,
				SplitBy:   splitBy,
			}

			// Set path or URL based on type
			switch sourceType {
			case "pdf", "markdown", "epub":
				src.Path = location
				src.OutputFile = name + "-minified.md"
			case "notion", "url":
				src.URL = location
				src.OutputFile = name + "-minified.md"
			case "github":
				src.Repo = location
				src.OutputFile = name + "-minified.md"
			default:
				return fmt.Errorf("unknown source type %q (supported: pdf, markdown, notion, url, epub, github)", sourceType)
			}

			if cfg.Sources == nil {
				cfg.Sources = map[string]config.Source{}
			}
			cfg.Sources[name] = src

			path := config.ProjectConfigPath()
			if err := config.SaveTo(cfg, path); err != nil {
				return err
			}

			detail := sourceType
			if splitBy != "" {
				detail += ", split by " + splitBy
			}
			fmt.Printf("✓ Added source %q (%s)\n", name, detail)
			fmt.Printf("Run `distill %s` to compact it.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "source name (derived from location if omitted)")
	cmd.Flags().StringVar(&tmpl, "template", "", "compaction template (default: rules)")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "output subdirectory (default: source name)")
	cmd.Flags().StringVar(&splitBy, "split-by", "", "split strategy (e.g., chapter)")

	return cmd
}

func deriveName(sourceType, location string) string {
	switch sourceType {
	case "pdf", "markdown", "epub":
		base := filepath.Base(location)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	case "url", "notion":
		// Extract last path segment
		parts := strings.Split(strings.TrimRight(location, "/"), "/")
		if len(parts) > 0 {
			slug := parts[len(parts)-1]
			// Remove Notion IDs (hex suffix after last dash)
			if idx := strings.LastIndex(slug, "-"); idx > 0 && len(slug)-idx > 16 {
				slug = slug[:idx]
			}
			return strings.ToLower(slug)
		}
		return "source"
	default:
		return "source"
	}
}
