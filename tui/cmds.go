package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components/table"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) openInBrowser(tableColIdx table.TableColIndex, view views) {
	selectedURL := m.getSelectedMrRow(tableColIdx, view)
	exec.Openbrowser(selectedURL)
}

// Merge Requests Table
func (m *Model) refetchMrList() {
	r, err := server.GetMergeRequests(m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}

	m.MergeRequests.List.SetRows(r)
}

func (m *Model) viewDescription() {
	content := string(m.getSelectedMrRow(table.MergeReqsCols.Desc.Idx, MrTableView))
	m.setResponseContent(content)
	m.PrevView = MrTableView
	m.CurrView = MdView
}

func (m *Model) viewComments() tea.Cmd {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(table.MergeReqsCols.ID.Idx, MrTableView), m.Projects.ProjectID)
	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return "success_comments"
	}
	m.MergeRequests.Comments = m.SetMergeRequestsCommentsModel(r)
	return c
}

func (m *Model) viewPipelines() tea.Cmd {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(table.MergeReqsCols.ID.Idx, MrTableView), m.Projects.ProjectID)
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return "success_pipelines"
	}
	m.MergeRequests.Pipeline = m.SetMergeRequestPipelinesModel(r)
	return c
}

func (m *Model) mergeMR() tea.Cmd {
	statusCode, err := server.MergeMR(m.Projects.ProjectID, m.getSelectedMrRow(table.MergeReqsCols.ID.Idx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			logger.Error(err)
			return err
		}

		switch statusCode {
		case 401:
			return "merge_unauthorized"

		case 405:
			return "merge_method_not_allowed"

		case 409:
			return "merge_error_in_sha"

		case 422:
			return "merge_branch_cant_be_merged"

		}

		return "success_merge"
	}
	return c
}

// Comments Table
// BUG: Not working
func (m *Model) refetchComments() {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(table.MergeReqsCols.CreatedAd.Idx, MrTableView), m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Comments.SetRows(r)
}

func (m *Model) viewCommentContent() {
	content := string(m.getSelectedMrRow(table.CommentsCols.Body.Idx, MrCommentsView))
	m.setResponseContent(content)
	m.PrevView = MrCommentsView
	m.CurrView = MdView
}

func (m *Model) navigateToMrComment() {
	selectedURL := m.getSelectedMrRow(table.MergeReqsCols.URL.Idx, MrTableView)
	commentID := m.getSelectedMrRow(table.CommentsCols.ID.Idx, MrCommentsView)
	exec.Openbrowser(fmt.Sprintf("%s#note_%s", selectedURL, commentID))
}

// Pipelines Table
// BUG: working weird
func (m *Model) refetchPipelines() {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(table.MergeReqsCols.CreatedAd.Idx, MrTableView), m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Pipeline.SetRows(r)
}

func (m *Model) viewPipelineJobs() tea.Cmd {
	r, err := server.GetPipelineJobs(m.Projects.ProjectID, m.getSelectedMrRow(table.PipelinesCols.ID.Idx, MrPipelinesView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return "success_jobs"
	}
	m.MergeRequests.PipelineJobs = m.SetPipelineJobsModel(r)
	return c
}

// Projects List
func (m *Model) viewMergeReqs(window tea.WindowSizeMsg) tea.Cmd {
	s := m.Projects.List.SelectedItem()
	i, ok := s.(item)
	var c tea.Cmd
	if ok {
		m.Projects.ProjectID = i.id
		r, err := server.GetMergeRequests(m.Projects.ProjectID)
		c = func() tea.Msg {
			if err != nil {
				return err
			}

			return "success_mergeReqs"
		}
		m.MergeRequests.List = table.InitModel(table.InitModelParams{
			Rows:      r,
			Colums:    table.GetMergeReqsColums(window.Width - 10),
			StyleFunc: table.StyleIconsColumns(table.Styles(style.Table()), table.MergeReqsIconCols),
		})
	}
	return c
}

// Jobs Table
// BUG: Not working
func (m *Model) refetchJobs() {
	r, err := server.GetPipelineJobs(m.Projects.ProjectID, m.getSelectedMrRow(table.PipelinesCols.ID.Idx, MrPipelinesView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.PipelineJobs.SetRows(r)
}
