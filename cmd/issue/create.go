package issue

import "github.com/spf13/cobra"

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "課題を作成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement
			return nil
		},
	}
}
