package auth

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "認証状態を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("設定の読み込みに失敗しました: %w", err)
			}

			if cfg.APIKey == "" {
				fmt.Println(errorStyle.Render("✗ 未認証です。bl auth login を実行してください"))
				return nil
			}

			client := api.NewClient(cfg.SpaceURL, cfg.APIKey)
			user, err := client.GetMyself()
			if err != nil {
				fmt.Println(errorStyle.Render("✗ 認証が無効です。bl auth login を再実行してください"))
				return nil
			}

			fmt.Println(successStyle.Render("✔ " + cfg.SpaceURL + " に " + user.Name + " としてログイン中"))
			return nil
		},
	}
}
