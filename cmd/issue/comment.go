package issue

import "github.com/spf13/cobra"

func newCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment [issueKey]",
		Short: "課題にコメントを追加する",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement
			return nil
		},
	}

	return cmd
}
