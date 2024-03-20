package api

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
)

type Res = struct {
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
	// url := fmt.Sprintf("%s/%s/projects/%s/merge_requests", Config.BaseUrl, Config.ApiVersion, Config.ProjectsId.PlanningTool)
	// token:= config.Config.ApiToken
	// mrUrlParams := []string{"state=opened", "something=value"}
	// params := "?" + strings.Join(mrUrlParams, "&")
	//
	// responseData, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	responseData, err := os.ReadFile("planning_mr.json")
	if err != nil {
		log.Fatal(err)
	}

	var r []Res
	if err := json.Unmarshal(responseData, &r); err != nil {
		log.Fatal(err)
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		n := table.Row{
			strconv.Itoa(item.Id),
			item.Title,
			item.Desc,
			item.Author.Name,
			item.MergeStatus,
			strconv.FormatBool(item.IsDraft),
			strconv.FormatBool(item.HasConflicts),
			item.Url,
		}
		rows = append(rows, n)
	}

	return rows
}
