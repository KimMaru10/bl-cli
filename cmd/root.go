package cmd

import (
	"github.com/KimMaru10/bl-cli/cmd/auth"
	"github.com/KimMaru10/bl-cli/cmd/issue"
	blmcp "github.com/KimMaru10/bl-cli/cmd/mcp"
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
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Claude Desktop 連携（MCP サーバー）",
		Run: func(cmd *cobra.Command, args []string) {
			blmcp.Run()
		},
	}
	mcpCmd.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "Claude Desktop に Backlog MCP サーバーを登録する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return blmcp.Setup()
		},
	})
	rootCmd.AddCommand(mcpCmd)
}

// SetVersion sets the version string shown by --version.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
