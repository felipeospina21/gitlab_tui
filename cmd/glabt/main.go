package main

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/server"
	"gitlab_tui/internal/style"
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
	m, _ := InitModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func InitModel() (tui.Model, error) {
	r, err := server.GetMergeRequestsMock()
	// r := server.GetMergeRequests()

	t := tui.InitMergeRequestsListTable(r, 155)
	t.SetStyles(style.Table)

	newM := tui.Model{
		MergeRequests: tui.MergeRequestsModel{List: t},
	}

	return newM, err
}
