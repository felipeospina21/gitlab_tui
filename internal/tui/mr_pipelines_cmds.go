package tui

import (
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
)

func (m *Model) refetchPipelines() {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Pipeline.SetRows(r)
}

func (m *Model) navigateToPipeline() {
	selectedURL := m.getSelectedMrRow(mrPipelinesURLIdx, MrPipelinesView)
	exec.Openbrowser(selectedURL)
}
