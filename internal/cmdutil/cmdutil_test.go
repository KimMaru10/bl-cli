package cmdutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
)

func TestSaveFileStaysInsideDir(t *testing.T) {
	cases := map[string]string{
		"../../evil.txt":    "evil.txt",
		`..\..\win.txt`:     "win.txt",
		"/abs/path/a.png":   "a.png",
		"":                  "fallback.bin",
		"..":                "fallback.bin",
		"スクリーンショット.png":     "スクリーンショット.png",
		"a\x1b[2Jb\r\n.txt": "a[2Jb.txt",
	}
	for name, want := range cases {
		dir := t.TempDir()
		path, err := SaveFile(dir, name, "fallback.bin", []byte("x"))
		if err != nil {
			t.Fatalf("%q: %v", name, err)
		}
		if path != filepath.Join(dir, want) {
			t.Errorf("%q: saved to %s, want %s", name, path, filepath.Join(dir, want))
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("%q: file not written: %v", name, err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Errorf("%q: perm = %o, want 600", name, info.Mode().Perm())
		}
	}
}

func TestSaveFileNeverOverwrites(t *testing.T) {
	dir := t.TempDir()
	first, _ := SaveFile(dir, "a.png", "f", []byte("1"))
	second, err := SaveFile(dir, "a.png", "f", []byte("2"))
	if err != nil {
		t.Fatal(err)
	}
	if second != filepath.Join(dir, "a (1).png") {
		t.Fatalf("second saved to %s", second)
	}
	if b, _ := os.ReadFile(first); string(b) != "1" {
		t.Fatalf("first file was overwritten: %q", b)
	}
}

func TestSaveFileDoesNotFollowSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(t.TempDir(), "secret")
	os.WriteFile(target, []byte("keep"), 0o600)
	if err := os.Symlink(target, filepath.Join(dir, "report.pdf")); err != nil {
		t.Skip("symlink not supported")
	}
	path, err := SaveFile(dir, "report.pdf", "f", []byte("attack"))
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(target); string(b) != "keep" {
		t.Fatalf("symlink target was written: %q", b)
	}
	if path != filepath.Join(dir, "report (1).pdf") {
		t.Fatalf("saved to %s", path)
	}
}

func TestStripControl(t *testing.T) {
	got := StripControl("a\x1b]52;c;evil\x07b\r\nc\td\u0085")
	if got != "a]52;c;evilb\nc\td" {
		t.Fatalf("got %q", got)
	}
}

func TestDownloadAttachmentsUnknownID(t *testing.T) {
	err := DownloadAttachments([]api.Attachment{{ID: 1, Name: "a"}}, 99, t.TempDir(), func(int) (string, []byte, error) {
		t.Fatal("should not download")
		return "", nil, nil
	})
	if err == nil {
		t.Fatal("want error for unknown id")
	}
}

func TestResolveProjectKey(t *testing.T) {
	cfg := &config.Config{
		CurrentSpace: "a",
		Spaces:       map[string]config.SpaceConfig{"a": {DefaultProject: "DEF"}},
	}
	if k, _ := ResolveProjectKey(cfg, "FLAG"); k != "FLAG" {
		t.Errorf("flag: got %q", k)
	}
	if k, _ := ResolveProjectKey(cfg, ""); k != "DEF" {
		t.Errorf("default: got %q", k)
	}
	if _, err := ResolveProjectKey(&config.Config{}, ""); err == nil {
		t.Error("no space: want error")
	}
}
