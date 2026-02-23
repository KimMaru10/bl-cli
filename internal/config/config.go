package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	SpaceURL       string `mapstructure:"space_url"`
	APIKey         string `mapstructure:"api_key"`
	DefaultProject string `mapstructure:"default_project"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
	}
	return filepath.Join(home, ".config", "bl"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads the config file and returns a Config.
func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗しました: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to the config file.
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %w", err)
	}

	v := viper.New()
	v.Set("space_url", cfg.SpaceURL)
	v.Set("api_key", cfg.APIKey)
	v.Set("default_project", cfg.DefaultProject)

	p, err := configPath()
	if err != nil {
		return err
	}
	return v.WriteConfigAs(p)
}

// Delete removes the config file.
func Delete() error {
	p, err := configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("設定ファイルの削除に失敗しました: %w", err)
	}
	return nil
}
