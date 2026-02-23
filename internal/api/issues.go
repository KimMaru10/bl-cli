package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetIssuesOptions holds parameters for GetIssues.
type GetIssuesOptions struct {
	ProjectIDs   []int
	AssigneeIDs  []int
	StatusIDs    []int
	MilestoneIDs []int
	Keyword      string
	Count        int
	Offset       int
	Sort         string
	Order        string
}

// GetIssues returns issues matching the given options.
func (c *Client) GetIssues(opts *GetIssuesOptions) ([]Issue, error) {
	params := url.Values{}

	for _, id := range opts.ProjectIDs {
		params.Add("projectId[]", strconv.Itoa(id))
	}
	for _, id := range opts.AssigneeIDs {
		params.Add("assigneeId[]", strconv.Itoa(id))
	}
	for _, id := range opts.StatusIDs {
		params.Add("statusId[]", strconv.Itoa(id))
	}
	for _, id := range opts.MilestoneIDs {
		params.Add("milestoneId[]", strconv.Itoa(id))
	}
	if opts.Keyword != "" {
		params.Set("keyword", opts.Keyword)
	}

	count := opts.Count
	if count <= 0 {
		count = 20
	}
	if count > 100 {
		count = 100
	}
	params.Set("count", strconv.Itoa(count))

	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	if opts.Order != "" {
		params.Set("order", opts.Order)
	}

	data, err := c.get("/issues", params)
	if err != nil {
		return nil, fmt.Errorf("課題一覧の取得に失敗しました: %w", err)
	}

	var issues []Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, fmt.Errorf("課題一覧の解析に失敗しました: %w", err)
	}
	return issues, nil
}

// GetIssue returns a single issue by key or ID.
func (c *Client) GetIssue(issueIDOrKey string) (*Issue, error) {
	data, err := c.get("/issues/"+issueIDOrKey, nil)
	if err != nil {
		return nil, fmt.Errorf("課題の取得に失敗しました: %w", err)
	}
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("課題の解析に失敗しました: %w", err)
	}
	return &issue, nil
}
