package tui

import (
	"encoding/json"
	"fmt"
	"gitlab_tui/tui/components"
	"gitlab_tui/tui/components/progress"
	"gitlab_tui/tui/components/table"
	"gitlab_tui/tui/components/tabs"
	"gitlab_tui/tui/components/toast"
	"gitlab_tui/tui/style"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	views uint
)

type Model struct {
	Projects        ProjectsModel
	MergeRequests   MergeRequestsModel
	Issues          IssuesModel
	Md              MdModel
	CurrView        views
	PrevView        views
	Title           string
	Window          tea.WindowSizeMsg
	Help            components.Help
	Toast           toast.Model
	Tabs            tabs.Model
	Paginator       paginator.Model
	isSidePanelOpen bool
}

const (
	MainTableView views = iota
	MrCommentsView
	MrPipelinesView
	JobsView
	MdView
	HomeView
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

func InitMainModel() Model {
	l := InitProjectsList()
	p := InitPaginatorModel()

	return Model{
		Projects:      ProjectsModel{List: l},
		CurrView:      HomeView,
		Help:          components.Help{Model: help.New()},
		MergeRequests: MergeRequestsModel{},
		Toast: toast.New(toast.Model{
			Progress: progress.New(
				progress.WithDefaultGradient(),
				progress.WithFillCharacters('-', ' '),
				progress.WithoutPercentage(),
			),
			Interval: 10,
			// Type:     toast.Info,
			// Show:     true,
			// Message:  "Info msg",
		}),
		Tabs: tabs.Model{
			Tabs: []string{"Merge Requests", "Issues", "Pipelines"},
		},
		Paginator:       p,
		isSidePanelOpen: true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Toast.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		cmd, cmds = m.updateKeyMsg(msg)

	case tea.WindowSizeMsg:
		m, cmd = m.updateWindowSize(msg)

	case string:
		m, cmd = m.updateStringMsg(msg)

	case error:
		m, cmd = m.updateErrorMsg(msg)

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
	projects := style.ListDocStyle.Render(m.Projects.List.View())
	if m.isSidePanelOpen {
		return lipgloss.JoinHorizontal(0, projects, m.renderTableView(renderTableParams{
			title:  "Merge Requests",
			footer: m.Help.Model.View(MergeReqsKeys),
			view:   m.MergeRequests.List.View(),
		}),
		)
	} else {

		switch m.Tabs.ActiveTab {
		case tabs.MergeRequests:
			switch m.CurrView {
			case MainTableView:
				return m.renderTableView(renderTableParams{
					title:  "Merge Requests",
					footer: m.Help.Model.View(MergeReqsKeys),
					view:   m.MergeRequests.List.View(),
				})
			// FIX: when navigating back from Issues List it never gets here (active tab is Issues)
			case HomeView:
				projects := style.ListDocStyle.Render(m.Projects.List.View())

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
					footer:   m.Help.Model.View(CommentsKeys),
					view:     m.MergeRequests.Comments.View(),
				})

			case MrPipelinesView:
				return m.renderTableView(renderTableParams{
					title:    "Merge Requests",
					subtitle: "Pipelines",
					footer:   m.Help.Model.View(PipelineKeys),
					view:     m.MergeRequests.Pipeline.View(),
				})

			case JobsView:
				return m.renderTableView(renderTableParams{
					title:    "Merge Requests",
					subtitle: "Jobs",
					footer:   m.Help.Model.View(JobsKeys),
					view:     m.MergeRequests.PipelineJobs.View(),
				})

			default:
				return "default"

			}
		case tabs.Issues:
			switch m.CurrView {
			case MainTableView:
				return m.renderTableView(renderTableParams{
					title:  "Issues",
					footer: m.Help.Model.View(IssuesKeys),
					view:   m.Issues.List.View(),
				})

			case MdView:
				return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.Issues.List.SelectedRow()[table.IssuesListCols.Title.Idx]), m.Md.Viewport.View(), m.footerView())
			}

		case tabs.Pipelines:
			return m.renderTableView(renderTableParams{
				title:  "Pipelines",
				footer: "help model",
				view:   "",
			})
		}

		return "Unsupported View"
	}
}

func (m Model) getSelectedRow(idx table.TableColIndex, view views) string {
	switch view {
	case MainTableView:
		if m.Tabs.ActiveTab == tabs.MergeRequests {
			return m.MergeRequests.List.SelectedRow()[idx]
		}
		if m.Tabs.ActiveTab == tabs.Issues {
			return m.Issues.List.SelectedRow()[idx]
		}
		return "tab not supported (getSelectedRow)"

	case MrCommentsView:
		return m.MergeRequests.Comments.SelectedRow()[idx]

	case MrPipelinesView:
		return m.MergeRequests.Pipeline.SelectedRow()[idx]

	case JobsView:
		return m.MergeRequests.PipelineJobs.SelectedRow()[idx]

	default:
		return "View Not supported"

	}
}

func (m *Model) setSelectedMr() {
	m.MergeRequests.SelectedMr = m.getSelectedRow(table.MergeReqsCols.Title.Idx, MainTableView)
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
