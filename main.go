package main

import (
	"fmt"
	"gitlab_tui/api"
	"gitlab_tui/command"
	"gitlab_tui/config"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const (
	RESPONSE_RIGHT_MARGIN      = 2
	useHighPerformanceRenderer = false
)

const (
	tableView currView = iota
	detailView
)

const (
	mrId mrRow = iota
	mrTitle
	mrDesc
	mrAuthor
	mrStatus
	mrDraft
	mrConflicts
	mrUrl
)

type (
	currView uint
	mrRow    uint
)

type model struct {
	table    table.Model
	detail   detail
	currView currView
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			cmds = append(cmds, tea.Quit)
		}

		switch m.currView {
		case tableView:
			switch msg.String() {
			case "x":
				selectedUrl := m.table.SelectedRow()[mrUrl]
				command.Openbrowser(selectedUrl)

			case "enter":
				m.setResponseContent()
				m.currView = detailView

			}

		case detailView:
			switch msg.String() {
			case "backspace":
				m.currView = tableView
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView(m.table.SelectedRow()[mrTitle]))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		cmd = m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	}
	m.table, cmd = m.table.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currView == detailView {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.table.SelectedRow()[mrTitle]), m.detail.model.View(), m.footerView())
	}
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	config.Load(&config.Config)

	r := api.GetMergeRequests()

	columns := []table.Column{
		{Title: "Iid", Width: 4},
		{Title: "Title", Width: 50},
		{Title: "Description", Width: 0},
		{Title: "Author", Width: 15},
		{Title: "Status", Width: 15},
		{Title: "Draft", Width: 10},
		{Title: "Conflicts", Width: 10},
		{Title: "Url", Width: 0},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(r),
		table.WithFocused(true),
		table.WithHeight(len(r)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, currView: tableView}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
