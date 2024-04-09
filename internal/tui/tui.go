package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
	"gitlab_tui/internal/style"
	"log"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	views         uint
	tableColIndex uint
)

type Model struct {
	Tabs          TabsModel
	MergeRequests MergeRequestsModel
	Md            MdModel
	CurrView      views
	PrevView      views
	Title         string
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
				r, err := server.GetMergeRequestsMock()
				// TODO: Do something with Error msg
				// TODO: Create enum for queries status
				if err != nil {
					return m, func() tea.Msg {
						return "fetch_mr_table_error"
					}
				}

				t := InitMergeRequestsListTable(r, 155)
				newM := m.UpdateMergeRequestsModel(t, table.Model{}, table.Model{})

				return newM, func() tea.Msg {
					return "fetch_mr_table_success"
				}

			case "x":
				selectedURL := m.MergeRequests.List.SelectedRow()[mergeReqsURLIdx]
				exec.Openbrowser(selectedURL)

			case "enter":
				content := string(m.MergeRequests.List.SelectedRow()[mergeReqsDescIdx])
				m.setResponseContent(content)
				m.PrevView = MrTableView
				m.CurrView = MdView

			case "c":
				c := func() tea.Msg {
					r, err := server.GetMergeRequestCommentsMock(m.MergeRequests.List.SelectedRow()[mergeReqsIDIdx])
					// r, err := server.GetMergeRequestComments(m.MergeRequests.List.SelectedRow()[mergeReqsIDIdx])
					if err != nil {
						return err
					}
					return r
				}
				cmds = append(cmds, c)

			case "p":
				c := func() tea.Msg {
					r, err := server.GetMergeRequestPipelinesMock(m.MergeRequests.List.SelectedRow()[mergeReqsIDIdx])
					if err != nil {
						return err
					}
					return r
				}

				cmds = append(cmds, c)
			}

			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)

		case MrCommentsView:
			switch msg.String() {
			case "r":
				r, err := server.GetMergeRequestCommentsMock(m.MergeRequests.List.SelectedRow()[mergeReqsIDIdx])
				// TODO: Do something with Error msg
				// TODO: Create enum for queries status
				if err != nil {
					return m, func() tea.Msg {
						return "fetch_mr_table_error"
					}
				}
				t := InitMergeRequestsListTable(r, 155)
				newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t, m.MergeRequests.Pipeline)

				return newM, func() tea.Msg {
					return "fetch_mr_table_success"
				}

			case "x":
				selectedURL := m.MergeRequests.Comments.SelectedRow()[mergeReqsURLIdx]
				exec.Openbrowser(selectedURL)

			case "enter":
				content := string(m.MergeRequests.Comments.SelectedRow()[commentsBodyIdx])
				m.setResponseContent(content)
				m.PrevView = MrCommentsView
				m.CurrView = MdView

			case "backspace":
				m.CurrView = MrTableView

			}

			m.MergeRequests.Comments, cmd = m.MergeRequests.Comments.Update(msg)

		case MdView:
			switch msg.String() {
			case "backspace":
				m.CurrView = m.PrevView
			}

			m.Md.Viewport, cmd = m.Md.Viewport.Update(msg)
		}

	case tea.WindowSizeMsg:
		cmd = m.setViewportViewSize(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		switch m.CurrView {
		case MrTableView:
			m.MergeRequests.List.SetWidth(msg.Width)
			t := InitMergeRequestsListTable(m.MergeRequests.List.Rows(), msg.Width-10)
			newM := m.UpdateMergeRequestsModel(t, table.Model{}, table.Model{})

			return newM, func() tea.Msg {
				return tea.ClearScreen()
			}

		case MrCommentsView:
			m.MergeRequests.Comments.SetWidth(msg.Width)
			t := InitMergeRequestsListTable(m.MergeRequests.Comments.Rows(), msg.Width-10)
			newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t, m.MergeRequests.Pipeline)

			return newM, func() tea.Msg {
				return tea.ClearScreen()
			}

		case MrPipelinesView:
			m.MergeRequests.Pipeline.SetWidth(msg.Width)
			t := InitMergeRequestsListTable(m.MergeRequests.Pipeline.Rows(), msg.Width-10)
			newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, m.MergeRequests.Comments, t)

			return newM, func() tea.Msg {
				return tea.ClearScreen()
			}

		case MdView:
			headerHeight := lipgloss.Height(m.headerView(m.MergeRequests.List.SelectedRow()[mergeReqsTitleIdx]))
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			m.Md.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)

			var content string
			switch m.PrevView {
			case MrTableView:
				content = string(m.MergeRequests.List.SelectedRow()[mergeReqsDescIdx])

			case MrCommentsView:
				content = string(m.MergeRequests.Comments.SelectedRow()[commentsBodyIdx])

			default:
				content = string(m.MergeRequests.List.SelectedRow()[mergeReqsDescIdx])
			}
			m.setResponseContent(content)

		}

	case server.MrCommentsQueryResponse:
		m.MergeRequests.Comments = SetMergeRequestsCommentsModel(msg)
		m.MergeRequests.Comments.SetStyles(style.Table)
		m.MergeRequests.SelectedMr = m.MergeRequests.List.SelectedRow()[mergeReqsTitleIdx]
		m.PrevView = MrTableView
		m.CurrView = MrCommentsView

		isRespReady := func() tea.Msg {
			return "comments"
		}
		cmds = append(cmds, isRespReady)

	case *server.MrPipelinesQueryResponse:
		m.MergeRequests.Pipeline = SetMergeRequestPipelinesModel(*msg)
		// m.Title = "Pipelines"
		// // m.MergeRequests.Pipeline.SetStyles(style.Table)
		// // m.MergeRequests.SelectedMr = m.MergeRequests.List.SelectedRow()[mergeReqsTitleIdx]
		// // m.PrevView = MrTableView
		// // m.CurrView = MrPipelinesView
		// isRespReady := func() tea.Msg {
		// 	return "pipeline"
		// }
		// cmds = append(cmds, isRespReady)

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

	case TabsView:
		return m.TabsView()

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
		t = fmt.Sprintf("Merge Request %s / %s", m.Title, m.MergeRequests.SelectedMr)
		// if m.CurrView == MrCommentsView {
		// 	t = fmt.Sprintf("Merge Request Comments / %s", m.MergeRequests.SelectedMr)
		// } else {
		// 	t = fmt.Sprintf("Merge Request Pipelines / %s", m.MergeRequests.SelectedMr)
		// }
	}
	return lipgloss.JoinVertical(0, titleStyle.Render(t), style.Base.Render(view)+"\n")
}
