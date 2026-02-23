package cmdutil

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
)

// LoadConfigAndClient loads the config file and creates an API client.
// Returns an error if not authenticated.
func LoadConfigAndClient() (*config.Config, *api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("設定の読み込みに失敗しました: %w", err)
	}
	if cfg.APIKey == "" {
		return nil, nil, fmt.Errorf("未認証です。bl auth login を先に実行してください")
	}
	client := api.NewClient(cfg.SpaceURL, cfg.APIKey)
	return cfg, client, nil
}
