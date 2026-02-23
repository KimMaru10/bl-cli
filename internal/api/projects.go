package api

import (
	"encoding/json"
	"fmt"
)

// GetProjects returns all projects.
func (c *Client) GetProjects() ([]Project, error) {
	data, err := c.get("/projects", nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の取得に失敗しました: %w", err)
	}
	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の解析に失敗しました: %w", err)
	}
	return projects, nil
}

// GetProject returns a single project.
func (c *Client) GetProject(projectIDOrKey string) (*Project, error) {
	data, err := c.get("/projects/"+projectIDOrKey, nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("プロジェクトの解析に失敗しました: %w", err)
	}
	return &project, nil
}

// GetProjectUsers returns users belonging to a project.
func (c *Client) GetProjectUsers(projectIDOrKey string) ([]User, error) {
	data, err := c.get("/projects/"+projectIDOrKey+"/users", nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクトユーザーの取得に失敗しました: %w", err)
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("プロジェクトユーザーの解析に失敗しました: %w", err)
	}
	return users, nil
}

// GetStatuses returns statuses for a project.
func (c *Client) GetStatuses(projectIDOrKey string) ([]Status, error) {
	data, err := c.get("/projects/"+projectIDOrKey+"/statuses", nil)
	if err != nil {
		return nil, fmt.Errorf("ステータス一覧の取得に失敗しました: %w", err)
	}
	var statuses []Status
	if err := json.Unmarshal(data, &statuses); err != nil {
		return nil, fmt.Errorf("ステータス一覧の解析に失敗しました: %w", err)
	}
	return statuses, nil
}

// GetIssueTypes returns issue types for a project.
func (c *Client) GetIssueTypes(projectIDOrKey string) ([]IssueType, error) {
	data, err := c.get("/projects/"+projectIDOrKey+"/issueTypes", nil)
	if err != nil {
		return nil, fmt.Errorf("課題種別一覧の取得に失敗しました: %w", err)
	}
	var issueTypes []IssueType
	if err := json.Unmarshal(data, &issueTypes); err != nil {
		return nil, fmt.Errorf("課題種別一覧の解析に失敗しました: %w", err)
	}
	return issueTypes, nil
}

// GetPriorities returns all priorities.
func (c *Client) GetPriorities() ([]Priority, error) {
	data, err := c.get("/priorities", nil)
	if err != nil {
		return nil, fmt.Errorf("優先度一覧の取得に失敗しました: %w", err)
	}
	var priorities []Priority
	if err := json.Unmarshal(data, &priorities); err != nil {
		return nil, fmt.Errorf("優先度一覧の解析に失敗しました: %w", err)
	}
	return priorities, nil
}

// GetMilestones returns milestones for a project.
func (c *Client) GetMilestones(projectIDOrKey string) ([]Milestone, error) {
	data, err := c.get("/projects/"+projectIDOrKey+"/versions", nil)
	if err != nil {
		return nil, fmt.Errorf("マイルストーン一覧の取得に失敗しました: %w", err)
	}
	var milestones []Milestone
	if err := json.Unmarshal(data, &milestones); err != nil {
		return nil, fmt.Errorf("マイルストーン一覧の解析に失敗しました: %w", err)
	}
	return milestones, nil
}

// GetCategories returns categories for a project.
func (c *Client) GetCategories(projectIDOrKey string) ([]Category, error) {
	data, err := c.get("/projects/"+projectIDOrKey+"/categories", nil)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ一覧の取得に失敗しました: %w", err)
	}
	var categories []Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, fmt.Errorf("カテゴリ一覧の解析に失敗しました: %w", err)
	}
	return categories, nil
}
