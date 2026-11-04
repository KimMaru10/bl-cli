package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
)

// StripControl removes terminal control and bidi/format characters (see api.StripControl).
func StripControl(s string) string {
	return api.StripControl(s)
}

// SaveFile writes data into dir using only a sanitized base name of name,
// so a server-supplied name cannot escape dir. It never overwrites or follows
// an existing path: on collision it appends " (n)". fallback is used when
// name is unusable.
func SaveFile(dir, name, fallback string, data []byte) (string, error) {
	clean := strings.ReplaceAll(StripControl(name), "\n", "")
	clean = strings.ReplaceAll(clean, "\t", "")
	base := filepath.Base(filepath.Clean("/" + strings.ReplaceAll(clean, "\\", "/")))
	// Leading dots would create hidden files such as .bashrc.
	base = strings.TrimLeft(base, ".")
	if base == "/" || base == "" {
		base = fallback
	}
	base = truncateName(base, maxNameBytes)
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

// SaveAttachments saves the attachments (all of them, or only id when non-zero)
// into output and returns the saved paths, including those saved before an error.
// It never writes to stdout, which is the protocol channel of the MCP server.
func SaveAttachments(attachments []api.Attachment, id int, output string, dl func(int) (string, []byte, error)) ([]string, error) {
	var paths []string
	found := false
	for _, a := range attachments {
		if id != 0 && a.ID != id {
			continue
		}
		found = true
		name, data, err := dl(a.ID)
		if err != nil {
			return paths, err
		}
		if name == "" {
			name = a.Name
		}
		path, err := SaveFile(output, name, fmt.Sprintf("attachment-%d", a.ID), data)
		if err != nil {
			return paths, err
		}
		paths = append(paths, path)
	}
	if !found {
		return paths, fmt.Errorf("添付ファイル ID %d はありません", id)
	}
	return paths, nil
}

// DownloadAttachments saves the attachments like SaveAttachments and prints the saved paths.
func DownloadAttachments(attachments []api.Attachment, id int, output string, dl func(int) (string, []byte, error)) error {
	paths, err := SaveAttachments(attachments, id, output, dl)
	for _, path := range paths {
		fmt.Println(path)
	}
	return err
}

// maxNameBytes keeps names (plus a " (n)" suffix) under the common 255-byte limit.
const maxNameBytes = 200

// truncateName shortens name to at most max bytes on a rune boundary, keeping the extension.
func truncateName(name string, max int) string {
	if len(name) <= max {
		return name
	}
	ext := filepath.Ext(name)
	if len(ext) > max/2 {
		ext = ""
	}
	stem := strings.TrimSuffix(name, ext)
	limit := max - len(ext)
	cut := 0
	for i := range stem {
		if i > limit {
			break
		}
		cut = i
	}
	return stem[:cut] + ext
}
