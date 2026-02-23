package cmd

import (
	"github.com/KimMaru10/bl-cli/cmd/auth"
	"github.com/KimMaru10/bl-cli/cmd/issue"
	"github.com/KimMaru10/bl-cli/cmd/project"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bl",
	Short: "Backlog CLI - ターミナルから Backlog を操作",
}

func init() {
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(project.NewProjectCmd())
	rootCmd.AddCommand(issue.NewIssueCmd())
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
