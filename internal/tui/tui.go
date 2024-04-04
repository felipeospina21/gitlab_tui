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
}

const (
	MrTableView views = iota
	MdView
	CommentsTableView
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
				newM := m.UpdateMergeRequestsModel(t, table.Model{})

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
			}

			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)

		case CommentsTableView:
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
				newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t)

				return newM, func() tea.Msg {
					return "fetch_mr_table_success"
				}

			case "x":
				selectedURL := m.MergeRequests.Comments.SelectedRow()[mergeReqsURLIdx]
				exec.Openbrowser(selectedURL)

			case "enter":
				content := string(m.MergeRequests.Comments.SelectedRow()[commentsBodyIdx])
				m.setResponseContent(content)
				m.PrevView = CommentsTableView
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
			newM := m.UpdateMergeRequestsModel(t, table.Model{})

			return newM, func() tea.Msg {
				return tea.ClearScreen()
			}

		case CommentsTableView:
			m.MergeRequests.Comments.SetWidth(msg.Width)
			t := InitMergeRequestsListTable(m.MergeRequests.Comments.Rows(), msg.Width-10)
			newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t)

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

			case CommentsTableView:
				content = string(m.MergeRequests.Comments.SelectedRow()[commentsBodyIdx])

			default:
				content = string(m.MergeRequests.List.SelectedRow()[mergeReqsDescIdx])
			}
			m.setResponseContent(content)

		}

	case server.MrCommentsQueryResponse:
		m.MergeRequests.Comments = SetMergeRequestsCommentsModel(msg)
		m.MergeRequests.Comments.SetStyles(style.Table)
		m.PrevView = MrTableView
		m.CurrView = CommentsTableView

		isRespReady := func() tea.Msg {
			return "comments_ready"
		}
		cmds = append(cmds, isRespReady)

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

	case CommentsTableView:
		return style.Base.Render(m.MergeRequests.Comments.View()) + "\n"

	case TabsView:
		return m.TabsView()

	default:
		return style.Base.Render(m.MergeRequests.List.View()) + "\n"

	}
}
