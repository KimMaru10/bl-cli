package api

import (
	"encoding/json"
	"fmt"
)

// GetMyself returns the authenticated user.
func (c *Client) GetMyself() (*User, error) {
	data, err := c.get("/users/myself", nil)
	if err != nil {
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %w", err)
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("ユーザー情報の解析に失敗しました: %w", err)
	}
	return &user, nil
}
