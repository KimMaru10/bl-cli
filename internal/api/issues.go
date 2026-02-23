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

// CreateIssueOptions holds parameters for CreateIssue.
type CreateIssueOptions struct {
	ProjectID   int
	Summary     string
	IssueTypeID int
	PriorityID  int
	Description string
	AssigneeID  int
	DueDate     string
	StartDate   string
	MilestoneIDs []int
	CategoryIDs  []int
}

// CreateIssue creates a new issue.
func (c *Client) CreateIssue(opts *CreateIssueOptions) (*Issue, error) {
	params := url.Values{}
	params.Set("projectId", strconv.Itoa(opts.ProjectID))
	params.Set("summary", opts.Summary)
	params.Set("issueTypeId", strconv.Itoa(opts.IssueTypeID))
	params.Set("priorityId", strconv.Itoa(opts.PriorityID))

	if opts.Description != "" {
		params.Set("description", opts.Description)
	}
	if opts.AssigneeID > 0 {
		params.Set("assigneeId", strconv.Itoa(opts.AssigneeID))
	}
	if opts.DueDate != "" {
		params.Set("dueDate", opts.DueDate)
	}
	if opts.StartDate != "" {
		params.Set("startDate", opts.StartDate)
	}
	for _, id := range opts.MilestoneIDs {
		params.Add("milestoneId[]", strconv.Itoa(id))
	}
	for _, id := range opts.CategoryIDs {
		params.Add("categoryId[]", strconv.Itoa(id))
	}

	data, err := c.post("/issues", params)
	if err != nil {
		return nil, fmt.Errorf("課題の作成に失敗しました: %w", err)
	}
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("課題の解析に失敗しました: %w", err)
	}
	return &issue, nil
}

// UpdateIssueOptions holds parameters for UpdateIssue.
// Pointer types are used to distinguish between unset and zero values.
type UpdateIssueOptions struct {
	Summary      *string
	Description  *string
	StatusID     *int
	AssigneeID   *int
	PriorityID   *int
	DueDate      *string
	StartDate    *string
	MilestoneIDs []int
	CategoryIDs  []int
	Comment      *string
}

// UpdateIssue updates an existing issue.
func (c *Client) UpdateIssue(issueIDOrKey string, opts *UpdateIssueOptions) (*Issue, error) {
	params := url.Values{}

	if opts.Summary != nil {
		params.Set("summary", *opts.Summary)
	}
	if opts.Description != nil {
		params.Set("description", *opts.Description)
	}
	if opts.StatusID != nil {
		params.Set("statusId", strconv.Itoa(*opts.StatusID))
	}
	if opts.AssigneeID != nil {
		params.Set("assigneeId", strconv.Itoa(*opts.AssigneeID))
	}
	if opts.PriorityID != nil {
		params.Set("priorityId", strconv.Itoa(*opts.PriorityID))
	}
	if opts.DueDate != nil {
		params.Set("dueDate", *opts.DueDate)
	}
	if opts.StartDate != nil {
		params.Set("startDate", *opts.StartDate)
	}
	for _, id := range opts.MilestoneIDs {
		params.Add("milestoneId[]", strconv.Itoa(id))
	}
	for _, id := range opts.CategoryIDs {
		params.Add("categoryId[]", strconv.Itoa(id))
	}
	if opts.Comment != nil {
		params.Set("comment", *opts.Comment)
	}

	data, err := c.patch("/issues/"+issueIDOrKey, params)
	if err != nil {
		return nil, fmt.Errorf("課題の更新に失敗しました: %w", err)
	}
	var updated Issue
	if err := json.Unmarshal(data, &updated); err != nil {
		return nil, fmt.Errorf("課題の解析に失敗しました: %w", err)
	}
	return &updated, nil
}

// AddComment adds a comment to an issue.
func (c *Client) AddComment(issueIDOrKey string, content string) (*Comment, error) {
	params := url.Values{}
	params.Set("content", content)

	data, err := c.post("/issues/"+issueIDOrKey+"/comments", params)
	if err != nil {
		return nil, fmt.Errorf("コメントの追加に失敗しました: %w", err)
	}
	var comment Comment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, fmt.Errorf("コメントの解析に失敗しました: %w", err)
	}
	return &comment, nil
}

// GetComments returns comments for an issue.
func (c *Client) GetComments(issueIDOrKey string, count int, order string) ([]Comment, error) {
	params := url.Values{}
	if count > 0 {
		params.Set("count", strconv.Itoa(count))
	}
	if order != "" {
		params.Set("order", order)
	}

	data, err := c.get("/issues/"+issueIDOrKey+"/comments", params)
	if err != nil {
		return nil, fmt.Errorf("コメント一覧の取得に失敗しました: %w", err)
	}
	var comments []Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, fmt.Errorf("コメント一覧の解析に失敗しました: %w", err)
	}
	return comments, nil
}
