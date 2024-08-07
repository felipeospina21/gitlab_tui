package server

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"gitlab_tui/tui/components/table"
	"strconv"
	"strings"
)

type GetMergeRequestCommentsResponse = struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Body   string `json:"body"`
	Author struct {
		Name string `json:"name"`
	}
	CreatedAt string `json:"created_at"`
	Resolved  bool   `json:"resolved"`
}

func GetMergeRequestComments(projectID string, mrID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/notes", config.Config.BaseURL, config.Config.APIVersion, projectID, mrID)
	token := config.Config.APIToken
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
