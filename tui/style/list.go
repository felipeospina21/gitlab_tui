package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	ListTitleStyle        = lipgloss.NewStyle().Margin(2).Foreground(lipgloss.Color(Blue[400]))
	ListItemStyle         = lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color(Violet[300]))
	ListSelectedItemStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).PaddingLeft(2).Foreground(lipgloss.Color(Violet[50])).Background(lipgloss.Color(Violet[800]))
	ListPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	ListHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	ListDocStyle          = lipgloss.NewStyle().
				MarginTop(1).
				PaddingRight(4).
				Foreground(lipgloss.Color(Violet[300])).
				BorderRight(true).
				BorderStyle(lipgloss.NormalBorder()).
				Width(50) // TODO: Set width in config file
)

var (
	hl = string(lipgloss.Color(Violet[400]))
	fg = string(lipgloss.Color(Violet[50]))
)

type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		// Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: fg}).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: hl}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: hl})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}
