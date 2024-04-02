package api

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

type GetMergeRequestsResponse = struct {
	Id     int    `json:"iid"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	Author struct {
		Name string `json:"name"`
	}
	MergeStatus  string `json:"merge_status"`
	Url          string `json:"web_url"`
	HasConflicts bool   `json:"has_conflicts"`
	IsDraft      bool   `json:"draft"`
}

func GetMergeRequests() []table.Row {
	// url := fmt.Sprintf("%s/%s/projects/%s/merge_requests", config.Config.BaseUrl, config.Config.ApiVersion, config.Config.ProjectsId.PlanningTool)
	// token := config.Config.ApiToken
	// mrUrlParams := []string{"state=opened"}
	// params := "?" + strings.Join(mrUrlParams, "&")
	//
	// responseData, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	responseData, err := os.ReadFile("planning_mr.json")
	if err != nil {
		log.Fatal(err)
	}

	var r []GetMergeRequestsResponse
	if err := json.Unmarshal(responseData, &r); err != nil {
		log.Fatal(err)
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		n := table.Row{
			strconv.Itoa(item.Id),
			item.Title,
			item.Author.Name,
			item.MergeStatus,
			strconv.FormatBool(item.IsDraft),
			strconv.FormatBool(item.HasConflicts),
			item.Url,
			item.Desc,
		}
		rows = append(rows, n)
	}

	return rows
}

type GetMergeRequestsCommentsResponse = struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Body   string `json:"body"`
	Author struct {
		Name string `json:"name"`
	}
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Resolved  bool   `json:"resolved"`
}

func GetMRComments(mrId string) ([]table.Row, error) {
	// url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/notes", config.Config.BaseUrl, config.Config.ApiVersion, config.Config.ProjectsId.PlanningTool, mrId)
	// token := config.Config.ApiToken
	// mrUrlParams := []string{"order_by=updated_at"}
	// params := "?" + strings.Join(mrUrlParams, "&")

	// responseData, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	responseData, err := os.ReadFile("comments.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var r []GetMergeRequestsCommentsResponse
	// if err := json.Unmarshal(responseData, &r); err != nil {
	// 	log.Fatal(err)
	// }

	err = json.Unmarshal(responseData, &r)

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		if item.Type != "" {
			createdAt, _, _ := strings.Cut(item.CreatedAt, "T")
			UpdatedAt, _, _ := strings.Cut(item.UpdatedAt, "T")

			n := table.Row{
				strconv.Itoa(item.Id),
				item.Type,
				item.Author.Name,
				createdAt,
				UpdatedAt,
				strconv.FormatBool(item.Resolved),
				item.Body,
			}
			rows = append(rows, n)

		}
	}

	return rows, err
}
