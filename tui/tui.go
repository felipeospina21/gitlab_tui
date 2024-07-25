package tui

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components"
	"gitlab_tui/tui/components/table"
	"gitlab_tui/tui/components/tabs"
	"gitlab_tui/tui/components/toast"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	views uint
)

type Model struct {
	Projects      ProjectsModel
	MergeRequests MergeRequestsModel
	Issues        IssuesModel
	Md            MdModel
	CurrView      views
	PrevView      views
	Title         string
	Window        tea.WindowSizeMsg
	Help          components.Help
	Keys          GlobalKeyMap
	Toast         toast.Model
	Tabs          tabs.Model
	Paginator     paginator.Model
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

type SuccessMsg struct {
	MRFetch        string
	CommentsFetch  string
	PipelinesFetch string
	JobsFetch      string
	Merge          string
	ReloadEnv      string
	IssuesList     string
}

var SuccessMessage = SuccessMsg{
	MRFetch:        "success_mr_fetch",
	CommentsFetch:  "success_comments_fetch",
	PipelinesFetch: "success_pipelines_fetch",
	JobsFetch:      "success_jobs_fetch",
	Merge:          "success_merge",
	ReloadEnv:      "success_env_reload",
	IssuesList:     "success_issues_list",
}

const (
	useHighPerformanceRenderer = false
	toastInterval              = 10
)

type tickMsg time.Time

func (m Model) Init() tea.Cmd {
	return m.Toast.Init()
	// return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Global commands
		switch {
		case key.Matches(msg, GlobalKeys.Quit):
			cmds = append(cmds, tea.Quit)

		case key.Matches(msg, GlobalKeys.ThrowError):
			cmds = append(cmds, func() tea.Msg {
				return errors.New("mocked")
			})
		}

		// Tabs
		switch m.Tabs.ActiveTab {
		case tabs.MergeRequests:
			if msg.String() == "tab" {
				if strings.TrimSpace(m.Issues.List.View()) == "" {
					cmds = append(cmds, m.viewIssues())
				}
			}

		case tabs.Issues:
			if msg.String() == "right" {
				cmds = append(cmds, m.getIssuesNextPage())
			}

		}

		// Views commands
		switch m.CurrView {
		case ProjectsView:
			switch {
			case key.Matches(msg, ProjectsKeys.ViewMRs):
				c := m.viewMergeReqs(m.Window)
				cmds = append(cmds, c)
			}
			m.Projects.List, cmd = m.Projects.List.Update(msg)

		case MdView:
			switch {
			case key.Matches(msg, MdKeys.NavigateBack):
				m.CurrView = m.PrevView
			}
			m.Md.Viewport, cmd = m.Md.Viewport.Update(msg)

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

	case tea.WindowSizeMsg:
		m.Window = msg
		m.Toast.Progress.Width = msg.Width - 4
		m.Toast.Width = msg.Width - 4

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

		case ProjectsView:
			h, v := style.ListItemStyle.GetFrameSize()
			m.Projects.List.SetSize(msg.Width-h, (msg.Height-v)/2)

		}

	case string:
		switch msg {
		case SuccessMessage.MRFetch:
			m.CurrView = MrTableView

		case SuccessMessage.CommentsFetch:
			m.CurrView = MrCommentsView
			m.setSelectedMr()

		case SuccessMessage.PipelinesFetch:
			m.CurrView = MrPipelinesView
			m.setSelectedMr()

		case SuccessMessage.JobsFetch:
			m.CurrView = JobsView
			// m.MergeRequests.PipelineJobs.SelectedRow()[table.PipelineJobsCols.Name.Idx]
			// m.setSelectedMr()

		case SuccessMessage.Merge:
			// TODO: change message for api response
			cmd = m.displayToast("Successfully Merged", toast.Success)
			cmds = append(cmds, cmd)

		case SuccessMessage.IssuesList:
			// unused currently

		}

	case error:
		cmd = m.displayToast(msg.Error(), toast.Error)
		cmds = append(cmds, cmd)

		lh, lv := style.ListItemStyle.GetFrameSize()
		nh, nv := toast.ErrorStyle(m.Window.Height, m.Window.Width).GetFrameSize()

		h := (lh + nh) * 2
		v := (lv + nv) * 2

		m.Projects.List.SetSize(m.Window.Width-h, m.Window.Height-v)

	}

	toastModel, toastCmd := m.Toast.Update(msg)
	m.Toast = toastModel.(toast.Model)

	tabsModel, tabsCmd := m.Tabs.Update(msg)
	m.Tabs = tabsModel.(tabs.Model)

	issuesModel, issuesCmd := m.Issues.List.Update(msg)
	m.Issues.List = issuesModel

	paginatorModel, paginatorCmd := m.Paginator.Update(msg)
	m.Paginator = paginatorModel

	cmds = append(cmds, cmd, toastCmd, tabsCmd, issuesCmd, paginatorCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.Tabs.ActiveTab {
	case tabs.MergeRequests:
		switch m.CurrView {
		case ProjectsView:
			projects := style.ListItemStyle.Render(m.Projects.List.View())

			if m.Toast.Show {
				toast := m.Toast.View()
				return lipgloss.JoinVertical(lipgloss.Left, toast, projects)
			}
			return projects

		case MdView:
			return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.MergeRequests.List.SelectedRow()[table.MergeReqsCols.Title.Idx]), m.Md.Viewport.View(), m.footerView())

		case MrCommentsView:
			return m.renderTableView(renderTableParams{
				title:    "Merge Requests",
				subtitle: "Comments",
				footer:   m.Help.Model.View(m.MergeRequests.CommentsKeys),
				view:     m.MergeRequests.Comments.View(),
			})

		case MrPipelinesView:
			return m.renderTableView(renderTableParams{
				title:    "Merge Requests",
				subtitle: "Pipelines",
				footer:   m.Help.Model.View(m.MergeRequests.PipelineKeys),
				view:     m.MergeRequests.Pipeline.View(),
			})

		case JobsView:
			return m.renderTableView(renderTableParams{
				title:    "Merge Requests",
				subtitle: "Jobs",
				footer:   m.Help.Model.View(m.MergeRequests.JobsKeys),
				view:     m.MergeRequests.PipelineJobs.View(),
			})

		default:
			return m.renderTableView(renderTableParams{
				title:  "Merge Requests",
				footer: m.Help.Model.View(m.MergeRequests.ListKeys),
				view:   m.MergeRequests.List.View(),
			})

		}
	case tabs.Issues:
		return m.renderTableView(renderTableParams{
			title:  "Issues",
			footer: "help model",
			view:   m.Issues.List.View(),
		})

	case tabs.Pipelines:
		return m.renderTableView(renderTableParams{
			title:  "Pipelines",
			footer: "help model",
			view:   "",
		})
	}

	return "Unsupported View"
}

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
	table := lipgloss.JoinVertical(
		0,
		style.TableTitle.Render(t),
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

func (m *Model) displayToast(msg string, t toast.ToastType) tea.Cmd {
	m.Toast.Show = true
	m.Toast.Message = getErrorMessage(msg)
	m.Toast.Type = t
	return m.Toast.Init()
}

// TODO: Move these to its own module?
type ResponseError struct {
	Message string `json:"message"`
}

func getErrorMessage(msg string) string {
	var e ResponseError

	err := json.Unmarshal([]byte(msg), &e)
	if err != nil {
		return msg
	}
	return e.Message
}
