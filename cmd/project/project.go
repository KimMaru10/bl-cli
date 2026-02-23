package project

import "github.com/spf13/cobra"

// NewProjectCmd returns the project subcommand group.
func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "プロジェクトの管理",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newCurrentCmd())

	return cmd
}
