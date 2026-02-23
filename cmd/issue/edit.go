package issue

import "github.com/spf13/cobra"

func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit [issueKey]",
		Short: "課題を更新する",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement
			return nil
		},
	}
}
