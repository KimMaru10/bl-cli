package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/KimMaru10/bl-cli/internal/api"
)

// StripControl removes terminal control characters (ESC, CR, C0/C1, DEL)
// from server-supplied text, keeping newlines and tabs.
func StripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// SaveFile writes data into dir using only a sanitized base name of name,
// so a server-supplied name cannot escape dir. It never overwrites or follows
// an existing path: on collision it appends " (n)". fallback is used when
// name is unusable.
func SaveFile(dir, name, fallback string, data []byte) (string, error) {
	clean := strings.ReplaceAll(StripControl(name), "\n", "")
	clean = strings.ReplaceAll(clean, "\t", "")
	base := filepath.Base(filepath.Clean("/" + strings.ReplaceAll(clean, "\\", "/")))
	if base == "/" || base == "." || base == ".." || base == "" {
		base = fallback
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("保存先フォルダの作成に失敗しました: %w", err)
	}

	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for i := 0; i < 1000; i++ {
		candidate := base
		if i > 0 {
			candidate = fmt.Sprintf("%s (%d)%s", stem, i, ext)
		}
		path := filepath.Join(dir, candidate)
		// O_EXCL fails on any existing path, including symlinks.
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("ファイルの保存に失敗しました: %w", err)
		}
		_, werr := f.Write(data)
		if cerr := f.Close(); werr == nil {
			werr = cerr
		}
		if werr != nil {
			os.Remove(path)
			return "", fmt.Errorf("ファイルの保存に失敗しました: %w", werr)
		}
		return path, nil
	}
	return "", fmt.Errorf("同名のファイルが多すぎて保存できません: %s", base)
}

// DownloadAttachments saves the attachments (all of them, or only id when non-zero)
// using dl, and prints each saved path.
func DownloadAttachments(attachments []api.Attachment, id int, output string, dl func(int) (string, []byte, error)) error {
	found := false
	for _, a := range attachments {
		if id != 0 && a.ID != id {
			continue
		}
		found = true
		name, data, err := dl(a.ID)
		if err != nil {
			return err
		}
		if name == "" {
			name = a.Name
		}
		path, err := SaveFile(output, name, fmt.Sprintf("attachment-%d", a.ID), data)
		if err != nil {
			return err
		}
		fmt.Println(path)
	}
	if !found {
		return fmt.Errorf("添付ファイル ID %d はありません", id)
	}
	return nil
}
