package tui

import (
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
)

func (m *Model) refetchPipelines() {
	r, err := server.GetMergeRequestPipelines(m.getSelectedMrRow(mergeReqsCols.iid.idx, MrTableView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Pipeline.SetRows(r)
}

func (m *Model) navigateToPipeline() {
	selectedURL := m.getSelectedMrRow(pipelinesCols.url.idx, MrPipelinesView)
	exec.Openbrowser(selectedURL)
}
