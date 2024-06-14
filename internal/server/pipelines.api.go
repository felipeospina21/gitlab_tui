package server

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/logger"
	"gitlab_tui/tui/components/table"
	"strconv"
)

type GetMergeRequestPipelinesResponse = struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	Status    string `json:"status"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"web_url"`
}

func GetMergeRequestPipelines(mrID string, projectID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/pipelines", config.Config.BaseURL, config.Config.APIVersion, projectID, mrID)
	token := config.Config.APIToken

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", params: "", token: token})
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
		createdAt := table.FormatTime(item.CreatedAt)

		n := table.Row{
			strconv.Itoa(item.ID),
			strconv.Itoa(item.IID),
			createdAt,
			item.Status,
			item.Source,
			item.URL,
		}
		rows = append(rows, n)

	}

	return rows, nil
}

type PipelineJobsResponse = struct {
	ID        int     `json:"id"`
	Status    string  `json:"status"`
	Stage     string  `json:"stage"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	Duration  float32 `json:"duration"`
	Coverage  float32 `json:"coverage"`
	URL       string  `json:"web_url"`
}

func GetPipelineJobs(projectID string, pipelineID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/pipelines/%s/jobs", config.Config.BaseURL, config.Config.APIVersion, projectID, pipelineID)
	token := config.Config.APIToken

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", params: "", token: token})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r []PipelineJobsResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		logger.Error(err)
		return nil, err
	}

	// transforms response interface to match table Row
	var rows []table.Row
	for _, item := range r {
		createdAt := table.FormatTime(item.CreatedAt)

		n := table.Row{
			strconv.Itoa(item.ID),
			createdAt,
			item.Status,
			item.Name,
			item.Stage,
			fmt.Sprintf("%f", item.Duration),
			fmt.Sprintf("%f", item.Coverage),
			item.URL,
		}
		rows = append(rows, n)

	}

	return rows, nil
}