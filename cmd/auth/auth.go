package auth

import "github.com/spf13/cobra"

// NewAuthCmd returns the auth subcommand group.
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "認証の管理",
	}

	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newStatusCmd())

	return cmd
}
