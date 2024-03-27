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

const (
	RESPONSE_RIGHT_MARGIN      = 2
	useHighPerformanceRenderer = false
)

const (
	tableView currView = iota
	detailView
	commentsView
)

const (
	idColIdx mrRow = iota
	titleColIdx
	authorColIdx
	statusColIdx
	draftColIdx
	conflictsColIdx
	urlColIdx
	descColIdx
)

type (
	currView uint
	mrRow    uint
)

type model struct {
	mergeRequests         mergeRequests
	mergeRequestsComments mergeRequestsComments
	detail                detail
	currView              currView
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.mergeRequests.model.Focused() {
				m.mergeRequests.model.Blur()
			} else {
				m.mergeRequests.model.Focus()
			}
		case "q", "ctrl+c":
			cmds = append(cmds, tea.Quit)
		}

		switch m.currView {
		case tableView:
			switch msg.String() {
			case "x":
				selectedUrl := m.mergeRequests.model.SelectedRow()[urlColIdx]
				command.Openbrowser(selectedUrl)

			case "enter":
				m.setResponseContent()
				m.currView = detailView

			case "c":
				m.mergeRequests.selectedMr = m.mergeRequests.model.SelectedRow()[idColIdx]
				logToFile("selected mr", func() {
					log.Println(m.mergeRequests.selectedMr)
				})
				c := func() tea.Msg {
					r, err := api.GetMRComments(m.mergeRequests.model.SelectedRow()[idColIdx])
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
				m.currView = tableView
			}
		}

	case error:
		logToFile("error", func() {
			log.Println(msg)
		})

	case []table.Row:
		// TODO: extract this logic
		// BUG: Table is not interactive

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)

		columns := []table.Column{
			{Title: "Id", Width: 4},
			{Title: "Type", Width: 10},
			{Title: "Author", Width: 20},
			{Title: "Created At", Width: 20},
			{Title: "Updated At", Width: 20},
			{Title: "Resolved", Width: 10},
			{Title: "Body", Width: 0},
		}
		m.mergeRequestsComments.model.SetStyles(s)
		m.mergeRequestsComments.model.SetHeight(len(msg))

		m.mergeRequestsComments.model.SetColumns(columns)
		m.mergeRequestsComments.model.SetRows(msg)
		m.mergeRequestsComments.model.Focus()
		m.mergeRequestsComments.model.UpdateViewport()
		m.currView = commentsView
		logToFile("rows", func() {
			log.Println(msg)
		})

		isRespReady := func() tea.Msg {
			return "comments_ready"
		}
		cmds = append(cmds, isRespReady)

	case tea.WindowSizeMsg:

		m.mergeRequests.model.SetWidth(msg.Width)
		m.mergeRequestsComments.model.SetWidth(msg.Width)
		// m.table.SetHeight(msg.Height)
		headerHeight := lipgloss.Height(m.headerView(m.mergeRequests.model.SelectedRow()[titleColIdx]))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		cmd = m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	}
	m.mergeRequests.model, cmd = m.mergeRequests.model.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currView == detailView {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.mergeRequests.model.SelectedRow()[titleColIdx]), m.detail.model.View(), m.footerView())
	}

	if m.currView == commentsView {
		return baseStyle.Render(m.mergeRequestsComments.model.View()) + "\n"
	}
	return baseStyle.Render(m.mergeRequests.model.View()) + "\n"
}

func main() {
	config.Load(&config.Config)

	// r := api.GetMRComments("3913")

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := setMergeRequestsModel()
	t.SetStyles(s)

	m := model{mergeRequests: mergeRequests{model: t}, currView: tableView}
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
