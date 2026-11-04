package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/browser"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/KimMaru10/bl-cli/internal/git"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type issueAttachmentDownloadArgs struct {
	IssueKey     string `json:"issue_key" jsonschema:"課題キー（例：PROJ-123）"`
	AttachmentID int    `json:"attachment_id,omitempty" jsonschema:"添付ファイル ID（省略時は全件）"`
}

type issueKeyFromBranchArgs struct{}

type docTreeArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
}

type docListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Keyword    string `json:"keyword,omitempty" jsonschema:"キーワード検索"`
	Count      int    `json:"count,omitempty" jsonschema:"取得件数（デフォルト20、最大100）"`
}

type docViewArgs struct {
	DocumentID string `json:"document_id" jsonschema:"ドキュメント ID（doc_tree または doc_list で取得）"`
}

type docAttachmentDownloadArgs struct {
	DocumentID   string `json:"document_id" jsonschema:"ドキュメント ID"`
	AttachmentID int    `json:"attachment_id,omitempty" jsonschema:"添付ファイル ID（省略時は全件）"`
}

type wikiListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Keyword    string `json:"keyword,omitempty" jsonschema:"キーワード検索"`
}

type wikiViewArgs struct {
	WikiID int `json:"wiki_id" jsonschema:"Wiki ID（wiki_list で取得）"`
}

type fileListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Dir        string `json:"dir,omitempty" jsonschema:"フォルダ（省略時はルート）"`
}

type fileDownloadArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	FileID     int    `json:"file_id" jsonschema:"共有ファイル ID（file_list で取得）"`
}

type repoListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
}

type prListArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Repo       string `json:"repo" jsonschema:"リポジトリ名（repo_list で取得）"`
	Count      int    `json:"count,omitempty" jsonschema:"取得件数（デフォルト20、最大100）"`
}

type prViewArgs struct {
	ProjectKey string `json:"project_key,omitempty" jsonschema:"プロジェクトキー（省略時はデフォルトプロジェクト）"`
	Repo       string `json:"repo" jsonschema:"リポジトリ名（repo_list で取得）"`
	Number     int    `json:"number" jsonschema:"プルリクエスト番号"`
}

func registerReadTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_attachment_download",
		Description: "Backlog の課題の添付ファイルを ~/Downloads/bl に保存し、保存先を返す（一覧は issue_view）",
	}, handleIssueAttachmentDownload)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "issue_key_from_branch",
		Description: "MCP サーバーの作業フォルダにある git ブランチ名から課題キーを推測する。書き込みの前に対象の課題が正しいか確認すること",
	}, handleIssueKeyFromBranch)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "doc_tree",
		Description: "Backlog のドキュメントツリーを ID 付きで取得する",
	}, handleDocTree)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "doc_list",
		Description: "Backlog のドキュメント一覧を更新順に取得する。キーワードで検索可能",
	}, handleDocList)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "doc_view",
		Description: "Backlog のドキュメントの本文・タグ・添付ファイル情報を取得する",
	}, handleDocView)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "doc_attachment_download",
		Description: "Backlog のドキュメントの添付ファイルを ~/Downloads/bl に保存し、保存先を返す（一覧は doc_view）",
	}, handleDocAttachmentDownload)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wiki_list",
		Description: "Backlog の Wiki ページ一覧を取得する。キーワードで検索可能",
	}, handleWikiList)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wiki_view",
		Description: "Backlog の Wiki ページの本文・タグ・添付ファイル情報を取得する",
	}, handleWikiView)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "file_list",
		Description: "Backlog の共有ファイル一覧を取得する",
	}, handleFileList)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "file_download",
		Description: "Backlog の共有ファイルを ~/Downloads/bl に保存し、保存先を返す",
	}, handleFileDownload)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "repo_list",
		Description: "Backlog Git のリポジトリ一覧を取得する",
	}, handleRepoList)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pr_list",
		Description: "Backlog Git のプルリクエスト一覧を取得する",
	}, handlePRList)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pr_view",
		Description: "Backlog Git のプルリクエスト詳細を取得する",
	}, handlePRView)
}

