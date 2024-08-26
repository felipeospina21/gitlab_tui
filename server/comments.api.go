package server

import (
	"encoding/json"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"gitlab_tui/tui/components/table"
	"strconv"
	"strings"
)

type Author struct {
	Name string `json:"name"`
}
type GetMergeRequestCommentsResponse struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	Author    Author
	CreatedAt string `json:"created_at"`
	Resolved  bool   `json:"resolved"`
}

func GetMergeRequestComments(url string) ([]table.Row, error) {
	token := config.GlobalConfig.APIToken
	mrURLParams := []string{"order_by=updated_at"}
	params := "?" + strings.Join(mrURLParams, "&")

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	if err != nil {
		return nil, err
	}

	var r []GetMergeRequestCommentsResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		return nil, err
	}

	// transforms response interface to match tbl Row
	var rows []table.Row
	for _, item := range r {
		if item.Type != "" {
			createdAt := table.FormatTime(item.CreatedAt)

			n := table.Row{
				createdAt,
				item.Type,
				item.Author.Name,
				renderIcon(item.Resolved, icon.Check),
				item.Body,
				strconv.Itoa(item.ID),
			}
			rows = append(rows, n)

		}
	}

	return rows, nil
}
