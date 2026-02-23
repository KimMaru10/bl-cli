package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open opens the given URL in the default browser.
func Open(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		return fmt.Errorf("ブラウザの起動に対応していないOS: %s", runtime.GOOS)
	}
	return exec.Command(cmd, url).Start()
}
