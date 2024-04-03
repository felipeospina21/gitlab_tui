package tui

import (
	"gitlab_tui/internal/style"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TabsModel struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m *Model) TabsView() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs.Tabs {
		var stl lipgloss.Style
		isFirst := i == 0
		isLast := i == len(m.Tabs.Tabs)-1
		isActive := i == m.Tabs.activeTab

		if isActive {
			stl = style.ActiveTab.Copy()
		} else {
			stl = style.InactiveTab.Copy()
		}
		border, _, _, _, _ := stl.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		stl = stl.Border(border)
		renderedTabs = append(renderedTabs, stl.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(style.TabsWindow.Width((lipgloss.Width(row) - style.TabsWindow.GetHorizontalFrameSize())).Render(m.Tabs.TabContent[m.Tabs.activeTab]))
	return style.Tabs.Render(doc.String())
}
