package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/config"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [name]",
		Short: "Re-compact one or all tracked sources",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if len(cfg.Sources) == 0 {
				fmt.Println("No sources configured.")
				return nil
			}

			// If a specific name is given, compact just that one
			if len(args) == 1 {
				return compactSource(args[0])
			}

			// Otherwise, compact all
			names := make([]string, 0, len(cfg.Sources))
			for n := range cfg.Sources {
				names = append(names, n)
			}
			sort.Strings(names)

			fmt.Printf("→ Updating %d sources...\n", len(names))
			updated := 0
			for _, name := range names {
				if err := compactSource(name); err != nil {
					fmt.Printf("⚠ %s: %v\n", name, err)
					continue
				}
				updated++
			}

			fmt.Printf("\n✓ %d source(s) updated.\n", updated)
			return nil
		},
	}
}
