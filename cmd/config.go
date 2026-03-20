package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotbrains/distill/internal/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}
	cmd.AddCommand(newConfigInitCmd())
	return cmd
}

func newConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Scaffold a default distill.yaml config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.ProjectConfigPath()

			if !flagForce {
				if _, err := os.Stat(path); err == nil {
					return fmt.Errorf("%s already exists (use --force to overwrite)", path)
				}
			}

			cfg := config.DefaultConfig()
			if err := config.SaveTo(cfg, path); err != nil {
				return err
			}

			fmt.Printf("✓ Wrote default config to %s\n", path)
			fmt.Println("Edit the file to add your sources and customize templates.")
			return nil
		},
	}
}
