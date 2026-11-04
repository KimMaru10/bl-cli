package api

import (
	"fmt"
)

// GetProjects returns all projects.
func (c *Client) GetProjects() ([]Project, error) {
	data, err := c.get("/projects", nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の取得に失敗しました: %w", err)
	}
	var projects []Project
	if err := decode(data, &projects); err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の解析に失敗しました: %w", err)
	}
	return projects, nil
}

// GetProject returns a single project.
func (c *Client) GetProject(projectIDOrKey string) (*Project, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key, nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
	}
	var project Project
	if err := decode(data, &project); err != nil {
		return nil, fmt.Errorf("プロジェクトの解析に失敗しました: %w", err)
	}
	return &project, nil
}

// GetProjectUsers returns users belonging to a project.
func (c *Client) GetProjectUsers(projectIDOrKey string) ([]User, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key+"/users", nil)
	if err != nil {
		return nil, fmt.Errorf("プロジェクトユーザーの取得に失敗しました: %w", err)
	}
	var users []User
	if err := decode(data, &users); err != nil {
		return nil, fmt.Errorf("プロジェクトユーザーの解析に失敗しました: %w", err)
	}
	return users, nil
}

// GetStatuses returns statuses for a project.
func (c *Client) GetStatuses(projectIDOrKey string) ([]Status, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key+"/statuses", nil)
	if err != nil {
		return nil, fmt.Errorf("ステータス一覧の取得に失敗しました: %w", err)
	}
	var statuses []Status
	if err := decode(data, &statuses); err != nil {
		return nil, fmt.Errorf("ステータス一覧の解析に失敗しました: %w", err)
	}
	return statuses, nil
}

// GetIssueTypes returns issue types for a project.
func (c *Client) GetIssueTypes(projectIDOrKey string) ([]IssueType, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key+"/issueTypes", nil)
	if err != nil {
		return nil, fmt.Errorf("課題種別一覧の取得に失敗しました: %w", err)
	}
	var issueTypes []IssueType
	if err := decode(data, &issueTypes); err != nil {
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
	if err := decode(data, &priorities); err != nil {
		return nil, fmt.Errorf("優先度一覧の解析に失敗しました: %w", err)
	}
	return priorities, nil
}

// GetMilestones returns milestones for a project.
func (c *Client) GetMilestones(projectIDOrKey string) ([]Milestone, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key+"/versions", nil)
	if err != nil {
		return nil, fmt.Errorf("マイルストーン一覧の取得に失敗しました: %w", err)
	}
	var milestones []Milestone
	if err := decode(data, &milestones); err != nil {
		return nil, fmt.Errorf("マイルストーン一覧の解析に失敗しました: %w", err)
	}
	return milestones, nil
}

// GetCategories returns categories for a project.
func (c *Client) GetCategories(projectIDOrKey string) ([]Category, error) {
	key, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	data, err := c.get("/projects/"+key+"/categories", nil)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ一覧の取得に失敗しました: %w", err)
	}
	var categories []Category
	if err := decode(data, &categories); err != nil {
		return nil, fmt.Errorf("カテゴリ一覧の解析に失敗しました: %w", err)
	}
	return categories, nil
}
