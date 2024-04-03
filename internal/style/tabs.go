package style

import "github.com/charmbracelet/lipgloss"

func TabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	InactiveTabBorder = TabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = TabBorderWithBottom("┘", " ", "└")
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	Tabs              = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	InactiveTab       = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTab         = InactiveTab.Copy().Border(ActiveTabBorder, true)
	TabsWindow        = lipgloss.NewStyle().BorderForeground(HighlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)
