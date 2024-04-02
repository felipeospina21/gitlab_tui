package main

import (
	"fmt"
	"gitlab_tui/api"
	"gitlab_tui/command"
	"gitlab_tui/config"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Table Styles
var (
	tableStyle       = table.DefaultStyles()
	tableHeaderStyle = tableStyle.Header.BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
	tableSelectedStyle = tableStyle.Selected.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)
)

const (
	RESPONSE_RIGHT_MARGIN      = 2
	useHighPerformanceRenderer = false
)

const (
	mrTableView currView = iota
	detailView
	commentsTableView
	tabsView
)

type (
	currView      uint
	tableColIndex uint
)

type model struct {
	tabs          tabsModel
	mergeRequests mergeRequests
	detail        detailModel
	currView      currView
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global commands
		switch msg.String() {
		case "esc":
			if m.mergeRequests.list.Focused() {
				m.mergeRequests.list.Blur()
			} else {
				m.mergeRequests.list.Focus()
			}
		case "q", "ctrl+c":
			cmds = append(cmds, tea.Quit)

		case "tab":
			m.currView = detailView
		}

		switch m.currView {
		case mrTableView:
			switch msg.String() {
			case "x":
				selectedUrl := m.mergeRequests.list.SelectedRow()[mergeReqsUrlIdx]
				command.Openbrowser(selectedUrl)

			case "enter":
				content := string(m.mergeRequests.list.SelectedRow()[mergeReqsDescIdx])
				m.setResponseContent(content)
				m.currView = detailView

			case "c":
				m.mergeRequests.selectedMr = m.mergeRequests.list.SelectedRow()[mergeReqsIdIdx]
				c := func() tea.Msg {
					r, err := api.GetMRComments(m.mergeRequests.list.SelectedRow()[mergeReqsIdIdx])
					if err != nil {
						return err
					}
					return r
				}
				cmds = append(cmds, c)
			}

		case detailView:
			switch msg.String() {
			case "backspace":
				m.currView = mrTableView
			}

		case commentsTableView:
			switch msg.String() {
			case "x":
				selectedUrl := m.mergeRequests.comments.SelectedRow()[mergeReqsUrlIdx]
				command.Openbrowser(selectedUrl)

			case "enter":
				content := string(m.mergeRequests.comments.SelectedRow()[commentsBodyIdx])
				m.setResponseContent(content)

				// BUG: detail view commands not properly working
				m.currView = detailView
			}
		}

	case error:
		logToFile("error", func() {
			log.Println(msg)
		})

	case []table.Row:
		m.mergeRequests.comments = setMergeRequestsCommentsModel(msg)
		m.mergeRequests.comments.SetStyles(tableStyle)
		m.currView = commentsTableView

		isRespReady := func() tea.Msg {
			return "comments_ready"
		}
		cmds = append(cmds, isRespReady)

	case tea.WindowSizeMsg:

		// NOTE: Resize tabs width
		// numTabs := len(m.tabs.Tabs)
		// x := msg.Width
		// a := x / numTabs
		// inactiveTabStyle.Width(a - docStyle.GetHorizontalPadding())
		// activeTabStyle.Width(a - docStyle.GetHorizontalPadding())

		m.mergeRequests.list.SetWidth(msg.Width)
		m.mergeRequests.comments.SetWidth(msg.Width)
		// m.table.SetHeight(msg.Height)
		headerHeight := lipgloss.Height(m.headerView(m.mergeRequests.list.SelectedRow()[mergeReqsTitleIdx]))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		cmd = m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	}
	m.mergeRequests.list, cmd = m.mergeRequests.list.Update(msg)
	m.mergeRequests.comments, cmd = m.mergeRequests.comments.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.currView {
	case detailView:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.mergeRequests.list.SelectedRow()[mergeReqsTitleIdx]), m.detail.model.View(), m.footerView())

	case commentsTableView:
		return baseStyle.Render(m.mergeRequests.comments.View()) + "\n"

	case tabsView:
		return m.tabsView()

	default:
		return baseStyle.Render(m.mergeRequests.list.View()) + "\n"

	}
}

func main() {
	config.Load(&config.Config)

	// r := api.GetMRComments("3913")

	t := setMergeRequestsListModel()
	t.SetStyles(tableStyle)

	m := model{
		mergeRequests: mergeRequests{list: t},
		currView:      mrTableView,
		// tabs: tabsModel{
		// 	Tabs:       []string{"Merge Requests", "Comments"},
		// 	TabContent: []string{"MRs", "Comments"},
		// },
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// Logs to debug.log file
//
//	logToFile("log", func() {
//		log.Println(strconv.Itoa(msg.Width))
//		log.Println("tw " + strconv.Itoa(m.table.Width()))
//	})
func logToFile(logPrefix string, cb func()) {
	f, err := tea.LogToFile("debug.log", logPrefix)
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	cb()
}