// --- Shared helpers ---

func clientForProject(projectKey string) (*api.Client, string, error) {
	client, defaultProject, err := newClient()
	if err != nil {
		return nil, "", err
	}
	if projectKey == "" {
		projectKey = defaultProject
	}
	if projectKey == "" {
		return nil, "", fmt.Errorf("project_key を指定してください（デフォルトプロジェクトが未設定です）")
	}
	return client, projectKey, nil
}

func milestoneID(client *api.Client, projectKey, name string) (int, error) {
	milestones, err := client.GetMilestones(projectKey)
	if err != nil {
		return 0, err
	}
	for _, m := range milestones {
		if m.Name == name {
			return m.ID, nil
		}
	}
	return 0, fmt.Errorf("マイルストーン「%s」が見つかりません", name)
}

// openInBrowser opens a page of the current space. Only the path comes from the
// tool arguments, so a request cannot open a site outside the configured space.
func openInBrowser(path string) (*mcp.CallToolResult, any, error) {
	cfg, err := config.Load()
	if err != nil {
		return errResult(err.Error())
	}
	_, space := resolveSpace(cfg)
	if space == nil {
		return errResult("認証されていません。先に bl auth login を実行してください")
	}
	pageURL := space.SpaceURL + path
	if err := browser.Open(pageURL); err != nil {
		return errResult(fmt.Sprintf("ブラウザを起動できませんでした（%v）: %s", err, pageURL))
	}
	return textResult(map[string]string{"url": pageURL, "message": "ブラウザで開きました"})
}

// warnIfStartDateIgnored reports a start date the space dropped: Backlog can answer
// success without storing it, and the model would otherwise claim it was set.
func warnIfStartDateIgnored(result map[string]string, requested, stored string) {
	if requested != "" && stored == "" {
		result["warning"] = "開始日が反映されませんでした。スペースの設定やプランで使えない可能性があります"
	}
}

type attachmentInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func attachmentInfos(attachments []api.Attachment) []attachmentInfo {
	result := make([]attachmentInfo, len(attachments))
	for i, a := range attachments {
		result[i] = attachmentInfo{ID: a.ID, Name: a.Name, Size: a.Size}
	}
	return result
}

func tagNames(tags []api.Tag) []string {
	result := make([]string, len(tags))
	for i, t := range tags {
		result[i] = t.Name
	}
	return result
}

// downloadDir is the only folder MCP downloads are written to.
// ponytail: fixed folder, because the caller is a model that reads untrusted
// Backlog text and must not choose where files land (for example ~/.ssh).
// Add an allowlist if other folders are ever needed.
func downloadDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
	}
	return filepath.Join(home, "Downloads", "bl"), nil
}

func saveAttachments(attachments []api.Attachment, id int, dl func(int) (string, []byte, error)) (*mcp.CallToolResult, any, error) {
	if len(attachments) == 0 {
		return textResult(map[string]any{"paths": []string{}, "message": "添付ファイルはありません"})
	}
	dir, err := downloadDir()
	if err != nil {
		return errResult(err.Error())
	}
	paths, err := cmdutil.SaveAttachments(attachments, id, dir, dl)
	if err != nil {
		return errResult(fmt.Sprintf("%v（保存済み: %v）", err, paths))
	}
	return textResult(map[string]any{"paths": paths})
}

// --- Handlers ---

func handleIssueAttachmentDownload(ctx context.Context, req *mcp.CallToolRequest, args issueAttachmentDownloadArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	attachments, err := client.GetIssueAttachments(args.IssueKey)
	if err != nil {
		return errResult(err.Error())
	}
	return saveAttachments(attachments, args.AttachmentID, func(id int) (string, []byte, error) {
		return client.DownloadIssueAttachment(args.IssueKey, id)
	})
}

