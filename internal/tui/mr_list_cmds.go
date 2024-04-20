package tui

import (
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) refetchMrList() (Model, tea.Cmd) {
	r, err := server.GetMergeRequests()
	if err != nil {
		logger.Error(err)
	}
	// TODO: Do something with Error msg
	// TODO: Create enum for queries status
	// if err != nil {
	// 	return nil, func() tea.Msg {
	// 		return "fetch_mr_table_error"
	// 	}
	// }

	t := InitMergeRequestsListTable(r, 155)
	newM := m.UpdateMergeRequestsModel(t, table.Model{}, table.Model{})

	return newM, func() tea.Msg {
		return nil
	}
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
