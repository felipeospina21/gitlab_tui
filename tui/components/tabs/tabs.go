package tabs

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	MergeRequests int = iota
	Issues
	Pipelines
)

type Model struct {
	Tabs      []string
	Content   string
	ActiveTab int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "shift+tab":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

// TODO: move styles to its own file
var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	// TODO: Update colors with tokens
	highlightColor  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	windowStyle = func(w int) lipgloss.Style {
		return lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder(), false).UnsetBorderTop().Width(w)
	}
)

func (m Model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isActive := i == m.ActiveTab

		if isActive {
			style = activeTab
		} else {
			style = tab
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	x := physicalWidth - lipgloss.Width(row)
	gap := tabGap.Render(strings.Repeat(" ", max(0, x-10)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(m.Content)
	// doc.WriteString(windowStyle(lipgloss.Width(row)).Render(m.Content))
	//
	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
