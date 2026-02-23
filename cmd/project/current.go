package project

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/spf13/cobra"
)

func newCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "デフォルトプロジェクトを表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("設定の読み込みに失敗しました: %w", err)
			}

			if cfg.DefaultProject == "" {
				fmt.Println("デフォルトプロジェクトが未設定です。bl project set を実行してください")
				return nil
			}

			fmt.Println(cfg.DefaultProject)
			return nil
		},
	}
}
