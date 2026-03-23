package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"go.yaml.in/yaml/v3"
)

// SpaceConfig holds credentials and settings for a single Backlog space.
type SpaceConfig struct {
	SpaceURL       string `yaml:"space_url"`
	APIKey         string `yaml:"api_key"`
	DefaultProject string `yaml:"default_project,omitempty"`
}

// Config holds the application configuration supporting multiple spaces.
type Config struct {
	CurrentSpace string                 `yaml:"current_space"`
	Spaces       map[string]SpaceConfig `yaml:"spaces"`
}

// Current returns the currently active SpaceConfig, or nil if not set.
func (c *Config) Current() *SpaceConfig {
	if c.CurrentSpace == "" {
		return nil
	}
	s, ok := c.Spaces[c.CurrentSpace]
	if !ok {
		return nil
	}
	return &s
}

// SpaceNames returns sorted space alias names.
func (c *Config) SpaceNames() []string {
	names := make([]string, 0, len(c.Spaces))
	for name := range c.Spaces {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
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
// Automatically migrates from the old single-space format.
func Load() (*Config, error) {
	p, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Spaces: make(map[string]SpaceConfig)}, nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	// Try new format first
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗しました: %w", err)
	}

	// If spaces map is populated, it's the new format
	if len(cfg.Spaces) > 0 {
		return &cfg, nil
	}

	// Try old single-space format for migration
	var old struct {
		SpaceURL       string `yaml:"space_url"`
		APIKey         string `yaml:"api_key"`
		DefaultProject string `yaml:"default_project"`
	}
	if err := yaml.Unmarshal(data, &old); err == nil && old.SpaceURL != "" {
		alias := ExtractAlias(old.SpaceURL)
		migrated := &Config{
			CurrentSpace: alias,
			Spaces: map[string]SpaceConfig{
				alias: {
					SpaceURL:       old.SpaceURL,
					APIKey:         old.APIKey,
					DefaultProject: old.DefaultProject,
				},
			},
		}
		// Save migrated config
		_ = Save(migrated)
		return migrated, nil
	}

	return &Config{Spaces: make(map[string]SpaceConfig)}, nil
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

	p, err := configPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %w", err)
	}
	return os.WriteFile(p, data, 0600)
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

// ExtractAlias derives an alias from a Backlog space URL.
// e.g., "https://myteam.backlog.com" -> "myteam"
func ExtractAlias(spaceURL string) string {
	// Remove scheme
	u := spaceURL
	for _, prefix := range []string{"https://", "http://"} {
		if len(u) > len(prefix) && u[:len(prefix)] == prefix {
			u = u[len(prefix):]
			break
		}
	}
	// Take subdomain
	for i, c := range u {
		if c == '.' {
			return u[:i]
		}
	}
	return u
}
