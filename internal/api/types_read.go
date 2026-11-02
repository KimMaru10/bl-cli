package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Attachment represents a file attached to an issue or a document.
type Attachment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	CreatedUser *User  `json:"createdUser"`
	Created     string `json:"created"`
}

// CustomField represents a custom field value on an issue.
// Value differs by field type (string, number, object or array), so it is kept raw.
type CustomField struct {
	ID          int             `json:"id"`
	FieldTypeID int             `json:"fieldTypeId"`
	Name        string          `json:"name"`
	Value       json.RawMessage `json:"value"`
}

// DisplayValue renders the custom field value as plain text.
func (f CustomField) DisplayValue() string {
	return StripControl(displayRaw(f.Value))
}

func displayRaw(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var named struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(raw, &named) == nil && named.Name != "" {
		return named.Name
	}
	var list []json.RawMessage
	if json.Unmarshal(raw, &list) == nil {
		var parts []string
		for _, item := range list {
			if v := displayRaw(item); v != "" {
				parts = append(parts, v)
			}
		}
		return strings.Join(parts, ", ")
	}
	// Numbers and anything else are shown as-is.
	return strings.TrimSpace(string(raw))
}

// Tag represents a tag on a document or wiki page.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Document represents a Backlog document.
type Document struct {
	ID          string       `json:"id"`
	ProjectID   int          `json:"projectId"`
	Title       string       `json:"title"`
	Plain       string       `json:"plain"`
	Emoji       string       `json:"emoji"`
	Attachments []Attachment `json:"attachments"`
	Tags        []Tag        `json:"tags"`
	CreatedUser *User        `json:"createdUser"`
	Created     string       `json:"created"`
	UpdatedUser *User        `json:"updatedUser"`
	Updated     string       `json:"updated"`
}

// DocumentNode is a node of the document tree.
type DocumentNode struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Emoji    string         `json:"emoji"`
	Children []DocumentNode `json:"children"`
}

// DocumentTree is the response of the document tree API.
type DocumentTree struct {
	ProjectID  int           `json:"projectId"`
	ActiveTree *DocumentNode `json:"activeTree"`
}

// Wiki represents a Backlog wiki page.
type Wiki struct {
	ID          int          `json:"id"`
	ProjectID   int          `json:"projectId"`
	Name        string       `json:"name"`
	Content     string       `json:"content"`
	Tags        []Tag        `json:"tags"`
	Attachments []Attachment `json:"attachments"`
	CreatedUser *User        `json:"createdUser"`
	Created     string       `json:"created"`
	UpdatedUser *User        `json:"updatedUser"`
	Updated     string       `json:"updated"`
}

// Repository represents a Backlog Git repository.
type Repository struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HTTPURL     string `json:"httpUrl"`
	SSHURL      string `json:"sshUrl"`
	Updated     string `json:"updated"`
}

// PullRequestStatus represents the status of a pull request.
type PullRequestStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PullRequest represents a Backlog Git pull request.
type PullRequest struct {
	ID          int                `json:"id"`
	Number      int                `json:"number"`
	Summary     string             `json:"summary"`
	Description string             `json:"description"`
	Base        string             `json:"base"`
	Branch      string             `json:"branch"`
	Status      *PullRequestStatus `json:"status"`
	Assignee    *User              `json:"assignee"`
	Issue       *struct {
		ID       int    `json:"id"`
		IssueKey string `json:"issueKey"`
	} `json:"issue"`
	CreatedUser *User  `json:"createdUser"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// SharedFile represents a shared file or directory in a project.
type SharedFile struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Dir     string `json:"dir"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Updated string `json:"updated"`
}

// HumanSize formats a byte count for display.
func HumanSize(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%dB", n)
	}
}
