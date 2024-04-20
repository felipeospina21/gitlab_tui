package tui

import (
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/server"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) refetchPipelines() (Model, tea.Cmd) {
	var newM Model
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	if err != nil {
		return newM, func() tea.Msg {
			return err
		}
	}
	t := InitMergeRequestsListTable(r, 155)
	newM = m.UpdateMergeRequestsModel(m.MergeRequests.List, m.MergeRequests.Comments, t)

	return newM, func() tea.Msg {
		return nil
	}
}

func (m *Model) navigateToPipeline() {
	selectedURL := m.getSelectedMrRow(mrPipelinesURLIdx, MrPipelinesView)
	exec.Openbrowser(selectedURL)
}
