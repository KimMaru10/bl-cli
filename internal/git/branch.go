package git

import (
	"os/exec"
	"regexp"
	"strings"
)

var issueKeyPattern = regexp.MustCompile(`[A-Z][A-Z0-9]+-\d+`)

// GetCurrentBranch returns the current git branch name.
func GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ExtractIssueKey extracts a Backlog issue key from a branch name.
// Returns an empty string if no match is found.
func ExtractIssueKey(branch string) string {
	match := issueKeyPattern.FindString(branch)
	return match
}
