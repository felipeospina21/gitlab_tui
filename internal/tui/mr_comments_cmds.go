package tui

import (
	"fmt"
	"gitlab_tui/internal/exec"
	"gitlab_tui/internal/logger"
	"gitlab_tui/internal/server"
)

func (m *Model) refetchComments() {
	r, err := server.GetMergeRequestComments(m.getSelectedMrRow(mergeReqsCols.iid.idx, MrTableView), m.Projects.ProjectID)
	if err != nil {
		logger.Error(err)
	}
	m.MergeRequests.Comments.SetRows(r)
}

func (m *Model) viewCommentContent() {
	content := string(m.getSelectedMrRow(commentsCols.body.idx, MrCommentsView))
	m.setResponseContent(content)
	m.PrevView = MrCommentsView
	m.CurrView = MdView
}

func (m *Model) navigateToMrComment() {
	selectedURL := m.getSelectedMrRow(mergeReqsCols.url.idx, MrTableView)
	commentID := m.getSelectedMrRow(commentsCols.id.idx, MrCommentsView)
	exec.Openbrowser(fmt.Sprintf("%s#note_%s", selectedURL, commentID))
}
