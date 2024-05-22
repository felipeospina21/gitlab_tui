package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) resizeMrTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.List.SetWidth(msg.Width)
	t := InitMergeRequestsListTable(m.MergeRequests.List.Rows(), msg.Width-10)
	newM := m.UpdateMergeRequestsModel(t, table.Model{}, table.Model{})

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrCommentsTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Comments.SetWidth(msg.Width)
	t := InitMergeRequestsListTable(m.MergeRequests.Comments.Rows(), msg.Width-10)
	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t, m.MergeRequests.Pipeline)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrPipelinesTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Pipeline.SetWidth(msg.Width)
	t := InitMergeRequestsListTable(m.MergeRequests.Pipeline.Rows(), msg.Width-10)
	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, m.MergeRequests.Comments, t)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMdView(msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(m.headerView(m.getSelectedMrRow(mergeReqsCols.title.idx, MrTableView)))
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	m.Md.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)

	var content string
	switch m.PrevView {
	case MrTableView:
		content = m.getSelectedMrRow(mergeReqsCols.desc.idx, MrTableView)

	case MrCommentsView:
		content = m.getSelectedMrRow(commentsCols.body.idx, MrCommentsView)

	default:
		content = "Model not selected"
	}
	m.setResponseContent(content)
}