func handleIssueKeyFromBranch(ctx context.Context, req *mcp.CallToolRequest, args issueKeyFromBranchArgs) (*mcp.CallToolResult, any, error) {
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return errResult("git ブランチを取得できませんでした。MCP サーバーの作業フォルダが git リポジトリではない可能性があります")
	}
	// Branch names come from the repository and may carry invisible characters.
	branch = cmdutil.StripControl(branch)
	key := git.ExtractIssueKey(branch)
	if key == "" {
		return errResult(fmt.Sprintf("ブランチ名「%s」に課題キーが含まれていません", branch))
	}
	return textResult(map[string]string{"issue_key": key, "branch": branch})
}

func handleDocTree(ctx context.Context, req *mcp.CallToolRequest, args docTreeArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	tree, err := client.GetDocumentTree(projectKey)
	if err != nil {
		return errResult(err.Error())
	}
	nodes := []api.DocumentNode{}
	if tree.ActiveTree != nil && tree.ActiveTree.Children != nil {
		nodes = tree.ActiveTree.Children
	}
	return textResult(nodes)
}

func handleDocList(ctx context.Context, req *mcp.CallToolRequest, args docListArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	project, err := client.GetProject(projectKey)
	if err != nil {
		return errResult(err.Error())
	}
	docs, err := client.GetDocuments(project.ID, args.Keyword, args.Count)
	if err != nil {
		return errResult(err.Error())
	}
	// ponytail: bodies are left out to keep lists small; doc_view returns them.
	type item struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Updated string `json:"updated"`
	}
	result := make([]item, len(docs))
	for i, d := range docs {
		result[i] = item{ID: d.ID, Title: d.Title, Updated: d.Updated}
	}
	return textResult(result)
}

func handleDocView(ctx context.Context, req *mcp.CallToolRequest, args docViewArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	d, err := client.GetDocument(args.DocumentID)
	if err != nil {
		return errResult(err.Error())
	}
	return textResult(struct {
		ID          string           `json:"id"`
		Title       string           `json:"title"`
		Tags        []string         `json:"tags"`
		Created     string           `json:"created"`
		Updated     string           `json:"updated"`
		Attachments []attachmentInfo `json:"attachments"`
		Body        string           `json:"body"`
	}{d.ID, d.Title, tagNames(d.Tags), d.Created, d.Updated, attachmentInfos(d.Attachments), d.Plain})
}

func handleDocAttachmentDownload(ctx context.Context, req *mcp.CallToolRequest, args docAttachmentDownloadArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	d, err := client.GetDocument(args.DocumentID)
	if err != nil {
		return errResult(err.Error())
	}
	return saveAttachments(d.Attachments, args.AttachmentID, func(id int) (string, []byte, error) {
		return client.DownloadDocumentAttachment(d.ID, id)
	})
}

func handleWikiList(ctx context.Context, req *mcp.CallToolRequest, args wikiListArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	wikis, err := client.GetWikis(projectKey, args.Keyword)
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Updated string `json:"updated"`
	}
	result := make([]item, len(wikis))
	for i, w := range wikis {
		result[i] = item{ID: w.ID, Name: w.Name, Updated: w.Updated}
	}
	return textResult(result)
}

func handleWikiView(ctx context.Context, req *mcp.CallToolRequest, args wikiViewArgs) (*mcp.CallToolResult, any, error) {
	client, _, err := newClient()
	if err != nil {
		return errResult(err.Error())
	}
	w, err := client.GetWiki(args.WikiID)
	if err != nil {
		return errResult(err.Error())
	}
	return textResult(struct {
		ID          int              `json:"id"`
		Name        string           `json:"name"`
		Tags        []string         `json:"tags"`
		Created     string           `json:"created"`
		Updated     string           `json:"updated"`
		Attachments []attachmentInfo `json:"attachments"`
		Body        string           `json:"body"`
	}{w.ID, w.Name, tagNames(w.Tags), w.Created, w.Updated, attachmentInfos(w.Attachments), w.Content})
}

