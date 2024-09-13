package tui

import (
	"fmt"
	"gitlab_tui/internal/logger"
	"gitlab_tui/tui/style"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type (
	responseMsg        string
	contentRenderedMsg string
	isResponseReady    bool
	errMsg             struct{ err error }
)

type MdModel struct {
	Viewport viewport.Model
	Ready    isResponseReady
	Content  responseMsg
	Err      error
}

func (e errMsg) Error() string { return e.err.Error() }

func (m *Model) headerView(queryName string) string {
	title := style.MdTitle.Render(queryName)
	line := strings.Repeat("─", max(0, m.Md.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	info := style.MdInfo.Render(fmt.Sprintf("%3.f%%", m.Md.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Md.Viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *Model) setResponseContent(content string) {
	styledContent := renderWithGlamour(m.Md, content)

	m.Md.Viewport.SetContent(styledContent)
}

func (m *Model) setViewportViewSize(msg tea.WindowSizeMsg) tea.Cmd {
	w := msg.Width
	headerHeight := lipgloss.Height(m.headerView(""))
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	if !m.Md.Ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.Md.Viewport = viewport.New(w, msg.Height-verticalMarginHeight)
		m.Md.Viewport.YPosition = headerHeight
		m.Md.Viewport.HighPerformanceRendering = useHighPerformanceRenderer

		// m.setResponseContent()
		m.Md.Ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.Md.Viewport.YPosition = headerHeight + 1
	} else {
		m.Md.Viewport.Width = w
		m.Md.Viewport.Height = msg.Height - verticalMarginHeight
	}
	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		// cmds = append(cmds, viewport.Sync(m.viewport.mod))
		return viewport.Sync(m.Md.Viewport)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderWithGlamour(m MdModel, md string) string {
	s, err := glamourRender(m, md)
	if err != nil {
		logger.Error(err)
	}
	return s
}

// This is where the magic happens.
func glamourRender(m MdModel, markdown string) (string, error) {
	// initialize glamour
	gs := glamour.WithStandardStyle(glamour.DarkStyle)

	width := m.Viewport.Width
	r, err := glamour.NewTermRenderer(
		gs,
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	)
	if err != nil {
		return "", err
	}

	out, err := r.Render(markdown)
	if err != nil {
		return "", err
	}

	// trim lines
	lines := strings.Split(out, "\n")

	var content string
	for i, s := range lines {
		content += strings.TrimSpace(s)

		// don't add an artificial newline after the last split
		if i+1 < len(lines) {
			content += "\n"
		}
	}

	return content, nil
}
