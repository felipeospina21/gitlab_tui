package server

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"gitlab_tui/internal/logger"
	"gitlab_tui/tui/components/table"
	"strconv"
	"strings"
)

type GetMergeRequestsResponse = struct {
	CreatedAt string `json:"created_at"`
	Title     string `json:"title"`
	Desc      string `json:"description"`
	Author    struct {
		Name string `json:"name"`
	}
	DetailedMergeStatus string `json:"detailed_merge_status"`
	URL                 string `json:"web_url"`
	HasConflicts        bool   `json:"has_conflicts"`
	IsDraft             bool   `json:"draft"`
	ID                  int    `json:"iid"`
}

func GetMergeRequests(projectID string) ([]table.Row, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests", config.Config.BaseURL, config.Config.APIVersion, projectID)
	token := config.Config.APIToken
	mrURLParams := []string{"state=opened"}
	params := "?" + strings.Join(mrURLParams, "&")

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", params: params, token: token})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r []GetMergeRequestsResponse
	if err := json.Unmarshal(responseData, &r); err != nil {
		logger.Error(err)
		return nil, err
	}

	// transforms response interface to match tbl Row
	var rows []table.Row
	for _, item := range r {
		createdAt := table.FormatTime(item.CreatedAt)

		n := table.Row{
			createdAt,
			item.Title,
			item.Author.Name,
			checkMRStatus(item.DetailedMergeStatus),
			renderIcon(item.IsDraft),
			renderIcon(item.HasConflicts),
			item.URL,
			item.Desc,
			strconv.Itoa(item.ID),
		}
		rows = append(rows, n)
	}

	return rows, nil
}

func MergeMR(projectID string, mergeReqIDD string) (int, error) {
	url := fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/merge", config.Config.BaseURL, config.Config.APIVersion, projectID, mergeReqIDD)
	token := config.Config.APIToken
	mrURLParams := []string{"should_remove_source_branch=true", "squash=true"}
	params := "?" + strings.Join(mrURLParams, "&")

	_, statusCode, err := fetchData(url, fetchConfig{method: "PUT", params: params, token: token})
	if err != nil {
		logger.Error(err)
		return statusCode, err
	}

	return statusCode, nil
}

// approvals_syncing: The merge request’s approvals are syncing.
// blocked_status: Blocked by another merge request.
// checking: Git is testing if a valid merge is possible.
// ci_must_pass: A CI/CD pipeline must succeed before merge.
// ci_still_running: A CI/CD pipeline is still running.
// conflict: Conflicts exist between the source and target branches.
// discussions_not_resolved: All discussions must be resolved before merge.
// draft_status: Can’t merge because the merge request is a draft.
// external_status_checks: All status checks must pass before merge.
// jira_association_missing: The title or description must reference a Jira issue. To configure, see Require associated Jira issue for merge requests to be merged.
// mergeable: The branch can merge cleanly into the target branch.
// need_rebase: The merge request must be rebased.
// not_approved: Approval is required before merge.
// not_open: The merge request must be open before merge.
// requested_changes: The merge request has reviewers who have requested changes.
// unchecked: Git has not yet tested if a valid merge is possible.
func checkMRStatus(status string) string {
	s := map[string]string{
		"not_approved":             icon.Empty,
		"unchecked":                icon.Dash,
		"mergeable":                icon.Check,
		"checking":                 icon.Clock,
		"need_rebase":              icon.Rebase,
		"conflict":                 icon.Cross,
		"blocked_status":           icon.Cross,
		"discussions_not_resolved": icon.Discussion,
		"ci_still_running":         icon.Clock,
		"draft_status":             icon.Edit,
	}
	return s[status]
}
