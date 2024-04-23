package tui

import (
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) refetchMrList() {
	r, err := server.GetMergeRequests()
	if err != nil {
		logger.Error(err)
	}

	m.MergeRequests.List.SetRows(r)
}

func (m *Model) navigateToMr() {
	selectedURL := m.getSelectedMrRow(mergeReqsURLIdx, MrTableView)
	exec.Openbrowser(selectedURL)
}

func (m *Model) viewDescription() {
	content := string(m.getSelectedMrRow(mergeReqsDescIdx, MrTableView))
	m.setResponseContent(content)
	m.PrevView = MrTableView
	m.CurrView = MdView
}

func (m *Model) viewComments() tea.Cmd {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}

		return "success_comments"
	}
	m.MergeRequests.Comments = SetMergeRequestsCommentsModel(r)
	return c
}

func (m *Model) viewPipelines() tea.Cmd {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	c := func() tea.Msg {
		if err != nil {
			return err
		}
		return "success_pipelines"
	}
	m.MergeRequests.Pipeline = SetMergeRequestPipelinesModel(r)
	return c
}
