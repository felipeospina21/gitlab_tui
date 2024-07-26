package main

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/tui"
	"gitlab_tui/tui/components"
	"gitlab_tui/tui/components/progress"
	"gitlab_tui/tui/components/tabs"
	"gitlab_tui/tui/components/toast"
	"os"

	"github.com/charmbracelet/bubbles/help"
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
	p := tui.InitPaginatorModel()

	newM := tui.Model{
		Projects: tui.ProjectsModel{List: l},
		CurrView: tui.ProjectsView,
		Help:     components.Help{Model: help.New()},
		MergeRequests: tui.MergeRequestsModel{
			ListKeys:     tui.MergeReqsKeys,
			CommentsKeys: tui.CommentsKeys,
			PipelineKeys: tui.PipelinKeys,
			JobsKeys:     tui.JobsKeys,
		},
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
		Paginator: p,
	}
	return newM
}