func handleFileList(ctx context.Context, req *mcp.CallToolRequest, args fileListArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	files, err := client.GetSharedFiles(projectKey, args.Dir)
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		ID      int    `json:"id"`
		Type    string `json:"type"`
		Path    string `json:"path"`
		Size    int64  `json:"size"`
		Updated string `json:"updated"`
	}
	result := make([]item, len(files))
	for i, f := range files {
		result[i] = item{ID: f.ID, Type: f.Type, Path: f.Dir + f.Name, Size: f.Size, Updated: f.Updated}
	}
	return textResult(result)
}

func handleFileDownload(ctx context.Context, req *mcp.CallToolRequest, args fileDownloadArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	name, data, err := client.DownloadSharedFile(projectKey, args.FileID)
	if err != nil {
		return errResult(err.Error())
	}
	dir, err := downloadDir()
	if err != nil {
		return errResult(err.Error())
	}
	path, err := cmdutil.SaveFile(dir, name, fmt.Sprintf("file-%d", args.FileID), data)
	if err != nil {
		return errResult(err.Error())
	}
	return textResult(map[string]any{"paths": []string{path}})
}

func handleRepoList(ctx context.Context, req *mcp.CallToolRequest, args repoListArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	repos, err := client.GetRepositories(projectKey)
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		Name        string `json:"name"`
		Updated     string `json:"updated"`
		Description string `json:"description"`
	}
	result := make([]item, len(repos))
	for i, r := range repos {
		result[i] = item{Name: r.Name, Updated: r.Updated, Description: r.Description}
	}
	return textResult(result)
}

func prStatusName(p api.PullRequest) string {
	if p.Status == nil {
		return ""
	}
	return p.Status.Name
}

func handlePRList(ctx context.Context, req *mcp.CallToolRequest, args prListArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	prs, err := client.GetPullRequests(projectKey, args.Repo, args.Count)
	if err != nil {
		return errResult(err.Error())
	}
	type item struct {
		Number  int    `json:"number"`
		Status  string `json:"status"`
		Branch  string `json:"branch"`
		Base    string `json:"base"`
		Summary string `json:"summary"`
	}
	result := make([]item, len(prs))
	for i, p := range prs {
		result[i] = item{Number: p.Number, Status: prStatusName(p), Branch: p.Branch, Base: p.Base, Summary: p.Summary}
	}
	return textResult(result)
}

func handlePRView(ctx context.Context, req *mcp.CallToolRequest, args prViewArgs) (*mcp.CallToolResult, any, error) {
	client, projectKey, err := clientForProject(args.ProjectKey)
	if err != nil {
		return errResult(err.Error())
	}
	p, err := client.GetPullRequest(projectKey, args.Repo, args.Number)
	if err != nil {
		return errResult(err.Error())
	}
	result := struct {
		Number      int    `json:"number"`
		Summary     string `json:"summary"`
		Status      string `json:"status"`
		Branch      string `json:"branch"`
		Base        string `json:"base"`
		Assignee    string `json:"assignee,omitempty"`
		Issue       string `json:"issue,omitempty"`
		Created     string `json:"created"`
		Updated     string `json:"updated"`
		Description string `json:"description"`
	}{
		Number:      p.Number,
		Summary:     p.Summary,
		Status:      prStatusName(*p),
		Branch:      p.Branch,
		Base:        p.Base,
		Created:     p.Created,
		Updated:     p.Updated,
		Description: p.Description,
	}
	if p.Assignee != nil {
		result.Assignee = p.Assignee.Name
	}
	if p.Issue != nil {
		result.Issue = p.Issue.IssueKey
		if result.Issue == "" {
			result.Issue = fmt.Sprintf("ID %d", p.Issue.ID)
		}
	}
	return textResult(result)
}
