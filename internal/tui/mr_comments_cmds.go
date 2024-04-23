package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
)

func (m *Model) refetchComments() {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(mergeReqsIDIdx, MrTableView))
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Comments.SetRows(r)
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
