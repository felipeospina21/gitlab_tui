package main

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/style"
	"gitlab_tui/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	RESPONSE_RIGHT_MARGIN = 2
)

func main() {
	config.Load(&config.Config)

	// r := api.GetMRComments("3913")

	t := tui.SetMergeRequestsListModel()
	t.SetStyles(style.Table)

	m := tui.Model{
		MergeRequests: tui.MergeRequestsModel{List: t},
		CurrView:      tui.MrTableView,
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
