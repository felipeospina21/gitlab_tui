package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type (
	responseMsg        string
	contentRenderedMsg string
	isResponseReady    bool
	errMsg             struct{ err error }
)

type detailModel struct {
	model   viewport.Model
	ready   isResponseReady
	content responseMsg
	err     error
}

func (e errMsg) Error() string { return e.err.Error() }

func (m *model) headerView(queryName string) string {
	title := titleStyle.Render(queryName)
	line := strings.Repeat("─", max(0, m.detail.model.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.detail.model.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.detail.model.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *model) setResponseContent(content string) {
	// styledContent := lipgloss.NewStyle().Width(m.detail.model.Width - RESPONSE_RIGHT_MARGIN).Render(content)

	styledContent := renderWithGlamour(m.detail, content)

	m.detail.model.SetContent(styledContent)
}

func (m *model) setViewportViewSize(msg tea.WindowSizeMsg, headerHeight int, verticalMarginHeight int) tea.Cmd {
	w := msg.Width

	if !m.detail.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.detail.model = viewport.New(w, msg.Height-verticalMarginHeight)
		m.detail.model.YPosition = headerHeight
		m.detail.model.HighPerformanceRendering = useHighPerformanceRenderer

		// m.setResponseContent()
		m.detail.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.detail.model.YPosition = headerHeight + 1
	} else {
		m.detail.model.Width = w
		m.detail.model.Height = msg.Height - verticalMarginHeight
	}
	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		// cmds = append(cmds, viewport.Sync(m.viewport.mod))
		return viewport.Sync(m.detail.model)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderWithGlamour(m detailModel, md string) string {
	s, err := glamourRender(m, md)
	if err != nil {
		log.Println("error rendering with Glamour:", err)
	}
	return s
}

// This is where the magic happens.
func glamourRender(m detailModel, markdown string) (string, error) {
	// initialize glamour
	var gs glamour.TermRendererOption
	gs = glamour.WithAutoStyle()

	// width := max(0, min(int(m.common.cfg.GlamourMaxWidth), m.model.Width))
	width := m.model.Width
	r, err := glamour.NewTermRenderer(
		gs,
		glamour.WithWordWrap(width),
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
