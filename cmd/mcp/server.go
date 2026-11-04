package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newClient() (*api.Client, string, error) {
	client, defaultProject, _, err := resolveClient()
	return client, defaultProject, err
}

// resolveClient also returns the alias of the space the client talks to, so a
// write can report the space it really used.
func resolveClient() (*api.Client, string, string, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, "", "", fmt.Errorf("設定の読み込みに失敗しました: %w", err)
	}
	alias, space := resolveSpace(cfg)
	if space == nil {
		if alias != "" && len(cfg.Spaces) > 0 {
			return nil, "", "", fmt.Errorf("選択中のスペース「%s」が設定にありません。auth_status で確認し、auth_switch で選び直してください", alias)
		}
		return nil, "", "", fmt.Errorf("認証されていません。先に bl auth login を実行してください")
	}
	return api.NewClient(space.SpaceURL, space.APIKey), space.DefaultProject, alias, nil
}

func textResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func errResult(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
		IsError: true,
	}, nil, nil
}

// --- Tool argument types ---

type projectListArgs struct{}

type issueListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Status     string `json:"status,omitempty" jsonschema:"ステータスで絞り込み（例：処理中、完了）"`
	AssigneeMe bool   `json:"assignee_me,omitempty" jsonschema:"自分にアサインされた課題のみ"`
	Assignee   string `json:"assignee,omitempty" jsonschema:"担当者名で絞り込み（@me で自分）"`
	Milestone  string `json:"milestone,omitempty" jsonschema:"マイルストーン名で絞り込み"`
	Parent     string `json:"parent,omitempty" jsonschema:"親課題のキー（子課題を一覧）"`
	Keyword    string `json:"keyword,omitempty" jsonschema:"キーワード検索（件名・説明・コメント）"`
	Count      int    `json:"count,omitempty" jsonschema:"取得件数（デフォルト20、最大100）"`
	Web        bool   `json:"web,omitempty" jsonschema:"一覧を返す代わりにブラウザで課題検索画面を開く"`
}

type issueViewArgs struct {
	IssueKey string `json:"issue_key" jsonschema:"課題キー（例：PROJ-123）"`
	Web      bool   `json:"web,omitempty" jsonschema:"詳細を返す代わりにブラウザで課題を開く"`
}

type issueCreateArgs struct {
	ProjectKey   string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Summary      string `json:"summary" jsonschema:"課題のタイトル"`
	IssueType    string `json:"issue_type" jsonschema:"課題種別名（例：タスク、バグ）"`
	Priority     string `json:"priority" jsonschema:"優先度名（例：高、中、低）"`
	Description  string `json:"description,omitempty" jsonschema:"課題の詳細"`
	AssigneeName string `json:"assignee_name,omitempty" jsonschema:"担当者名"`
	DueDate      string `json:"due_date,omitempty" jsonschema:"期日（yyyy-MM-dd）"`
	StartDate    string `json:"start_date,omitempty" jsonschema:"開始日（yyyy-MM-dd）"`
	Milestone    string `json:"milestone,omitempty" jsonschema:"マイルストーン名"`
}

type issueEditArgs struct {
	IssueKey     string `json:"issue_key" jsonschema:"課題キー（例：PROJ-123）"`
	Status       string `json:"status,omitempty" jsonschema:"変更先のステータス名"`
	AssigneeName string `json:"assignee_name,omitempty" jsonschema:"担当者名"`
	Priority     string `json:"priority,omitempty" jsonschema:"優先度名"`
	DueDate      string `json:"due_date,omitempty" jsonschema:"期日（yyyy-MM-dd）"`
	StartDate    string `json:"start_date,omitempty" jsonschema:"開始日（yyyy-MM-dd）"`
	Milestone    string `json:"milestone,omitempty" jsonschema:"マイルストーン名"`
	Unassign     bool   `json:"unassign,omitempty" jsonschema:"担当者を未設定に戻す"`
	ClearDueDate bool   `json:"clear_due_date,omitempty" jsonschema:"期日を消す"`
	Comment      string `json:"comment,omitempty" jsonschema:"更新時のコメント"`
}

