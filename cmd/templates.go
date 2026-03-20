package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/template"
)

func newTemplatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List available compaction templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, t := range template.BuiltinTemplates {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-12s %s\n", t.Name, t.Description)
			}
			return nil
		},
	}
}
