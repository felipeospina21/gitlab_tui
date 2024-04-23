package tui

import (
	"fmt"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/style"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	views         uint
	tableColIndex uint
)

type Model struct {
	MergeRequests MergeRequestsModel
	Md            MdModel
	CurrView      views
	PrevView      views
	Title         string
	Window        tea.WindowSizeMsg
}

const (
	MrTableView views = iota
	MrCommentsView
	MrPipelinesView
	MdView
	TabsView
)

const useHighPerformanceRenderer = false

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Global commands
		switch msg.String() {
		case "esc":
			if m.MergeRequests.List.Focused() {
				m.MergeRequests.List.Blur()
			} else {
				m.MergeRequests.List.Focus()
			}
		case "q", "ctrl+c":
			cmds = append(cmds, tea.Quit)

		case "tab":
			m.CurrView = MdView
		}

		switch m.CurrView {
		case MrTableView:
			switch msg.String() {
			case "r":
				m.refetchMrList()

			case "x":
				m.navigateToMr()

			case "enter":
				m.viewDescription()

			case "c":
				c := m.viewComments()
				cmds = append(cmds, c)

			case "p":
				c := m.viewPipelines()
				cmds = append(cmds, c)
			}

			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)

		case MrCommentsView:
			switch msg.String() {
			case "r":
				m.refetchComments()

			case "x":
				m.navigateToMrComment()

			case "enter":
				m.viewCommentContent()

			case "backspace":
				m.CurrView = MrTableView

			}

			m.MergeRequests.Comments, cmd = m.MergeRequests.Comments.Update(msg)

		case MrPipelinesView:
			switch msg.String() {
			case "r":
				m.refetchPipelines()

			case "x":
				m.navigateToPipeline()

			case "backspace":
				m.CurrView = MrTableView

			}

			m.MergeRequests.Pipeline, cmd = m.MergeRequests.Pipeline.Update(msg)

		case MdView:
			switch msg.String() {
			case "backspace":
				m.CurrView = m.PrevView
			}

			m.Md.Viewport, cmd = m.Md.Viewport.Update(msg)
		}

	case tea.WindowSizeMsg:
		m.Window = msg
		logger.Debug("window", func() {
			log.Print(m.Window)
			log.Print(msg)
		})
		cmd = m.setViewportViewSize(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		switch m.CurrView {
		case MrTableView:
			return m.resizeMrTable(msg)

		case MrCommentsView:
			return m.resizeMrCommentsTable(msg)

		case MrPipelinesView:
			return m.resizeMrPipelinesTable(msg)

		case MdView:
			m.resizeMdView(msg)

		}

	case string:
		if msg == "success_comments" {
			m.MergeRequests.Comments.SetStyles(style.Table)
			m.CurrView = MrCommentsView
			m.setSelectedMr()
		}
		if msg == "success_pipelines" {
			m.MergeRequests.Pipeline.SetStyles(style.Table)
			m.CurrView = MrPipelinesView
			m.setSelectedMr()
		}

	case error:
		logger.Debug("error", func() {
			log.Println(msg)
		})

	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.CurrView {
	case MdView:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.MergeRequests.List.SelectedRow()[mergeReqsTitleIdx]), m.Md.Viewport.View(), m.footerView())

	case MrCommentsView:
		return m.renderTableView(m.MergeRequests.Comments.View(), "Comments")

	case MrPipelinesView:
		return m.renderTableView(m.MergeRequests.Pipeline.View(), "Pipelines")

	default:
		return m.renderTableView(m.MergeRequests.List.View(), "")

	}
}

func (m Model) renderTableView(view string, title string) string {
	titleStyle := lipgloss.NewStyle().Margin(2, 0, 1, 2).Foreground(lipgloss.Color("51"))

	var t string
	if title == "" {
		t = "Merge Requests"
	} else {
		t = fmt.Sprintf("Merge Request %s | %s", title, m.MergeRequests.SelectedMr)
	}
	return lipgloss.JoinVertical(0, titleStyle.Render(t), style.Base.Render(view)+"\n")
}

func (m Model) getSelectedMrRow(idx tableColIndex, view views) string {
	switch view {
	case MrTableView:
		return m.MergeRequests.List.SelectedRow()[idx]

	case MrCommentsView:
		return m.MergeRequests.Comments.SelectedRow()[idx]

	case MrPipelinesView:
		return m.MergeRequests.Pipeline.SelectedRow()[idx]

	default:
		return ""

	}
}

func (m *Model) setSelectedMr() {
	m.MergeRequests.SelectedMr = m.getSelectedMrRow(mergeReqsTitleIdx, MrTableView)
}
