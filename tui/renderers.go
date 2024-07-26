package tui

import (
	"fmt"
	"gitlab_tui/tui/components/table"
	"gitlab_tui/tui/components/toast"
	"gitlab_tui/tui/style"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderPaginator() string {
	// TODO: remove help text & add help model
	var b strings.Builder
	b.WriteString("  " + m.Paginator.View())
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return b.String()
}

type renderTableParams struct {
	title    string
	subtitle string
	footer   string
	view     string
}

func (m Model) renderTableView(params renderTableParams) string {
	project := m.Projects.List.SelectedItem().FilterValue()

	var t string
	if params.title != "" && params.subtitle == "" {
		t = fmt.Sprintf("%s - %s", project, params.title)
	} else if params.title != "" && params.subtitle != "" {
		t = fmt.Sprintf("%s - %s %s | %s", project, params.title, params.subtitle, m.MergeRequests.SelectedMr)
	}
	// TODO: create footer status bar (similar to nvim)
	table := lipgloss.JoinVertical(
		0,
		table.TitleStyle.Render(t),
		style.Base.Render(params.view)+"\n",
		m.renderPaginator(),
		style.HelpStyle.Render(params.footer),
	)

	if m.Toast.Show {
		toast := m.Toast.View()
		return lipgloss.JoinVertical(lipgloss.Left, toast, table)
	}

	m.Tabs.Content = table

	return m.Tabs.View()
}

func (m *Model) displayToast(msg string, t toast.ToastType) tea.Cmd {
	m.Toast.Show = true
	m.Toast.Message = getErrorMessage(msg)
	m.Toast.Type = t
	return m.Toast.Init()
}
