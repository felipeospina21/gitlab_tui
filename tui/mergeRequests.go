package tui

import (
	"gitlab_tui/api"

	"github.com/charmbracelet/bubbles/table"
)

const (
	mergeReqsIdIdx tableColIndex = iota
	mergeReqsTitleIdx
	mergeReqsAuthorIdx
	mergeReqsStatusIdx
	mergeReqsDraftIdx
	mergeReqsConflictsIdx
	mergeReqsUrlIdx
	mergeReqsDescIdx
)

const (
	commentsIdIdx tableColIndex = iota
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
	SelectedMr string
}

func SetMergeRequestsListModel() table.Model {
	r := api.GetMergeRequests()

	columns := []table.Column{
		{Title: "Iid", Width: 4},
		{Title: "Title", Width: 80},
		{Title: "Author", Width: 20},
		{Title: "Status", Width: 30},
		{Title: "Draft", Width: 10},
		{Title: "Conflicts", Width: 10},
		{Title: "Url", Width: 0},
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