type commentAddArgs struct {
	IssueKey string `json:"issue_key" jsonschema:"課題キー（例：PROJ-123）"`
	Body     string `json:"body" jsonschema:"コメント本文"`
}

type commentListArgs struct {
	IssueKey string `json:"issue_key" jsonschema:"課題キー（例：PROJ-123）"`
	Count    int    `json:"count,omitempty" jsonschema:"取得件数（デフォルト20）"`
}

// Version is reported to MCP clients. It is set from the CLI version at startup.
var Version = "dev"

// Run starts the MCP server over stdio.
func Run() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "bl-backlog",
		Version: Version,
	}, nil)

	registerTools(server)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

func registerTools(server *mcp.Server) {
	// project_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "project_list",
		Description: "Backlog のプロジェクト一覧を取得する",
	}, handleProjectList)

	// issue_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_list",
		Description: "Backlog の課題一覧を更新順に取得する。ステータス・担当者・マイルストーン・親課題・キーワードで絞り込み可能",
	}, handleIssueList)

	// issue_view
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_view",
		Description: "Backlog の課題詳細を取得する",
	}, handleIssueView)

	// issue_create
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_create",
		Description: "Backlog に新しい課題を作成する",
	}, handleIssueCreate)

	// issue_edit
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_edit",
		Description: "Backlog の課題を更新する（ステータス変更、担当者変更など）",
	}, handleIssueEdit)

	// comment_add
	mcp.AddTool(server, &mcp.Tool{
		Name:        "comment_add",
		Description: "Backlog の課題にコメントを追加する",
	}, handleCommentAdd)

	// comment_list
	mcp.AddTool(server, &mcp.Tool{
		Name:        "comment_list",
		Description: "Backlog の課題のコメント一覧を取得する",
	}, handleCommentList)

	registerReadTools(server)
	registerSettingsTools(server)
}

// --- Handlers ---

func handleProjectList(ctx context.Context, req *mcp.CallToolRequest, args projectListArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	projects, err := client.GetProjects()
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}
	result := make([]item, len(projects))
	for i, p := range projects {
		result[i] = item{Key: p.ProjectKey, Name: p.Name}
	}
	return textResult(result)
}

func handleIssueList(ctx context.Context, req *mcp.CallToolRequest, args issueListArgs) (*mcp.CallToolResult, any, error) {
	client, defaultProject, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}

	projectKey := args.ProjectKey
	if projectKey == "" {
		projectKey = defaultProject
	}
	if projectKey == "" {
		return errResult("project_key を指定してください（デフォルトプロジェクトが未設定です）")
	}

	if args.Web {
		return openInBrowser("/find/" + url.PathEscape(projectKey))
	}

	project, err := client.GetProject(projectKey)
	if err != nil {
		return errResult(err.Error())
	}

	opts := &api.GetIssuesOptions{
		ProjectIDs: []int{project.ID},
		Count:      args.Count,
		Sort:       "updated",
		Order:      "desc",
	}

	if args.AssigneeMe || args.Assignee == "@me" {
		me, err := client.GetMyself()
		if err != nil {
			return errResult(err.Error())
		}
		opts.AssigneeIDs = []int{me.ID}
	} else if args.Assignee != "" {
		users, err := client.GetProjectUsers(projectKey)
		if err != nil {
			return errResult(err.Error())
		}
		for _, u := range users {
			if u.Name == args.Assignee {
				opts.AssigneeIDs = []int{u.ID}
				break
			}
		}
		if opts.AssigneeIDs == nil {
			return errResult(fmt.Sprintf("担当者「%s」が見つかりません", args.Assignee))
		}
	}

	if args.Milestone != "" {
		id, err := milestoneID(client, projectKey, args.Milestone)
		if err != nil {
			return errResult(err.Error())
		}
		opts.MilestoneIDs = []int{id}
	}

	if args.Parent != "" {
		parent, err := client.GetIssue(args.Parent)
		if err != nil {
			return errResult(err.Error())
		}
		opts.ParentIDs = []int{parent.ID}
	}

	if args.Status != "" {
		statuses, err := client.GetStatuses(projectKey)
		if err != nil {
			return errResult(err.Error())
		}
		for _, s := range statuses {
			if s.Name == args.Status {
				opts.StatusIDs = []int{s.ID}
				break
			}
		}
		// An unknown name must not silently widen the result to every status.
		if opts.StatusIDs == nil {
			return errResult(fmt.Sprintf("ステータス「%s」が見つかりません", args.Status))
		}
	}

	if args.Keyword != "" {
		opts.Keyword = args.Keyword
	}

	issues, err := client.GetIssues(opts)
	if err != nil {
		return errResult(err.Error())
	}

	type item struct {
		Key      string `json:"key"`
		Summary  string `json:"summary"`
		Status   string `json:"status"`
		Assignee string `json:"assignee"`
		Priority string `json:"priority"`
		DueDate  string `json:"due_date,omitempty"`
	}
	result := make([]item, len(issues))
	for i, iss := range issues {
		it := item{
			Key:     iss.IssueKey,
			Summary: iss.Summary,
		}
		if iss.Status != nil {
			it.Status = iss.Status.Name
		}
		if iss.Assignee != nil {
			it.Assignee = iss.Assignee.Name
		}
		if iss.Priority != nil {
			it.Priority = iss.Priority.Name
		}
		it.DueDate = iss.DueDate
		result[i] = it
	}
	return textResult(result)
}

