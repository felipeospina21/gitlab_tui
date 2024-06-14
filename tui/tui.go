package tui

import (
	"fmt"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components"
	"gitlab_tui/tui/components/table"
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
	JobsView
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
				m.openInBrowser(table.MergeReqsCols.URL.Idx, MrTableView)

			case key.Matches(msg, m.MergeRequests.ListKeys.Comments):
				c := m.viewComments()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.ListKeys.Pipelines):
				c := m.viewPipelines()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.ListKeys.Merge):
				c := m.mergeMR()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.ListKeys.Description):
				m.viewDescription()

			case key.Matches(msg, m.MergeRequests.ListKeys.Refetch):
				m.refetchMrList()

			case key.Matches(msg, m.MergeRequests.ListKeys.NavigateBack):
				m.CurrView = ProjectsView

			}
			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)

		case MrCommentsView:
			switch {
			case key.Matches(msg, m.MergeRequests.CommentsKeys.Refetch):
				m.refetchComments()

			case key.Matches(msg, m.MergeRequests.CommentsKeys.OpenInBrowser):
				m.navigateToMrComment()

			case key.Matches(msg, m.MergeRequests.CommentsKeys.Description):
				m.viewCommentContent()

			case key.Matches(msg, m.MergeRequests.CommentsKeys.NavigateBack):
				m.CurrView = MrTableView

			}
			m.MergeRequests.Comments, cmd = m.MergeRequests.Comments.Update(msg)

		case MrPipelinesView:
			switch {
			case key.Matches(msg, m.MergeRequests.PipelineKeys.Jobs):
				c := m.viewPipelineJobs()
				cmds = append(cmds, c)

			case key.Matches(msg, m.MergeRequests.PipelineKeys.Refetch):
				m.refetchPipelines()

			case key.Matches(msg, m.MergeRequests.PipelineKeys.OpenInBrowser):
				m.openInBrowser(table.PipelinesCols.URL.Idx, MrPipelinesView)

			case key.Matches(msg, m.MergeRequests.PipelineKeys.NavigateBack):
				m.CurrView = MrTableView

			}
			m.MergeRequests.Pipeline, cmd = m.MergeRequests.Pipeline.Update(msg)

		case JobsView:
			switch {
			case key.Matches(msg, m.MergeRequests.JobsKeys.NavigateBack):
				m.CurrView = MrPipelinesView

			case key.Matches(msg, m.MergeRequests.JobsKeys.OpenInBrowser):
				m.openInBrowser(table.PipelineJobsCols.URL.Idx, JobsView)

			case key.Matches(msg, m.MergeRequests.JobsKeys.Refetch):
				m.refetchJobs()

			}
			m.MergeRequests.PipelineJobs, cmd = m.MergeRequests.PipelineJobs.Update(msg)
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
			return m.resizeMrCommentsTable(msg)

		case MrPipelinesView:
			return m.resizeMrPipelinesTable(msg)

		case JobsView:
			return m.resizePipelineJobsTable(msg)

		case MdView:
			m.resizeMdView(msg)

		}

	case string:
		switch msg {
		case "success_mergeReqs":
			m.CurrView = MrTableView

		case "success_comments":
			m.CurrView = MrCommentsView
			m.setSelectedMr()

		case "success_pipelines":
			m.CurrView = MrPipelinesView
			m.setSelectedMr()

		case "success_jobs":
			m.CurrView = JobsView
			// m.MergeRequests.PipelineJobs.SelectedRow()[table.PipelineJobsCols.Name.Idx]
			// m.setSelectedMr()

		case "merge_unauthorized":
		case "merge_branch_cant_be_merged":
		case "merge_method_not_allowed":
		case "merge_error_in_sha":
			m.setResponseContent("Error: " + msg)
			m.CurrView = MdView

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
		return style.ListItemStyle.Render(m.Projects.List.View())

	case MdView:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.MergeRequests.List.SelectedRow()[table.MergeReqsCols.Title.Idx]), m.Md.Viewport.View(), m.footerView())

	case MrCommentsView:
		return m.renderTableView(m.MergeRequests.Comments.View(), "Comments", m.Help.Model.View(m.MergeRequests.CommentsKeys))

	case MrPipelinesView:
		return m.renderTableView(m.MergeRequests.Pipeline.View(), "Pipelines", m.Help.Model.View(m.MergeRequests.PipelineKeys))

	case JobsView:
		return m.renderTableView(m.MergeRequests.PipelineJobs.View(), "Jobs", m.Help.Model.View(m.MergeRequests.JobsKeys))

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

func (m Model) getSelectedRow(idx table.TableColIndex, view views) string {
	switch view {
	case MrTableView:
		return m.MergeRequests.List.SelectedRow()[idx]

	case MrCommentsView:
		return m.MergeRequests.Comments.SelectedRow()[idx]

	case MrPipelinesView:
		return m.MergeRequests.Pipeline.SelectedRow()[idx]

	case JobsView:
		return m.MergeRequests.PipelineJobs.SelectedRow()[idx]

	default:
		return ""

	}
}

func (m *Model) setSelectedMr() {
	m.MergeRequests.SelectedMr = m.getSelectedRow(table.MergeReqsCols.Title.Idx, MrTableView)
}
