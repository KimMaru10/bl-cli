package project

import "github.com/spf13/cobra"

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "デフォルトプロジェクトを設定する",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement
			return nil
		},
	}
}
