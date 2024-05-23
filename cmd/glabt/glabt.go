package main

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	responseRightMargin = 2
)

func main() {
	config.Load(&config.Config)

	// TODO: handle fetching error
	m := InitModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func InitModel() tui.Model {
	// r, err := server.GetMergeRequests()

	// t := tui.InitMergeRequestsListTable(r, 155)
	l := tui.InitProjectsList()

	newM := tui.Model{
		Projects: tui.ProjectsModel{List: l},
		CurrView: tui.ProjectsView,
		// MergeRequests: tui.MergeRequestsModel{},
	}

	return newM
}