func handleIssueView(ctx context.Context, req *mcp.CallToolRequest, args issueViewArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	if args.Web {
		return openInBrowser("/view/" + url.PathEscape(args.IssueKey))
	}
	issue, err := client.GetIssue(args.IssueKey)
	if err != nil {
		return errResult(err.Error())
	}

	type customField struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	type result struct {
		Key         string   `json:"key"`
		Summary     string   `json:"summary"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		Assignee    string   `json:"assignee"`
		Priority    string   `json:"priority"`
		IssueType   string   `json:"issue_type"`
		DueDate     string   `json:"due_date,omitempty"`
		StartDate   string   `json:"start_date,omitempty"`
		Created     string   `json:"created"`
		Updated     string   `json:"updated"`
		Milestones  []string `json:"milestones,omitempty"`
		Categories  []string `json:"categories,omitempty"`
		CreatedUser string   `json:"created_user,omitempty"`
		Parent      string   `json:"parent,omitempty"`

		CustomFields []customField    `json:"custom_fields,omitempty"`
		Attachments  []attachmentInfo `json:"attachments,omitempty"`
	}

	r := result{
		Key:         issue.IssueKey,
		Summary:     issue.Summary,
		Description: issue.Description,
		DueDate:     issue.DueDate,
		StartDate:   issue.StartDate,
		Created:     issue.Created,
		Updated:     issue.Updated,
	}
	if issue.Status != nil {
		r.Status = issue.Status.Name
	}
	if issue.Assignee != nil {
		r.Assignee = issue.Assignee.Name
	}
	if issue.Priority != nil {
		r.Priority = issue.Priority.Name
	}
	if issue.IssueType != nil {
		r.IssueType = issue.IssueType.Name
	}
	for _, m := range issue.Milestone {
		r.Milestones = append(r.Milestones, m.Name)
	}
	for _, c := range issue.Category {
		r.Categories = append(r.Categories, c.Name)
	}
	if issue.CreatedUser != nil {
		r.CreatedUser = issue.CreatedUser.Name
	}
	if issue.ParentIssueID != nil {
		r.Parent = fmt.Sprintf("ID %d（取得できませんでした）", *issue.ParentIssueID)
		if p, err := client.GetIssue(fmt.Sprintf("%d", *issue.ParentIssueID)); err == nil {
			r.Parent = p.IssueKey + " " + p.Summary
		}
	}
	for _, f := range issue.CustomFields {
		if v := f.DisplayValue(); v != "" {
			r.CustomFields = append(r.CustomFields, customField{Name: f.Name, Value: v})
		}
	}
	r.Attachments = attachmentInfos(issue.Attachments)

	return textResult(r)
}

func handleIssueCreate(ctx context.Context, req *mcp.CallToolRequest, args issueCreateArgs) (*mcp.CallToolResult, any, error) {
	client, defaultProject, alias, err := resolveClient()
	if err != nil {
		return errResult(err.Error())
	}

	projectKey := args.ProjectKey
	if projectKey == "" {
		projectKey = defaultProject
	}
	if projectKey == "" {
		return errResult("project_key を指定してください")
	}

	project, err := client.GetProject(projectKey)
	if err != nil {
		return errResult(err.Error())
	}

	// Resolve issue type
	issueTypes, err := client.GetIssueTypes(projectKey)
	if err != nil {
		return errResult(err.Error())
	}
	var issueTypeID int
	for _, t := range issueTypes {
		if t.Name == args.IssueType {
			issueTypeID = t.ID
			break
		}
	}
	if issueTypeID == 0 {
		names := make([]string, len(issueTypes))
		for i, t := range issueTypes {
			names[i] = t.Name
		}
		return errResult(fmt.Sprintf("課題種別「%s」が見つかりません。選択肢: %v", args.IssueType, names))
	}

	// Resolve priority
	priorities, err := client.GetPriorities()
	if err != nil {
		return errResult(err.Error())
	}
	var priorityID int
	for _, p := range priorities {
		if p.Name == args.Priority {
			priorityID = p.ID
			break
		}
	}
	if priorityID == 0 {
		names := make([]string, len(priorities))
		for i, p := range priorities {
			names[i] = p.Name
		}
		return errResult(fmt.Sprintf("優先度「%s」が見つかりません。選択肢: %v", args.Priority, names))
	}

	opts := &api.CreateIssueOptions{
		ProjectID:   project.ID,
		Summary:     args.Summary,
		IssueTypeID: issueTypeID,
		PriorityID:  priorityID,
		Description: args.Description,
		DueDate:     args.DueDate,
		StartDate:   args.StartDate,
	}

	if args.Milestone != "" {
		id, err := milestoneID(client, projectKey, args.Milestone)
		if err != nil {
			return errResult(err.Error())
		}
		opts.MilestoneIDs = []int{id}
	}

	// Resolve assignee
	if args.AssigneeName != "" {
		users, err := client.GetProjectUsers(projectKey)
		if err != nil {
			return errResult(err.Error())
		}
		for _, u := range users {
			if u.Name == args.AssigneeName {
				opts.AssigneeID = u.ID
				break
			}
		}
		if opts.AssigneeID == 0 {
			return errResult(fmt.Sprintf("担当者「%s」が見つかりません", args.AssigneeName))
		}
	}

	issue, err := client.CreateIssue(opts)
	if err != nil {
		return errResult(err.Error())
	}

	r := map[string]string{
		"issue_key": issue.IssueKey,
		"summary":   issue.Summary,
		"space":     alias,
		"message":   fmt.Sprintf("%s を作成しました", issue.IssueKey),
	}
	warnIfStartDateIgnored(r, args.StartDate, issue.StartDate)
	return textResult(r)
}

func handleIssueEdit(ctx context.Context, req *mcp.CallToolRequest, args issueEditArgs) (*mcp.CallToolResult, any, error) {
	if args.Unassign && args.AssigneeName != "" {
		return errResult("assignee_name と unassign は同時に指定できません")
	}
	if args.ClearDueDate && args.DueDate != "" {
		return errResult("due_date と clear_due_date は同時に指定できません")
	}
	client, _, alias, err := resolveClient()
	if err != nil {
		return errResult(err.Error())
	}

	// Get current issue to determine project
	issue, err := client.GetIssue(args.IssueKey)
	if err != nil {
		return errResult(err.Error())
	}

	opts := &api.UpdateIssueOptions{}

	if args.Status != "" {
		statuses, err := client.GetStatuses(fmt.Sprintf("%d", issue.ProjectID))
		if err != nil {
			return errResult(err.Error())
		}
		found := false
		for _, s := range statuses {
			if s.Name == args.Status {
				opts.StatusID = &s.ID
				found = true
				break
			}
		}
		if !found {
			names := make([]string, len(statuses))
			for i, s := range statuses {
				names[i] = s.Name
			}
			return errResult(fmt.Sprintf("ステータス「%s」が見つかりません。選択肢: %v", args.Status, names))
		}
	}

	if args.Priority != "" {
		priorities, err := client.GetPriorities()
		if err != nil {
			return errResult(err.Error())
		}
		found := false
		for _, p := range priorities {
			if p.Name == args.Priority {
				opts.PriorityID = &p.ID
				found = true
				break
			}
		}
		if !found {
			return errResult(fmt.Sprintf("優先度「%s」が見つかりません", args.Priority))
		}
	}

	if args.AssigneeName != "" {
		users, err := client.GetProjectUsers(fmt.Sprintf("%d", issue.ProjectID))
		if err != nil {
			return errResult(err.Error())
		}
		found := false
		for _, u := range users {
			if u.Name == args.AssigneeName {
				opts.AssigneeID = &u.ID
				found = true
				break
			}
		}
		if !found {
			return errResult(fmt.Sprintf("担当者「%s」が見つかりません", args.AssigneeName))
		}
	}

	if args.Unassign {
		unassigned := 0
		opts.AssigneeID = &unassigned
	}

	if args.DueDate != "" {
		opts.DueDate = &args.DueDate
	}

	if args.ClearDueDate {
		noDueDate := ""
		opts.DueDate = &noDueDate
	}

	if args.StartDate != "" {
		opts.StartDate = &args.StartDate
	}

	if args.Milestone != "" {
		id, err := milestoneID(client, fmt.Sprintf("%d", issue.ProjectID), args.Milestone)
		if err != nil {
			return errResult(err.Error())
		}
		opts.MilestoneIDs = []int{id}
	}

	if args.Comment != "" {
		opts.Comment = &args.Comment
	}

	updated, err := client.UpdateIssue(args.IssueKey, opts)
	if err != nil {
		return errResult(err.Error())
	}

	r := map[string]string{
		"issue_key": updated.IssueKey,
		"space":     alias,
		"message":   fmt.Sprintf("%s を更新しました", updated.IssueKey),
	}
	if updated.Status != nil {
		r["status"] = updated.Status.Name
	}
	warnIfStartDateIgnored(r, args.StartDate, updated.StartDate)
	return textResult(r)
}

func handleCommentAdd(ctx context.Context, req *mcp.CallToolRequest, args commentAddArgs) (*mcp.CallToolResult, any, error) {
	client, _, alias, err := resolveClient()
	if err != nil {
		return errResult(err.Error())
	}
	comment, err := client.AddComment(args.IssueKey, args.Body)
	if err != nil {
		return errResult(err.Error())
	}
	return textResult(map[string]any{
		"comment_id": comment.ID,
		"space":      alias,
		"message":    fmt.Sprintf("%s にコメントを追加しました", args.IssueKey),
	})
}

func handleCommentList(ctx context.Context, req *mcp.CallToolRequest, args commentListArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	count := args.Count
	if count <= 0 {
		count = 20
	}
	comments, err := client.GetComments(args.IssueKey, count, "desc")
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
		Author  string `json:"author"`
		Created string `json:"created"`
	}
	result := make([]item, len(comments))
	for i, c := range comments {
		it := item{
			ID:      c.ID,
			Content: c.Content,
			Created: c.Created,
		}
		if c.CreatedUser != nil {
			it.Author = c.CreatedUser.Name
		}
		result[i] = it
	}
	return textResult(result)
}
