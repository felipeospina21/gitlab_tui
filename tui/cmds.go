package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/server"
	"gitlab_tui/tui/components/table"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) openInBrowser(tableColIdx table.TableColIndex, view views) {
	selectedURL := m.getSelectedRow(tableColIdx, view)
	exec.Openbrowser(selectedURL)
}

func (m *Model) toggleSidePanel() {
	m.isSidePanelOpen = !m.isSidePanelOpen
	if m.isSidePanelOpen {
		m.CurrView = HomeView
	} else {
		m.CurrView = m.PrevView
	}
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
	content := string(m.getSelectedRow(table.MergeReqsCols.Desc.Idx, MrTableView))
	m.setResponseContent(content)
	m.PrevView = MrTableView
	m.CurrView = MdView
}

func (m *Model) viewComments() tea.Cmd {
	r, err := server.GetMergeRequestComments(m.Projects.ProjectID, m.getSelectedRow(table.MergeReqsCols.ID.Idx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return SuccessMessage.CommentsFetch
	}
	m.MergeRequests.Comments = m.SetMergeRequestsCommentsModel(r)
	return c
}

func (m *Model) viewPipelines() tea.Cmd {
	r, err := server.GetMergeRequestPipelines(m.Projects.ProjectID, m.getSelectedRow(table.MergeReqsCols.ID.Idx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return SuccessMessage.PipelinesFetch
	}
	m.MergeRequests.Pipeline = m.SetMergeRequestPipelinesModel(r)
	return c
}

func (m *Model) mergeMR() tea.Cmd {
	_, err := server.MergeMR(m.Projects.ProjectID, m.getSelectedRow(table.MergeReqsCols.ID.Idx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return SuccessMessage.Merge
	}
	return c
}

// Comments Table
func (m *Model) refetchComments() {
	r, err := server.GetMergeRequestComments(m.Projects.ProjectID, m.getSelectedRow(table.MergeReqsCols.ID.Idx, MrTableView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Comments.SetRows(r)
}

func (m *Model) viewCommentContent() {
	content := string(m.getSelectedRow(table.CommentsCols.Body.Idx, MrCommentsView))
	m.setResponseContent(content)
	m.PrevView = MrCommentsView
	m.CurrView = MdView
}

func (m *Model) navigateToMrComment() {
	selectedURL := m.getSelectedRow(table.MergeReqsCols.URL.Idx, MrTableView)
	commentID := m.getSelectedRow(table.CommentsCols.ID.Idx, MrCommentsView)
	exec.Openbrowser(fmt.Sprintf("%s#note_%s", selectedURL, commentID))
}

// Pipelines Table
func (m *Model) refetchPipelines() {
	r, err := server.GetMergeRequestPipelines(m.Projects.ProjectID, m.getSelectedRow(table.MergeReqsCols.ID.Idx, MrTableView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Pipeline.SetRows(r)
}

func (m *Model) viewPipelineJobs() tea.Cmd {
	r, err := server.GetPipelineJobs(m.Projects.ProjectID, m.getSelectedRow(table.PipelinesCols.ID.Idx, MrPipelinesView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return SuccessMessage.JobsFetch
	}
	m.MergeRequests.PipelineJobs = m.SetPipelineJobsModel(r)
	return c
}

// Projects List
func (m *Model) viewMergeReqs(window tea.WindowSizeMsg) tea.Cmd {
	s := m.Projects.List.SelectedItem()
	i, ok := s.(Item)
	var c tea.Cmd
	if ok {
		m.Projects.ProjectID = i.ID
		r, err := server.GetMergeRequests(m.Projects.ProjectID)
		c = func() tea.Msg {
			if err != nil {
				return err
			}

			return SuccessMessage.MRFetch
		}
		m.MergeRequests.List = table.InitModel(table.InitModelParams{
			Rows:      r,
			Colums:    table.GetMergeReqsColums(window.Width - 10),
			StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.MergeReqsIconCols),
		})
	}
	return c
}

// Jobs Table
func (m *Model) refetchJobs() {
	r, err := server.GetPipelineJobs(m.Projects.ProjectID, m.getSelectedRow(table.PipelinesCols.ID.Idx, MrPipelinesView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.PipelineJobs.SetRows(r)
}

// Issues List Table
func (m *Model) viewIssues() tea.Cmd {
	r, pages, err := server.GetIssues(m.Projects.ProjectID, "")

	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return SuccessMessage.IssuesList
	}
	m.Issues.List = m.SetIssuesListModel(r)
	m.Issues.PrevPage = pages.Prev
	m.Issues.NexPage = pages.Next
	m.Paginator.TotalPages = pages.Total

	return c
}

func (m *Model) getIssuesNextPage() tea.Cmd {
	r, pages, err := server.GetIssues(m.Projects.ProjectID, m.Issues.NexPage)

	return m.issuesPageCmd(r, pages, err, "")
}

func (m *Model) getIssuesPrevPage() tea.Cmd {
	r, pages, err := server.GetIssues(m.Projects.ProjectID, m.Issues.PrevPage)

	return m.issuesPageCmd(r, pages, err, "")
}

func (m *Model) issuesPageCmd(r []table.Row, pages server.Pages, err error, msg string) tea.Cmd {
	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return msg
	}
	m.Issues.List.SetRows(r)
	m.Issues.PrevPage = pages.Prev
	m.Issues.NexPage = pages.Next

	return c
}

func (m *Model) viewIssueDescription() {
	content := string(m.getSelectedRow(table.IssuesListCols.Desc.Idx, IssuesListView))
	m.setResponseContent(content)
	m.PrevView = IssuesListView
	m.CurrView = MdView
}
