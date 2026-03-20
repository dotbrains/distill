package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagAgent       string
	flagForce       bool
	flagDryRun      bool
	flagVerbose     bool
	flagTokenBudget int
)

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "distill <name>",
		Short: "AI-powered knowledge compactor for agents",
		Long:  "Compact technical books, documentation, and reference material into agent-optimized markdown.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runCompact,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Version: version,
	}

	root.SetVersionTemplate(fmt.Sprintf("distill version %s\n", version))

	// Global flags
	root.PersistentFlags().StringVar(&flagAgent, "agent", "", "use a specific configured agent")
	root.PersistentFlags().BoolVar(&flagForce, "force", false, "force re-compaction even if source unchanged")
	root.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "show what would happen without writing files")
	root.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "show detailed progress")
	root.PersistentFlags().IntVar(&flagTokenBudget, "token-budget", 0, "override token budget for this run")

	// Subcommands
	root.AddCommand(newAgentsCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newTemplatesCmd())
	root.AddCommand(newAddCmd())
	root.AddCommand(newListCmd())
	root.AddCommand(newUpdateCmd())
	root.AddCommand(newValidateCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newPublishCmd())

	return root
}

// Execute runs the root command.
func Execute(version string) error {
	return newRootCmd(version).Execute()
}
