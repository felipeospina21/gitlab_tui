package tui

import (
	"gitlab_tui/tui/components/toast"
	"gitlab_tui/tui/style"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateErrorMsg(msg error) (Model, tea.Cmd) {
	cmd := m.displayToast(msg.Error(), toast.Error)

	lh, lv := style.ListItemStyle.GetFrameSize()
	nh, nv := toast.ErrorStyle(m.Window.Height, m.Window.Width).GetFrameSize()

	h := (lh + nh) * 2
	v := (lv + nv) * 2

	m.Projects.List.SetSize(m.Window.Width-h, m.Window.Height-v)

	return *m, cmd
}
