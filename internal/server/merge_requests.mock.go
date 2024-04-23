package server

import (
	"encoding/json"
	"gitlab_tui/internal/logger"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

func GetMergeRequestsMock() ([]table.Row, error) {
	responseData, err := os.ReadFile("planning_mr.json")
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

func GetMergeRequestCommentsMock(mrID string) ([]table.Row, error) {
	responseData, err := os.ReadFile("comments.json")
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
				strconv.FormatBool(item.Resolved),
				item.Body,
			}
			rows = append(rows, n)

		}
	}

	return MrCommentsQueryResponse(rows), nil
}

func GetMergeRequestPipelinesMock(mrID string) ([]table.Row, error) {
	responseData, err := os.ReadFile("mr_pipeline.json")
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
