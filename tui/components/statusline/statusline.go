package statusline

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type modes struct {
	Normal  string
	Insert  string
	Loading string
}

var Modes = modes{
	Normal: "NORMAL",
	Insert: "INSERT",
}

type Model struct {
	Status  string
	Project string
	Content string
	Width   int
	Height  int
}

func New(m Model) Model {
	return Model{
		Status:  m.Status,
		Project: m.Project,
		Content: m.Content,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) View() string {
	width := m.Width
	w := lipgloss.Width

	statusKey := statusStyle.Render(m.Status)
	encoding := encodingStyle.Render("UTF-8")
	// projectName := fishCakeStyle.Render("🍥 Fish Cake")
	projectName := projectStyle.Render(m.Project)
	statusVal := statusText.
		Width(width - w(statusKey) - w(encoding) - w(projectName)).
		Render(m.Content)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		projectName,
	)

	return StatusBarStyle.Render(bar)
}

// TODO: add a function to control statuses
