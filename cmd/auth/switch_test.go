package auth

import (
	"testing"

	"github.com/KimMaru10/bl-cli/internal/config"
)

func TestSwitchToAlias(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg := &config.Config{
		CurrentSpace: "a",
		Spaces: map[string]config.SpaceConfig{
			"a": {SpaceURL: "https://a.backlog.com", APIKey: "key-a"},
			"b": {SpaceURL: "https://b.backlog.com", APIKey: "key-b"},
		},
	}

	if err := switchToAlias(cfg, "b"); err != nil {
		t.Fatalf("switchToAlias(b) unexpected error: %v", err)
	}
	if cfg.CurrentSpace != "b" {
		t.Errorf("CurrentSpace = %q, want %q", cfg.CurrentSpace, "b")
	}

	reloaded, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() error: %v", err)
	}
	if reloaded.CurrentSpace != "b" {
		t.Errorf("persisted CurrentSpace = %q, want %q", reloaded.CurrentSpace, "b")
	}
	if reloaded.Spaces["a"] != cfg.Spaces["a"] {
		t.Errorf("space 'a' was modified unexpectedly: %+v", reloaded.Spaces["a"])
	}
	if reloaded.Spaces["b"] != cfg.Spaces["b"] {
		t.Errorf("space 'b' was modified unexpectedly: %+v", reloaded.Spaces["b"])
	}

	if err := switchToAlias(cfg, "nonexistent"); err == nil {
		t.Error("switchToAlias(nonexistent) expected error, got nil")
	}
}

func TestSwitchCmdNoSpacesWithAlias(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cmd := newSwitchCmd()
	if err := cmd.RunE(cmd, []string{"myteam"}); err == nil {
		t.Error("expected error when no spaces are registered, got nil (exit code must be non-zero for script/AI callers)")
	}
}
