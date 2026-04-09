package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Setup registers the bl mcp server in Claude Desktop's config file.
func Setup() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("現在 macOS のみ対応しています")
	}

	// Find bl binary path
	blPath, err := exec.LookPath("bl")
	if err != nil {
		return fmt.Errorf("bl コマンドが見つかりません。先に npm install -g @kimmaru10/bl-cli を実行してください")
	}
	// Resolve to absolute path (keep symlink target for portability)
	blPath, err = filepath.Abs(blPath)
	if err != nil {
		return fmt.Errorf("bl のパス解決に失敗しました: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
	}

	configPath := filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")

	// Read existing config or create empty
	var configMap map[string]any

	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
		}
		configMap = map[string]any{}
	} else {
		if err := json.Unmarshal(data, &configMap); err != nil {
			return fmt.Errorf("設定ファイルの解析に失敗しました: %w", err)
		}
	}

	// Check if already registered
	if servers, ok := configMap["mcpServers"].(map[string]any); ok {
		if _, exists := servers["backlog"]; exists {
			fmt.Println("✔ backlog MCP サーバーは既に登録されています")
			fmt.Printf("  パス: %s\n", blPath)
			return nil
		}
	}

	// Add mcpServers.backlog
	servers, ok := configMap["mcpServers"].(map[string]any)
	if !ok {
		servers = map[string]any{}
	}
	servers["backlog"] = map[string]any{
		"command": blPath,
		"args":    []string{"mcp"},
	}
	configMap["mcpServers"] = servers

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
	}

	// Write config
	out, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %w", err)
	}
	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗しました: %w", err)
	}

	fmt.Println("✔ Claude Desktop に backlog MCP サーバーを登録しました")
	fmt.Printf("  パス: %s\n", blPath)
	fmt.Println()
	fmt.Println("Claude Desktop を再起動すると、Backlog ツールが使えるようになります")

	return nil
}
