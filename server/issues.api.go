package server

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/tui/components/table"
	"strconv"
	"strings"
)

type GetIssuesResponse = struct {
	CreatedAt string `json:"created_at"`
	Title     string `json:"title"`
	Desc      string `json:"description"`
	Author    struct {
		Name string `json:"name"`
	}
	Assignees []struct {
		Name string `json:"name"`
	}
	ClosedBy struct {
		Name string `json:"name"`
	}
	Labels     []string `json:"labels"`
	NotesCount int      `json:"user_notes_count"`
	URL        string   `json:"web_url"`
	ID         int      `json:"iid"`
}

type pages struct {
	Prev  string
	Next  string
	Total int
}

// TODO: add state param to fetch opened or closed issues
func GetIssues(projectID string, pageURL string) ([]table.Row, pages, error) {
	var url, params string

	if pageURL != "" {
		url = pageURL
	} else {
		url = fmt.Sprintf("%s/%s/projects/%s/issues", config.Config.BaseURL, config.Config.APIVersion, projectID)
		mrURLParams := []string{"state=opened", "per_page=30"}
		params = "?" + strings.Join(mrURLParams, "&")
	}

	token := config.Config.APIToken

	responseData, res, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	if err != nil {
		return nil, pages{}, err
	}

	var r []GetIssuesResponse
	if err := json.Unmarshal(responseData, &r); err != nil {
		return nil, pages{}, err
	}

	// Get headers
	th := res.Header.Get("x-total-pages")
	lh := res.Header.Get("link")
	prevPage, nextPage := getPagesLinks(strings.Split(lh, ","))
	totalPages, e := strconv.Atoi(th)

	if e != nil {
		return nil, pages{}, e
	}

	p := pages{Prev: prevPage, Next: nextPage, Total: totalPages}

	// transforms response interface to match tbl Row
	var rows []table.Row
	for _, item := range r {
		createdAt := table.FormatTime(item.CreatedAt)
		var labels, assignees string

		for i, assignee := range item.Assignees {
			if i == len(item.Assignees)-1 {
				assignees += assignee.Name
			} else {
				assignees += fmt.Sprintf("%s, ", assignee.Name)
			}
		}

		for i, label := range item.Labels {
			if i == len(item.Labels)-1 {
				labels += label
			} else {
				labels += fmt.Sprintf("%s, ", label)
			}
		}

		n := table.Row{
			createdAt,
			item.Title,
			assignees,
			labels,
			strconv.Itoa(item.NotesCount),
			item.Author.Name,
			item.ClosedBy.Name,
			item.URL,
			item.Desc,
			strconv.Itoa(item.ID),
		}
		rows = append(rows, n)
	}

	return rows, p, nil
}
