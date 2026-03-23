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

			if len(cfg.Spaces) == 0 {
				fmt.Println(errorStyle.Render("✗ 未認証です。bl auth login を実行してください"))
				return nil
			}

			names := cfg.SpaceNames()
			for _, name := range names {
				space := cfg.Spaces[name]
				marker := "  "
				if name == cfg.CurrentSpace {
					marker = "* "
				}

				client := api.NewClient(space.SpaceURL, space.APIKey)
				user, err := client.GetMyself()
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("%s%s (%s) - ✗ 認証が無効です", marker, name, space.SpaceURL)))
					continue
				}

				line := fmt.Sprintf("%s%s (%s) - ✔ %s としてログイン中", marker, name, space.SpaceURL, user.Name)
				if name == cfg.CurrentSpace {
					fmt.Println(successStyle.Render(line))
				} else {
					fmt.Println(line)
				}
			}

			return nil
		},
	}
}
