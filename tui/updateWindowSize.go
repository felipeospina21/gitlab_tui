package tui

import (
	"gitlab_tui/internal/style"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateWindowSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.Window = msg
	m.Toast.Progress.Width = msg.Width - 4
	m.Toast.Width = msg.Width - 4

	cmd = m.setViewportViewSize(msg)

	switch m.CurrView {
	case MrTableView:
		return m.resizeMrTable(msg)

	case MrCommentsView:
		return m.resizeMrCommentsTable(msg)

	case MrPipelinesView:
		return m.resizeMrPipelinesTable(msg)

	case JobsView:
		return m.resizePipelineJobsTable(msg)

	case MdView:
		m.resizeMdView(msg)

	case ProjectsView:
		h, v := style.ListItemStyle.GetFrameSize()
		m.Projects.List.SetSize(msg.Width-h, (msg.Height-v)/2)

	}

	return *m, cmd
}
