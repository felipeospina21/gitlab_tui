package server

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"gitlab_tui/internal/logger"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

type (
	MrCommentsQueryResponse  = []table.Row
	MrPipelinesQueryResponse = []table.Row
)

type GetMergeRequestsResponse = struct {
	ID     int    `json:"iid"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	Author struct {
		Name string `json:"name"`
	}
	MergeStatus         string `json:"merge_status"`
	DetailedMergeStatus string `json:"detailed_merge_status"`
	URL                 string `json:"web_url"`
	HasConflicts        bool   `json:"has_conflicts"`
	IsDraft             bool   `json:"draft"`
}

func GetMergeRequests(projectID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests", config.Config.BaseURL, config.Config.APIVersion, projectID)
	token := config.Config.APIToken
	mrURLParams := []string{"state=opened"}
	params := "?" + strings.Join(mrURLParams, "&")

	responseData, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r []GetMergeRequestsResponse
	if err := json.Unmarshal(responseData, &r); err != nil {
		logger.Error(err)
		return nil, err
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		n := table.Row{
			strconv.Itoa(item.ID),
			item.Title,
			item.Author.Name,
			item.MergeStatus,
			checkStatus(item.DetailedMergeStatus),
			renderIcon(item.IsDraft),
			renderIcon(item.HasConflicts),
			item.URL,
			item.Desc,
		}
		rows = append(rows, n)
	}

	return rows, nil
}

func checkStatus(status string) string {
	s := map[string]string{
		"not_approved": icon.Empty,
		"unchecked":    icon.Dash,
		"mergeable":    icon.Check,
		"checking":     icon.Clock,
	}
	return s[status]
}

func renderIcon(b bool) string {
	i := icon.Empty
	if b {
		i = icon.Check
	}

	return i
}

type GetMergeRequestCommentsResponse = struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Body   string `json:"body"`
	Author struct {
		Name string `json:"name"`
	}
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Resolved  bool   `json:"resolved"`
}

func GetMergeRequestComments(mrID string, projectID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/notes", config.Config.BaseURL, config.Config.APIVersion, projectID, mrID)
	token := config.Config.APIToken
	mrURLParams := []string{"order_by=updated_at"}
	params := "?" + strings.Join(mrURLParams, "&")

	responseData, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r []GetMergeRequestCommentsResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		logger.Error(err)
		return nil, err
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		if item.Type != "" {
			createdAt, _, _ := strings.Cut(item.CreatedAt, "T")
			UpdatedAt, _, _ := strings.Cut(item.UpdatedAt, "T")

			n := table.Row{
				strconv.Itoa(item.ID),
				item.Type,
				item.Author.Name,
				createdAt,
				UpdatedAt,
				renderIcon(item.Resolved),
				item.Body,
			}
			rows = append(rows, n)

		}
	}

	return MrCommentsQueryResponse(rows), nil
}

type GetMergeRequestPipelinesResponse = struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	Status    string `json:"status"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"web_url"`
}

func GetMergeRequestPipelines(mrID string, projectID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/pipelines", config.Config.BaseURL, config.Config.APIVersion, projectID, mrID)
	token := config.Config.APIToken

	responseData, err := fetchData(url, fetchConfig{method: "GET", params: "", token: token})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r []GetMergeRequestPipelinesResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		logger.Error(err)
		return nil, err
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		if item.Status != "success" {
			createdAt, _, _ := strings.Cut(item.CreatedAt, "T")
			UpdatedAt, _, _ := strings.Cut(item.UpdatedAt, "T")

			n := table.Row{
				strconv.Itoa(item.ID),
				strconv.Itoa(item.IID),
				item.Status,
				item.Source,
				createdAt,
				UpdatedAt,
				item.URL,
			}
			rows = append(rows, n)

		}
	}

	return MrPipelinesQueryResponse(rows), nil
}
