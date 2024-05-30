package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	ListTitleStyle        = lipgloss.NewStyle().Margin(2).Foreground(lipgloss.Color(Blue[400]))
	ListItemStyle         = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).PaddingLeft(2).Foreground(lipgloss.Color(Yellow[400]))
	ListSelectedItemStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).PaddingLeft(2).Foreground(lipgloss.Color(Blue[400]))
	ListPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	ListHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	ListQuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
