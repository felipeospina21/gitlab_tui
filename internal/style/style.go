package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Base = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
