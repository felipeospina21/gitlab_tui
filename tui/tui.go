package tui

import (
	"fmt"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components"
	tbl "gitlab_tui/tui/components/table"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	views uint
)

type Model struct {
	Projects      ProjectsModel
	MergeRequests MergeRequestsModel
	Md            MdModel
	CurrView      views
	PrevView      views
	Title         string
	Window        tea.WindowSizeMsg
	Help          components.Help
}

const (
	MrTableView views = iota
	MrCommentsView
	MrPipelinesView
	MdView
	ProjectsView
	TabsView
)

const useHighPerformanceRenderer = false

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.CurrView {
		case MrTableView:
			switch {
			case key.Matches(msg, m.MergeRequests.ListKeys.Help):
				m.Help.Model.ShowAll = !m.Help.Model.ShowAll

			case key.Matches(msg, m.MergeRequests.ListKeys.OpenInBrowser):
				m.navigateToMr()

			case key.Matches(msg, m.MergeRequests.ListKeys.Comments):
				c := m.viewComments()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.ListKeys.Pipelines):
				c := m.viewPipelines()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.ListKeys.Description):
				m.viewDescription()

			case key.Matches(msg, m.MergeRequests.ListKeys.Refetch):
				m.refetchMrList()

			case key.Matches(msg, m.MergeRequests.ListKeys.NavigateBack):
				m.CurrView = ProjectsView

			}
			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)
		}

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
		case ProjectsView:
			switch msg.String() {
			case "enter":
				c := m.viewMergeReqs(m.Window)
				cmds = append(cmds, c)
			}
			m.Projects.List, cmd = m.Projects.List.Update(msg)

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

		cmd = m.setViewportViewSize(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		switch m.CurrView {
		case MrTableView:
			return m.resizeMrTable(msg)

		case MrCommentsView:
			// BUG: After resize, comments table loses its width
			return m.resizeMrCommentsTable(msg)

		case MrPipelinesView:
			return m.resizeMrPipelinesTable(msg)

		case MdView:
			m.resizeMdView(msg)

		}

	case string:
		if msg == "success_mergeReqs" {
			m.CurrView = MrTableView
		}
		if msg == "success_comments" {
			m.CurrView = MrCommentsView
			m.setSelectedMr()
		}
		if msg == "success_pipelines" {
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
	case ProjectsView:
		return style.ListTitleStyle.Render(m.Projects.List.View())

	case MdView:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.MergeRequests.List.SelectedRow()[tbl.MergeReqsCols.Title.Idx]), m.Md.Viewport.View(), m.footerView())

	case MrCommentsView:
		return m.renderTableView(m.MergeRequests.Comments.View(), "Comments", m.Help.Model.View(m.MergeRequests.CommentsKeys))

	case MrPipelinesView:
		return m.renderTableView(m.MergeRequests.Pipeline.View(), "Pipelines", m.Help.Model.View(m.MergeRequests.PipelineKeys))

	default:
		return m.renderTableView(m.MergeRequests.List.View(), "", m.Help.Model.View(m.MergeRequests.ListKeys))

	}
}

func (m Model) renderTableView(view string, title string, footer string) string {
	project := m.Projects.List.SelectedItem().FilterValue()

	var t string
	if title == "" {
		t = fmt.Sprintf("%s - Merge Requests", project)
	} else {
		t = fmt.Sprintf("%s - Merge Request %s | %s", project, title, m.MergeRequests.SelectedMr)
	}
	return lipgloss.JoinVertical(
		0,
		style.TableTitle.Render(t),
		style.Base.Render(view)+"\n",
		style.HelpStyle.Render(footer),
	)
}

func (m Model) getSelectedMrRow(idx tbl.TableColIndex, view views) string {
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
	m.MergeRequests.SelectedMr = m.getSelectedMrRow(tbl.MergeReqsCols.Title.Idx, MrTableView)
}
