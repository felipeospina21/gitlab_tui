package tui

import (
	"gitlab_tui/internal/style"
	tbl "gitlab_tui/tui/components/table"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) resizeMrTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.List.SetWidth(msg.Width)

	t := tbl.InitModel(tbl.InitModelParams{
		Rows:      m.MergeRequests.List.Rows(),
		Colums:    tbl.GetMergeReqsColums(msg.Width - 10),
		StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.MergeReqsIconCols),
	})

	newM := m.UpdateMergeRequestsModel(t, table.Model{}, table.Model{}, m.Projects.List)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrCommentsTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Comments.SetWidth(msg.Width)

	t := tbl.InitModel(tbl.InitModelParams{
		Rows:      m.MergeRequests.Comments.Rows(),
		Colums:    tbl.GetCommentsColums(msg.Width - 10),
		StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.CommentsIconCols),
	})

	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t, m.MergeRequests.Pipeline, m.Projects.List)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrPipelinesTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Pipeline.SetWidth(msg.Width)

	t := tbl.InitModel(tbl.InitModelParams{
		Rows:      m.MergeRequests.Pipeline.Rows(),
		Colums:    tbl.GetPipelinesColums(msg.Width - 10),
		StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.PipelinesIconCols),
	})

	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, m.MergeRequests.Comments, t, m.Projects.List)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMdView(msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(m.headerView(m.getSelectedMrRow(tbl.MergeReqsCols.Title.Idx, MrTableView)))
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	m.Md.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)

	var content string
	switch m.PrevView {
	case MrTableView:
		content = m.getSelectedMrRow(tbl.MergeReqsCols.Desc.Idx, MrTableView)

	case MrCommentsView:
		content = m.getSelectedMrRow(tbl.CommentsCols.Body.Idx, MrCommentsView)

	default:
		content = "Model not selected"
	}
	m.setResponseContent(content)
}

func (m *Model) resizeProjectsList(msg tea.WindowSizeMsg) {
	l := InitProjectsList()
	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, m.MergeRequests.Comments, m.MergeRequests.Pipeline, l)
	m.Projects.List = newM.Projects.List
	m.Projects.List.SetSize(msg.Width, msg.Height)
}
