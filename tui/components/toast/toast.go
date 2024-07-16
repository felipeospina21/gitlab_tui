package toast

import (
	"gitlab_tui/tui/components/progress"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 0
	maxWidth = 80
	fps      = 60
)

const (
	Success ToastType = iota
	Error
	Warning
	Info
)

type (
	ToastType uint
	tickMsg   time.Time
)

type Model struct {
	Progress progress.Model
	Interval time.Duration
	Type     ToastType
	Width    int
	Height   int
	Message  string
	Show     bool
	Timer    *time.Timer
	percent  float64
}

func New(m Model) Model {
	return Model{
		Progress: m.Progress,
		Interval: m.Interval,
		Type:     m.Type,
		Message:  m.Message,
		Show:     m.Show,
		percent:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.Show {
			m.percent += 0.01
			if m.percent > 1.0 {
				m.percent = 0
				m.Message = ""
				m.Show = false
				return m, nil
			}
			return m, m.tickCmd()

		}
		return m, nil

	default:
		return m, nil
	}
}

func (m Model) View() string {
	pad := strings.Repeat(" ", padding)
	bar := pad + m.Progress.ViewAs(m.percent) + "\n\n"

	var toast string
	switch m.Type {
	case Success:
		toast = SuccessStyle(m.Height, m.Width).Render(m.Message)

	case Error:
		toast = ErrorStyle(m.Height, m.Width).Render(m.Message)

	case Warning:
		toast = WarningStyle(m.Height, m.Width).Render(m.Message)

	case Info:
		toast = InfoStyle(m.Height, m.Width).Render(m.Message)
	}

	return lipgloss.JoinVertical(0, toast, bar)
}

func (m Model) tickCmd() tea.Cmd {
	return m.Tick(time.Millisecond*m.Interval*10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) Tick(d time.Duration, fn func(time.Time) tea.Msg) tea.Cmd {
	t := time.NewTimer(d)
	m.Timer = t
	return func() tea.Msg {
		ts := <-t.C
		t.Stop()
		for len(t.C) > 0 {
			<-t.C
		}
		return fn(ts)
	}
}
