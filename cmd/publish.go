package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPublishCmd() *cobra.Command {
	var (
		repo   string
		branch string
		push   bool
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Copy compacted output into a context repo and commit",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = repo
			_ = branch
			_ = push
			fmt.Fprintln(cmd.OutOrStdout(), "publish: not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringVar(&repo, "repo", "", "path to target context repo")
	cmd.Flags().StringVar(&branch, "branch", "", "branch to commit to")
	cmd.Flags().BoolVar(&push, "push", false, "push after committing")

	return cmd
}
