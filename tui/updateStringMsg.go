package tui

import (
	"gitlab_tui/tui/components/toast"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateStringMsg(msg string) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg {
	case SuccessMessage.MRFetch:
		m.CurrView = MainTableView
		m.isSidePanelOpen = false

	case SuccessMessage.DiscussionsFetch:
		m.CurrView = MrDiscussionsView
		m.setSelectedMr()

	case SuccessMessage.PipelinesFetch:
		m.CurrView = MrPipelinesView
		m.setSelectedMr()

	case SuccessMessage.JobsFetch:
		m.CurrView = JobsView
		// m.MergeRequests.PipelineJobs.SelectedRow()[table.PipelineJobsCols.Name.Idx]
		// m.setSelectedMr()

	case SuccessMessage.Merge:
		// TODO: change message for api response
		cmd = m.displayToast("Successfully Merged", toast.Success)
		// cmds = append(cmds, cmd)

	case SuccessMessage.IssuesList:
		m.CurrView = MainTableView
		m.Issues.HasData = true

	}
	return *m, cmd
}
