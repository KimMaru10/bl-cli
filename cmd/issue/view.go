package issue

import "github.com/spf13/cobra"

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [issueKey]",
		Short: "課題の詳細を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement
			return nil
		},
	}
}
