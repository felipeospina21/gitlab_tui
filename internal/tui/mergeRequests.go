package tui

import (
	"gitlab_tui/internal/style"

	"github.com/charmbracelet/bubbles/table"
)

const (
	mergeReqsIDIdx tableColIndex = iota
	mergeReqsTitleIdx
	mergeReqsAuthorIdx
	mergeReqsStatusIdx
	mergeReqsDraftIdx
	mergeReqsConflictsIdx
	mergeReqsURLIdx
	mergeReqsDescIdx
)

const (
	commentsIDIdx tableColIndex = iota
	commentsTypeIdx
	commentsAuthorIdx
	commentsCreatedAtIdx
	commentsUpdatedAtIdx
	commentsResolvedIdx
	commentsBodyIdx
)

type MergeRequestsModel struct {
	List       table.Model
	Comments   table.Model
	Pipeline   table.Model
	SelectedMr string
	Error      error
}

func (m Model) UpdateMergeRequestsModel(listModel table.Model, commentsModel table.Model, pipelinesModel table.Model) Model {
	listModel.SetStyles(style.Table)
	commentsModel.SetStyles(style.Table)

	newM := Model{
		MergeRequests: MergeRequestsModel{List: listModel, Comments: commentsModel, Pipeline: pipelinesModel},
		CurrView:      m.CurrView,
		Md:            m.Md,
	}

	return newM
}

func InitMergeRequestsListTable(r []table.Row, width int) table.Model {
	id := int(float32(width) * 0.06)
	title := int(float32(width) * 0.5)
	author := int(float32(width) * 0.1)
	status := int(float32(width) * 0.20)
	draft := int(float32(width) * 0.06)
	conf := int(float32(width) * 0.06)
	url := 0

	if width > 170 {
		id = int(float32(width) * 0.03)
		title = int(float32(width) * 0.4)
		status = int(float32(width) * 0.1)
		url = int(float32(width) * 0.24)
	}

	columns := []table.Column{
		{Title: "Iid", Width: id},
		{Title: "Title", Width: title},
		{Title: "Author", Width: author},
		{Title: "Status", Width: status},
		{Title: "Draft", Width: draft},
		{Title: "Conflicts", Width: conf},
		{Title: "Url", Width: url},
		{Title: "Description", Width: 0},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(r),
		table.WithFocused(true),
		table.WithHeight(len(r)),
	)

	return t
}

func SetMergeRequestsCommentsModel(msg []table.Row) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 10},
		{Title: "Type", Width: 20},
		{Title: "Author", Width: 20},
		{Title: "Created At", Width: 30},
		{Title: "Updated At", Width: 30},
		{Title: "Resolved", Width: 10},
		{Title: "Body", Width: 0},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(msg),
		table.WithFocused(true),
		table.WithHeight(len(msg)),
	)

	return t
}

func SetMergeRequestPipelinesModel(msg []table.Row) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "IID", Width: 20},
		{Title: "Status", Width: 20},
		{Title: "Source", Width: 20},
		{Title: "Created At", Width: 30},
		{Title: "Updated At", Width: 30},
		{Title: "URL", Width: 0},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(msg),
		table.WithFocused(true),
		table.WithHeight(len(msg)),
	)

	return t
}
