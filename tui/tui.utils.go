package tui

import (
	"gitlab_tui/tui/components/table"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) resizeMrTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.List.SetWidth(msg.Width)

	t := table.InitModel(table.InitModelParams{
		Rows:      m.MergeRequests.List.Rows(),
		Colums:    table.GetMergeReqsColums(msg.Width - 10),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.MergeReqsIconCols),
	})

	newM := m.UpdateModel(
		t,
		table.Model{},
		table.Model{},
		m.Projects.List,
		table.Model{},
	)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrCommentsTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Comments.SetWidth(msg.Width)

	t := table.InitModel(table.InitModelParams{
		Rows:      m.MergeRequests.Comments.Rows(),
		Colums:    table.GetCommentsColums(msg.Width - 10),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.CommentsIconCols),
	})

	newM := m.UpdateModel(
		m.MergeRequests.List,
		t,
		m.MergeRequests.Pipeline,
		m.Projects.List,
		m.MergeRequests.PipelineJobs,
	)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeMrPipelinesTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.Pipeline.SetWidth(msg.Width)

	t := table.InitModel(table.InitModelParams{
		Rows:      m.MergeRequests.Pipeline.Rows(),
		Colums:    table.GetPipelinesColums(msg.Width - 10),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.PipelinesIconCols),
	})

	newM := m.UpdateModel(
		m.MergeRequests.List,
		m.MergeRequests.Comments,
		t,
		m.Projects.List,
		m.MergeRequests.PipelineJobs,
	)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizePipelineJobsTable(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.MergeRequests.PipelineJobs.SetWidth(msg.Width)

	t := table.InitModel(table.InitModelParams{
		Rows:   m.MergeRequests.PipelineJobs.Rows(),
		Colums: table.GetPipelineJobsColums(msg.Width - 10),
	})

	newM := m.UpdateModel(
		m.MergeRequests.List,
		m.MergeRequests.Comments,
		m.MergeRequests.Pipeline,
		m.Projects.List,
		t,
	)

	return newM, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m *Model) resizeProjectsList(msg tea.WindowSizeMsg) {
	l := InitProjectsList()

	newM := m.UpdateModel(
		m.MergeRequests.List,
		m.MergeRequests.Comments,
		m.MergeRequests.Pipeline,
		l,
		m.MergeRequests.PipelineJobs,
	)

	m.Projects.List = newM.Projects.List
	m.Projects.List.SetSize(msg.Width, msg.Height)
}

func (m *Model) resizeMdView(msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(m.headerView(m.getSelectedRow(table.MergeReqsCols.Title.Idx, MainTableView)))
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	m.Md.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)

	var content string
	switch m.PrevView {
	case MainTableView:
		content = m.getSelectedRow(table.MergeReqsCols.Desc.Idx, MainTableView)

	case MrCommentsView:
		content = m.getSelectedRow(table.CommentsCols.Body.Idx, MrCommentsView)

	default:
		content = "Model not selected"
	}
	m.setResponseContent(content)
}

func (m Model) UpdateModel(listModel table.Model, commentsModel table.Model, pipelinesModel table.Model, projectsModel list.Model, jobsModel table.Model) Model {
	newM := Model{
		MergeRequests: MergeRequestsModel{
			List:         listModel,
			Comments:     commentsModel,
			Pipeline:     pipelinesModel,
			PipelineJobs: jobsModel,
		},
		Projects: ProjectsModel{
			List:      projectsModel,
			ProjectID: m.Projects.ProjectID,
		},
		CurrView: m.CurrView,
		Md:       m.Md,
		Help:     m.Help,
		Window:   tea.WindowSizeMsg{Width: m.Window.Width, Height: m.Window.Height},
	}

	return newM
}
