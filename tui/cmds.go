package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
	"gitlab_tui/internal/style"
	tbl "gitlab_tui/tui/components/table"

	tea "github.com/charmbracelet/bubbletea"
)

// Merge Requests Table
func (m *Model) refetchMrList() {
	r, err := server.GetMergeRequests(m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}

	m.MergeRequests.List.SetRows(r)
}

func (m *Model) navigateToMr() {
	selectedURL := m.getSelectedMrRow(tbl.MergeReqsCols.URL.Idx, MrTableView)
	exec.Openbrowser(selectedURL)
}

func (m *Model) viewDescription() {
	content := string(m.getSelectedMrRow(tbl.MergeReqsCols.Desc.Idx, MrTableView))
	m.setResponseContent(content)
	m.PrevView = MrTableView
	m.CurrView = MdView
}

func (m *Model) viewComments() tea.Cmd {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(tbl.MergeReqsCols.Iid.Idx, MrTableView), m.Projects.ProjectID)
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
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(tbl.MergeReqsCols.Iid.Idx, MrTableView), m.Projects.ProjectID)
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return "success_pipelines"
	}
	m.MergeRequests.Pipeline = m.SetMergeRequestPipelinesModel(r)
	return c
}

// Comments Table
func (m *Model) refetchComments() {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(tbl.MergeReqsCols.Iid.Idx, MrTableView), m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Comments.SetRows(r)
}

func (m *Model) viewCommentContent() {
	content := string(m.getSelectedMrRow(tbl.CommentsCols.Body.Idx, MrCommentsView))
	m.setResponseContent(content)
	m.PrevView = MrCommentsView
	m.CurrView = MdView
}

func (m *Model) navigateToMrComment() {
	selectedURL := m.getSelectedMrRow(tbl.MergeReqsCols.URL.Idx, MrTableView)
	commentID := m.getSelectedMrRow(tbl.CommentsCols.ID.Idx, MrCommentsView)
	exec.Openbrowser(fmt.Sprintf("%s#note_%s", selectedURL, commentID))
}

// Pipelines Table
func (m *Model) refetchPipelines() {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(tbl.MergeReqsCols.Iid.Idx, MrTableView), m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Pipeline.SetRows(r)
}

func (m *Model) navigateToPipeline() {
	selectedURL := m.getSelectedMrRow(tbl.PipelinesCols.URL.Idx, MrPipelinesView)
	exec.Openbrowser(selectedURL)
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
		m.MergeRequests.List = tbl.InitModel(tbl.InitModelParams{
			Rows:      r,
			Colums:    tbl.GetMergeReqsColums(window.Width - 10),
			StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.MergeReqsIconCols),
		})
	}
	return c
}
