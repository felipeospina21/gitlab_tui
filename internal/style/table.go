package style

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	Table = func() table.Styles {
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(Violet[400])).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color(Violet[50])).
			Background(lipgloss.Color(Violet[800])).
			Bold(false)

		return s
	}
	TableTitle = lipgloss.NewStyle().Margin(2, 0, 1, 2).Foreground(lipgloss.Color(Violet[300])).Bold(true)
)
