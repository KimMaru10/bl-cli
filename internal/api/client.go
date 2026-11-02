package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

// Client is the Backlog API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Backlog API client.
func NewClient(spaceURL, apiKey string) *Client {
	base := strings.TrimRight(spaceURL, "/") + "/api/v2"
	return &Client{
		baseURL:    base,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) do(method, path string, params url.Values, body io.Reader) (*http.Response, error) {
	u := c.baseURL + path

	if params == nil {
		params = url.Values{}
	}
	params.Set("apiKey", c.apiKey)

	var req *http.Request
	var err error

	switch method {
	case http.MethodGet:
		u += "?" + params.Encode()
		req, err = http.NewRequest(method, u, nil)
	case http.MethodPost, http.MethodPatch:
		req, err = http.NewRequest(method, u+"?apiKey="+url.QueryEscape(c.apiKey), strings.NewReader(params.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	default:
		return nil, fmt.Errorf("サポートされていないHTTPメソッドです: %s", method)
	}
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// *url.Error embeds the full URL, which contains the API key.
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return nil, fmt.Errorf("リクエストの送信に失敗しました: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)

		var errResp BacklogErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && len(errResp.Errors) > 0 {
			// Feature disabled for the space (e.g. Git on some plans).
			if strings.HasPrefix(errResp.Errors[0].Message, "msg.featureRestrictedError") {
				return nil, fmt.Errorf("このスペースでは利用できない機能です（%s）", errResp.Errors[0].Message)
			}
			return nil, fmt.Errorf("Backlog API エラー: %s (code: %d)", errResp.Errors[0].Message, errResp.Errors[0].Code)
		}
		return nil, fmt.Errorf("Backlog API エラー: ステータスコード %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *Client) get(path string, params url.Values) ([]byte, error) {
	resp, err := c.do(http.MethodGet, path, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) post(path string, values url.Values) ([]byte, error) {
	resp, err := c.do(http.MethodPost, path, values, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) patch(path string, values url.Values) ([]byte, error) {
	resp, err := c.do(http.MethodPatch, path, values, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// maxDownloadBytes caps a single download held in memory.
const maxDownloadBytes int64 = 512 << 20

// download fetches a binary endpoint and returns the file name from
// Content-Disposition (empty if absent) and the body.
func (c *Client) download(path string) (string, []byte, error) {
	resp, err := c.do(http.MethodGet, path, nil, nil)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxDownloadBytes+1))
	if err != nil {
		return "", nil, fmt.Errorf("ダウンロードに失敗しました: %w", err)
	}
	if int64(len(data)) > maxDownloadBytes {
		return "", nil, fmt.Errorf("ファイルが大きすぎます（上限 %dMB）", maxDownloadBytes>>20)
	}
	return filenameFromDisposition(resp.Header.Get("Content-Disposition")), data, nil
}

// filenameFromDisposition extracts the file name, including RFC 5987 filename*.
func filenameFromDisposition(header string) string {
	if header == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return params["filename"]
}
