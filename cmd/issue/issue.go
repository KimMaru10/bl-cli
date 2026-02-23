package issue

import "github.com/spf13/cobra"

// NewIssueCmd returns the issue subcommand group.
func NewIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "課題の管理",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newViewCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newEditCmd())
	cmd.AddCommand(newCommentCmd())

	return cmd
}
