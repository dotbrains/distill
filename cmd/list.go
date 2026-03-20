package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/config"
	"github.com/dotbrains/distill/internal/state"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tracked sources with their status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if len(cfg.Sources) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No sources configured. Run `distill add` to add one.")
				return nil
			}

			st, err := state.Load()
			if err != nil {
				return fmt.Errorf("loading state: %w", err)
			}

			// Sort by name
			names := make([]string, 0, len(cfg.Sources))
			for n := range cfg.Sources {
				names = append(names, n)
			}
			sort.Strings(names)

			for _, name := range names {
				src := cfg.Sources[name]
				entry, compacted := st.Sources[name]

				status := "✗ not yet compacted"
				tokens := "-"
				if compacted {
					tokens = fmt.Sprintf("%d tok", entry.OutputTokens)
					status = "✓ current"
				}

				fmt.Fprintf(cmd.OutOrStdout(), "  %-18s %-10s %-12s %-14s %-10s %s\n",
					name, src.Type, src.Template, src.OutputDir+"/", tokens, status)
			}
			return nil
		},
	}
}
