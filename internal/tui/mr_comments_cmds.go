package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/server"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) refetchComments() (Model, tea.Cmd) {
	r, _ := server.GetMergeRequestComments(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	// TODO: Do something with Error msg
	// TODO: Create enum for queries status
	// if err != nil {
	// 	return m, func() tea.Msg {
	// 		return "fetch_mr_table_error"
	// 	}
	// }
	t := InitMergeRequestsListTable(r, 155)
	newM := m.UpdateMergeRequestsModel(m.MergeRequests.List, t, m.MergeRequests.Pipeline)

	return newM, func() tea.Msg {
		return nil
	}
}

func (m *Model) viewCommentContent() {
	content := string(m.getSelectedMrRow(commentsBodyIdx, MrCommentsView))
	m.setResponseContent(content)
	m.PrevView = MrCommentsView
	m.CurrView = MdView
}

func (m *Model) navigateToMrComment() {
	selectedURL := m.getSelectedMrRow(mergeReqsURLIdx, MrTableView)
	commentID := m.getSelectedMrRow(commentsIDIdx, MrCommentsView)
	exec.Openbrowser(fmt.Sprintf("%s#note_%s", selectedURL, commentID))
}
