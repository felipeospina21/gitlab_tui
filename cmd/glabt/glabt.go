package main

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	responseRightMargin = 2
)

func main() {
	config.Load(&config.Config)

	m := InitModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func InitModel() tui.Model {
	l := tui.InitProjectsList()

	newM := tui.Model{
		Projects: tui.ProjectsModel{List: l},
		CurrView: tui.ProjectsView,
	}

	return newM
}
