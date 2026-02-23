package api

import (
	"encoding/json"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("リクエストの送信に失敗しました: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)

		var errResp BacklogErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && len(errResp.Errors) > 0 {
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
