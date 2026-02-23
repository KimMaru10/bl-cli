package auth

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "認証情報を削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Delete(); err != nil {
				return fmt.Errorf("認証情報の削除に失敗しました: %w", err)
			}
			fmt.Println(successStyle.Render("✔ 認証情報を削除しました"))
			return nil
		},
	}
}
