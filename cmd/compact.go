package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/agent"
	_ "github.com/dotbrains/distill/internal/agent/anthropic" // register provider
	_ "github.com/dotbrains/distill/internal/agent/claudecli" // register provider
	_ "github.com/dotbrains/distill/internal/agent/codexcli"  // register provider
	_ "github.com/dotbrains/distill/internal/agent/openai"    // register provider
	"github.com/dotbrains/distill/internal/config"
	"github.com/dotbrains/distill/internal/ingest"
	"github.com/dotbrains/distill/internal/state"
	"github.com/dotbrains/distill/internal/template"
	"github.com/dotbrains/distill/internal/writer"
)

func runCompact(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return compactSource(args[0])
}

func compactSource(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	src, ok := cfg.Sources[name]
	if !ok {
		return fmt.Errorf("source %q not found in config (available: %s)", name, availableSources(cfg))
	}

	// Load template
	tmpl, err := template.Load(src.Template, cfg.Output.CustomTemplates)
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}

	// Create agent
	a, err := agent.NewAgentFromConfig(flagAgent, cfg)
	if err != nil {
		return err
	}

	// Ingest
	ingestor, err := ingest.New(src.Type)
	if err != nil {
		return err
	}

	fmt.Printf("→ source:   %s (%s)\n", name, src.Type)
	fmt.Printf("→ template: %s\n", src.Template)
	fmt.Printf("→ agent:    %s\n", a.Name())

	chunks, err := ingestor.Ingest(src)
	if err != nil {
		return fmt.Errorf("ingesting source: %w", err)
	}
	fmt.Printf("→ Compacting (%d chunk(s))...\n", len(chunks))

	tokenBudget := cfg.Output.TokenBudget
	if flagTokenBudget > 0 {
		tokenBudget = flagTokenBudget
	}

	// Compact each chunk
	var outputs []string
	for i, chunk := range chunks {
		input := &agent.CompactInput{
			SourceName:  name,
			ChunkIndex:  i,
			TotalChunks: len(chunks),
			Content:     chunk.Content,
			Template:    tmpl + fmt.Sprintf("\n\nToken budget: ~%d tokens.", tokenBudget),
			TokenBudget: tokenBudget,
			Metadata: agent.SourceMetadata{
				Type:         src.Type,
				ChapterNum:   chunk.ChapterNum,
				ChapterTitle: chunk.ChapterTitle,
			},
		}

		if flagDryRun {
			fmt.Printf("→ [dry-run] Would compact chunk %d/%d (%d chars)\n", i+1, len(chunks), len(chunk.Content))
			continue
		}

		result, err := a.Compact(context.Background(), input)
		if err != nil {
			return fmt.Errorf("compacting chunk %d: %w", i+1, err)
		}
		outputs = append(outputs, result.Content)
	}

	if flagDryRun {
		fmt.Println("✓ Dry run complete.")
		return nil
	}

	// Assemble output
	assembled := strings.Join(outputs, "\n\n")

	// Determine output filename
	filename := src.OutputFile
	if filename == "" {
		filename = name + "-minified.md"
	}

	// Write
	outputPath, err := writer.WriteOutput(cfg.Output.Dir, src.OutputDir, filename, assembled)
	if err != nil {
		return err
	}

	// Update state
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	// Build content hash from all chunks
	var allContent string
	for _, c := range chunks {
		allContent += c.Content
	}

	st.Update(name, state.SourceState{
		ContentHash:   state.HashContent(allContent),
		TemplateHash:  state.HashContent(tmpl),
		Agent:         a.Name(),
		LastCompacted: time.Now().UTC(),
		OutputFiles:   []string{outputPath},
	})
	if err := st.Save(); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Compaction complete.")
	fmt.Printf("→ output:  %s\n", outputPath)
	return nil
}

func availableSources(cfg *config.Config) string {
	if len(cfg.Sources) == 0 {
		return "none (add sources to distill.yaml)"
	}
	names := make([]string, 0, len(cfg.Sources))
	for n := range cfg.Sources {
		names = append(names, n)
	}
	return strings.Join(names, ", ")
}
