package tui

import (
	"errors"
	"gitlab_tui/tui/components/table"
	"gitlab_tui/tui/components/tabs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateKeyMsg(msg tea.KeyMsg) (tea.Cmd, []tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	// Global commands
	switch {
	case key.Matches(msg, GlobalKeys.Help):
		m.Help.Model.ShowAll = !m.Help.Model.ShowAll

	case key.Matches(msg, GlobalKeys.Quit):
		cmds = append(cmds, tea.Quit)

	case key.Matches(msg, GlobalKeys.ThrowError):
		cmds = append(cmds, func() tea.Msg {
			return errors.New("mocked")
		})
	}

	// Tabs
	switch m.Tabs.ActiveTab {
	case tabs.MergeRequests:
		switch {
		case key.Matches(msg, GlobalKeys.NextTab):
			if !m.Issues.HasData {
				cmds = append(cmds, m.viewIssues())
			}
		}

	case tabs.Issues:
		switch {
		case key.Matches(msg, GlobalKeys.NextPage):
			if m.Issues.NexPage != "" {
				cmds = append(cmds, m.getIssuesNextPage())
			}

		// case key.Matches(msg, GlobalKeys.NavigateBack):
		// 	m.Tabs.ActiveTab = m.Tabs.ActiveTab - 1

		case key.Matches(msg, GlobalKeys.PrevPage):
			if m.Issues.PrevPage != "" {
				cmds = append(cmds, m.getIssuesPrevPage())
			}
		}

	}

	// Views commands
	switch m.CurrView {
	case HomeView:
		switch {
		case key.Matches(msg, GlobalKeys.ToggleSidePanel):
			m.toggleSidePanel()

		case key.Matches(msg, ProjectsKeys.ViewMRs):
			m.toggleSidePanel()
			m.Issues.HasData = false
			m.Tabs.ActiveTab = tabs.MergeRequests
			c := m.viewMergeReqs(m.Window)
			cmds = append(cmds, c)
		}
		m.Projects.List, cmd = m.Projects.List.Update(msg)

	case MdView:
		switch {
		case key.Matches(msg, GlobalKeys.NavigateBack):
			m.CurrView = m.PrevView
		}
		m.Md.Viewport, cmd = m.Md.Viewport.Update(msg)

	case MainTableView:
		if key.Matches(msg, GlobalKeys.ToggleSidePanel) {
			m.toggleSidePanel()
		}

		switch m.Tabs.ActiveTab {
		case tabs.MergeRequests:
			switch {
			case key.Matches(msg, MergeReqsKeys.OpenInBrowser):
				m.openInBrowser(table.MergeReqsCols.URL.Idx, MainTableView)

			case key.Matches(msg, MergeReqsKeys.Comments):
				c := m.viewComments()
				cmds = append(cmds, c)

			case key.Matches(msg, MergeReqsKeys.Pipelines):
				c := m.viewPipelines()
				cmds = append(cmds, c)

			case key.Matches(msg, MergeReqsKeys.Merge):
				c := m.mergeMR()
				cmds = append(cmds, c)

			case key.Matches(msg, MergeReqsKeys.Description):
				m.viewDescription()

			case key.Matches(msg, MergeReqsKeys.Refetch):
				m.refetchMrList()

			case key.Matches(msg, GlobalKeys.NavigateBack):
				m.CurrView = HomeView

			}
			m.MergeRequests.List, cmd = m.MergeRequests.List.Update(msg)

		case tabs.Issues:
			switch {
			case key.Matches(msg, IssuesKeys.Description):
				m.viewIssueDescription()

			case key.Matches(msg, IssuesKeys.Refetch):
			// TODO: implement refetch

			case key.Matches(msg, IssuesKeys.OpenInBrowser):
				m.openInBrowser(table.IssuesListCols.URL.Idx, MainTableView)

			case key.Matches(msg, GlobalKeys.NavigateBack):
				m.CurrView = HomeView
			}
		}

	case MrCommentsView:
		switch {
		case key.Matches(msg, CommentsKeys.Refetch):
			m.refetchComments()

		case key.Matches(msg, CommentsKeys.OpenInBrowser):
			m.navigateToMrComment()

		case key.Matches(msg, CommentsKeys.Description):
			m.viewCommentContent()

		case key.Matches(msg, GlobalKeys.NavigateBack):
			m.CurrView = MainTableView

		}
		m.MergeRequests.Comments, cmd = m.MergeRequests.Comments.Update(msg)

	case MrPipelinesView:
		switch {
		case key.Matches(msg, PipelineKeys.Jobs):
			c := m.viewPipelineJobs()
			cmds = append(cmds, c)

		case key.Matches(msg, PipelineKeys.Refetch):
			m.refetchPipelines()

		case key.Matches(msg, PipelineKeys.OpenInBrowser):
			m.openInBrowser(table.PipelinesCols.URL.Idx, MrPipelinesView)

		case key.Matches(msg, GlobalKeys.NavigateBack):
			m.CurrView = MainTableView

		}
		m.MergeRequests.Pipeline, cmd = m.MergeRequests.Pipeline.Update(msg)

	case JobsView:
		switch {
		case key.Matches(msg, GlobalKeys.NavigateBack):
			m.CurrView = MrPipelinesView

		case key.Matches(msg, JobsKeys.OpenInBrowser):
			m.openInBrowser(table.PipelineJobsCols.URL.Idx, JobsView)

		case key.Matches(msg, JobsKeys.Refetch):
			m.refetchJobs()

		}
		m.MergeRequests.PipelineJobs, cmd = m.MergeRequests.PipelineJobs.Update(msg)

	}

	return cmd, cmds
}
