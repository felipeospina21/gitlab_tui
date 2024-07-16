package toast

import (
	"gitlab_tui/internal/style"

	"github.com/charmbracelet/lipgloss"
)

func toastStyle(h, w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#111111")).
		Margin(1, 1, 0, 1).
		Padding(1).
		Width(w).
		Height(h)
}

func successStyle(h, w int) lipgloss.Style {
	return toastStyle(h, w).
		Background(lipgloss.Color(style.Green[400]))
}

func errorStyle(h, w int) lipgloss.Style {
	return toastStyle(h, w).
		Background(lipgloss.Color(style.Red[400]))
}

func warningStyle(h, w int) lipgloss.Style {
	return toastStyle(h, w).
		Background(lipgloss.Color(style.Yellow[400]))
}

func infoStyle(h, w int) lipgloss.Style {
	return toastStyle(h, w).
		Background(lipgloss.Color(style.Violet[400]))
}
