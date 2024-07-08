package style

import "github.com/charmbracelet/lipgloss"

var ErrorNotification = func(h, w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(Red[400])).
		Foreground(lipgloss.Color("#111111")).
		Align(lipgloss.Center).
		Padding(1).
		Margin(1, 1, 0, 1).
		Width(w).
		Bold(true)
}
