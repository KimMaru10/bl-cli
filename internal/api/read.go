package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// seg escapes a single path segment, rejecting "." and ".." which
// url.PathEscape leaves as-is and the server could resolve to another endpoint.
func seg(v string) (string, error) {
	if v == "" || v == "." || v == ".." {
		return "", fmt.Errorf("不正な指定です: %q", v)
	}
	return url.PathEscape(v), nil
}

// segs escapes and joins segments with "/".
func segs(vs ...string) (string, error) {
	out := make([]string, len(vs))
	for i, v := range vs {
		e, err := seg(v)
		if err != nil {
			return "", err
		}
		out[i] = e
	}
	return strings.Join(out, "/"), nil
}

// getJSON fetches path and decodes the JSON body into out.
func (c *Client) getJSON(path string, params url.Values, out any, what string) error {
	data, err := c.get(path, params)
	if err != nil {
		return fmt.Errorf("%sの取得に失敗しました: %w", what, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("%sの解析に失敗しました: %w", what, err)
	}
	return nil
}

// GetIssueAttachments returns the attachments of an issue.
func (c *Client) GetIssueAttachments(issueIDOrKey string) ([]Attachment, error) {
	key, err := seg(issueIDOrKey)
	if err != nil {
		return nil, err
	}
	var out []Attachment
	err = c.getJSON("/issues/"+key+"/attachments", nil, &out, "添付ファイル一覧")
	return out, err
}

// DownloadIssueAttachment downloads an issue attachment.
func (c *Client) DownloadIssueAttachment(issueIDOrKey string, attachmentID int) (string, []byte, error) {
	key, err := seg(issueIDOrKey)
	if err != nil {
		return "", nil, err
	}
	return c.download("/issues/" + key + "/attachments/" + strconv.Itoa(attachmentID))
}

// GetDocumentTree returns the document tree of a project.
func (c *Client) GetDocumentTree(projectIDOrKey string) (*DocumentTree, error) {
	params := url.Values{}
	params.Set("projectIdOrKey", projectIDOrKey)
	var out DocumentTree
	if err := c.getJSON("/documents/tree", params, &out, "ドキュメントツリー"); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDocuments returns documents of a project, optionally filtered by keyword.
func (c *Client) GetDocuments(projectID int, keyword string, count int) ([]Document, error) {
	params := url.Values{}
	params.Add("projectId[]", strconv.Itoa(projectID))
	// offset is required by the API.
	params.Set("offset", "0")
	params.Set("count", strconv.Itoa(clampCount(count)))
	params.Set("sort", "updated")
	if keyword != "" {
		params.Set("keyword", keyword)
	}
	var out []Document
	err := c.getJSON("/documents", params, &out, "ドキュメント一覧")
	return out, err
}

// GetDocument returns a single document including its body.
func (c *Client) GetDocument(documentID string) (*Document, error) {
	id, err := seg(documentID)
	if err != nil {
		return nil, err
	}
	var out Document
	if err := c.getJSON("/documents/"+id, nil, &out, "ドキュメント"); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadDocumentAttachment downloads a document attachment.
func (c *Client) DownloadDocumentAttachment(documentID string, attachmentID int) (string, []byte, error) {
	id, err := seg(documentID)
	if err != nil {
		return "", nil, err
	}
	return c.download("/documents/" + id + "/attachments/" + strconv.Itoa(attachmentID))
}

// GetWikis returns wiki pages of a project, optionally filtered by keyword.
func (c *Client) GetWikis(projectIDOrKey, keyword string) ([]Wiki, error) {
	params := url.Values{}
	params.Set("projectIdOrKey", projectIDOrKey)
	if keyword != "" {
		params.Set("keyword", keyword)
	}
	var out []Wiki
	err := c.getJSON("/wikis", params, &out, "Wiki 一覧")
	return out, err
}

// GetWiki returns a single wiki page including its content.
func (c *Client) GetWiki(wikiID int) (*Wiki, error) {
	var out Wiki
	if err := c.getJSON("/wikis/"+strconv.Itoa(wikiID), nil, &out, "Wiki"); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRepositories returns Git repositories of a project.
func (c *Client) GetRepositories(projectIDOrKey string) ([]Repository, error) {
	project, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	var out []Repository
	err = c.getJSON("/projects/"+project+"/git/repositories", nil, &out, "リポジトリ一覧")
	return out, err
}

// GetPullRequests returns pull requests of a repository.
func (c *Client) GetPullRequests(projectIDOrKey, repo string, count int) ([]PullRequest, error) {
	path, err := pullRequestsPath(projectIDOrKey, repo)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("count", strconv.Itoa(clampCount(count)))
	var out []PullRequest
	err = c.getJSON(path, params, &out, "プルリクエスト一覧")
	return out, err
}

// GetPullRequest returns a single pull request.
func (c *Client) GetPullRequest(projectIDOrKey, repo string, number int) (*PullRequest, error) {
	path, err := pullRequestsPath(projectIDOrKey, repo)
	if err != nil {
		return nil, err
	}
	var out PullRequest
	if err := c.getJSON(path+"/"+strconv.Itoa(number), nil, &out, "プルリクエスト"); err != nil {
		return nil, err
	}
	return &out, nil
}

func pullRequestsPath(projectIDOrKey, repo string) (string, error) {
	project, err := seg(projectIDOrKey)
	if err != nil {
		return "", err
	}
	r, err := seg(repo)
	if err != nil {
		return "", err
	}
	return "/projects/" + project + "/git/repositories/" + r + "/pullRequests", nil
}

// GetSharedFiles returns shared files under dir ("" or "/" for the root).
func (c *Client) GetSharedFiles(projectIDOrKey, dir string) ([]SharedFile, error) {
	project, err := seg(projectIDOrKey)
	if err != nil {
		return nil, err
	}
	var parts []string
	for _, p := range strings.Split(dir, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	rel, err := segs(parts...)
	if err != nil {
		return nil, err
	}
	var out []SharedFile
	err = c.getJSON("/projects/"+project+"/files/metadata/"+rel, nil, &out, "共有ファイル一覧")
	return out, err
}

// DownloadSharedFile downloads a shared file.
func (c *Client) DownloadSharedFile(projectIDOrKey string, fileID int) (string, []byte, error) {
	project, err := seg(projectIDOrKey)
	if err != nil {
		return "", nil, err
	}
	return c.download("/projects/" + project + "/files/" + strconv.Itoa(fileID))
}

func clampCount(n int) int {
	if n <= 0 {
		return 20
	}
	if n > 100 {
		return 100
	}
	return n
}
