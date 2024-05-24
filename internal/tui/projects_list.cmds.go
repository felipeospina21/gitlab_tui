package tui

import (
	"gitlab_tui/internal/server"

	tea "github.com/charmbracelet/bubbletea"
)

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
		m.MergeRequests.List = InitMergeRequestsListTable(r, window.Width-10)
	}
	return c
}
