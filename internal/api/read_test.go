package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "test-key")
}

func TestDownloadIssueAttachmentDecodesUTF8FileName(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/issues/PRJ-1/attachments/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''%E7%94%BB%E9%9D%A2.png")
		w.Write([]byte("PNGDATA"))
	})
	name, data, err := c.DownloadIssueAttachment("PRJ-1", 42)
	if err != nil {
		t.Fatal(err)
	}
	if name != "画面.png" || string(data) != "PNGDATA" {
		t.Fatalf("got name=%q data=%q", name, data)
	}
}

func TestGetDocumentsSendsRequiredOffset(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("offset") != "0" || q.Get("projectId[]") != "7" || q.Get("keyword") != "仕様" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Write([]byte(`[{"id":"abc","title":"仕様書"}]`))
	})
	docs, err := c.GetDocuments(7, "仕様", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 1 || docs[0].ID != "abc" {
		t.Fatalf("got %+v", docs)
	}
}

func TestGetPullRequestsEscapesRepoName(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/api/v2/projects/PRJ/git/repositories/my%20repo/pullRequests" {
			t.Errorf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Write([]byte(`[{"number":3,"summary":"fix","base":"main","branch":"feature/PRJ-1","status":{"id":1,"name":"Open"}}]`))
	})
	prs, err := c.GetPullRequests("PRJ", "my repo", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 1 || prs[0].Number != 3 || prs[0].Status.Name != "Open" {
		t.Fatalf("got %+v", prs)
	}
}

func TestFeatureRestrictedErrorIsReadable(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"msg.featureRestrictedError.title.git","code":5}]}`))
	})
	_, err := c.GetRepositories("PRJ")
	if err == nil || !strings.Contains(err.Error(), "利用できない機能") {
		t.Fatalf("got %v", err)
	}
}

func TestCustomFieldDisplayValue(t *testing.T) {
	cases := map[string]string{
		`"text"`:                      "text",
		`12.5`:                        "12.5",
		`{"id":1,"name":"高"}`:         "高",
		`[{"name":"A"},{"name":"B"}]`: "A, B",
		`null`:                        "",
	}
	for raw, want := range cases {
		f := CustomField{Value: json.RawMessage(raw)}
		if got := f.DisplayValue(); got != want {
			t.Errorf("%s: got %q want %q", raw, got, want)
		}
	}
}

func TestFilenameFromDispositionFallsBackToEmpty(t *testing.T) {
	for _, h := range []string{"", "attachment", "!!invalid;;", `attachment; filename="plain.txt"`} {
		got := filenameFromDisposition(h)
		want := ""
		if h == `attachment; filename="plain.txt"` {
			want = "plain.txt"
		}
		if got != want {
			t.Errorf("%q: got %q want %q", h, got, want)
		}
	}
}

func TestDotSegmentsAreRejected(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request should not be sent: %s", r.URL.Path)
	})
	if _, err := c.GetIssueAttachments(".."); err == nil {
		t.Error("issue key ..: want error")
	}
	if _, err := c.GetDocument("."); err == nil {
		t.Error("document id .: want error")
	}
	if _, err := c.GetSharedFiles("PRJ", "a/../b"); err == nil {
		t.Error("dir with ..: want error")
	}
}

func TestGetSharedFilesEscapesEachSegment(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/api/v2/projects/PRJ/files/metadata/a%20b/%E4%BB%95%E6%A7%98" {
			t.Errorf("unexpected path: %s", r.URL.EscapedPath())
		}
		w.Write([]byte(`[]`))
	})
	if _, err := c.GetSharedFiles("PRJ", "/a b/仕様/"); err != nil {
		t.Fatal(err)
	}
}

func TestClampCount(t *testing.T) {
	for in, want := range map[int]int{-1: 20, 0: 20, 1: 1, 100: 100, 101: 100} {
		if got := clampCount(in); got != want {
			t.Errorf("clampCount(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestTransportErrorDoesNotLeakAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close() // connection refused
	c := NewClient(srv.URL, "SECRET-KEY")
	_, err := c.GetIssueAttachments("PRJ-1")
	if err == nil || strings.Contains(err.Error(), "SECRET-KEY") {
		t.Fatalf("error leaks api key or is nil: %v", err)
	}
}
