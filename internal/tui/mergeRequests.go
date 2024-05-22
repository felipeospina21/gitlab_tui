package tui

import (
	"gitlab_tui/internal/style"

	"github.com/charmbracelet/bubbles/table"
)

type tableCol struct {
	name string
	idx  tableColIndex
}

type mergeReqsTable struct {
	iid         tableCol
	title       tableCol
	author      tableCol
	status      tableCol
	mergeStatus tableCol
	draft       tableCol
	confilcts   tableCol
	url         tableCol
	desc        tableCol
}

type mergeReqsCommentsTable struct {
	id          tableCol
	commentType tableCol
	author      tableCol
	createdAt   tableCol
	updatedAt   tableCol
	resolved    tableCol
	body        tableCol
}

type mergeReqsPipelinesTable struct {
	id        tableCol
	iid       tableCol
	status    tableCol
	source    tableCol
	createdAt tableCol
	updatedAt tableCol
	url       tableCol
}

var mergeReqsCols = mergeReqsTable{
	iid:         tableCol{idx: 0, name: "Iid"},
	title:       tableCol{idx: 1, name: "Title"},
	author:      tableCol{idx: 2, name: "Author"},
	status:      tableCol{idx: 3, name: "Status"},
	mergeStatus: tableCol{idx: 4, name: "Merge Status"},
	draft:       tableCol{idx: 5, name: "Draft"},
	confilcts:   tableCol{idx: 6, name: "Conflicts"},
	url:         tableCol{idx: 7, name: "Url"},
	desc:        tableCol{idx: 8, name: "Description"},
}

var commentsCols = mergeReqsCommentsTable{
	id:          tableCol{idx: 0, name: "Id"},
	commentType: tableCol{idx: 1, name: "Type"},
	author:      tableCol{idx: 2, name: "Author"},
	createdAt:   tableCol{idx: 3, name: "Created At"},
	updatedAt:   tableCol{idx: 4, name: "Updated At"},
	resolved:    tableCol{idx: 5, name: "Resolved"},
	body:        tableCol{idx: 6, name: "Body"},
}

var pipelinesCols = mergeReqsPipelinesTable{
	id:        tableCol{idx: 0, name: "Id"},
	iid:       tableCol{idx: 1, name: "IID"},
	status:    tableCol{idx: 2, name: "Status"},
	source:    tableCol{idx: 3, name: "Source"},
	createdAt: tableCol{idx: 4, name: "Created At"},
	updatedAt: tableCol{idx: 5, name: "Updated At"},
	url:       tableCol{idx: 6, name: "URL"},
}

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
	title := int(float32(width) * 0.45)
	author := int(float32(width) * 0.1)
	status := int(float32(width) * 0.13)
	icon := int(float32(width) * 0.04)
	url := 0

	if width > 170 {
		id = int(float32(width) * 0.03)
		title = int(float32(width) * 0.35)
		status = int(float32(width) * 0.1)
		total := id + title + author + (status * 2) + (icon * 2)
		url = width - total - 10
	}

	columns := []table.Column{
		{Title: mergeReqsCols.iid.name, Width: id},
		{Title: mergeReqsCols.title.name, Width: title},
		{Title: mergeReqsCols.author.name, Width: author},
		{Title: mergeReqsCols.status.name, Width: status},
		{Title: mergeReqsCols.mergeStatus.name, Width: status},
		{Title: mergeReqsCols.draft.name, Width: icon},
		{Title: mergeReqsCols.confilcts.name, Width: icon},
		{Title: mergeReqsCols.url.name, Width: url},
		{Title: mergeReqsCols.desc.name, Width: 0},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(r),
		table.WithFocused(true),
		table.WithHeight(len(r)),
		table.WithStyles(style.Table),
	)

	return t
}

func SetMergeRequestsCommentsModel(msg []table.Row) table.Model {
	columns := []table.Column{
		{Title: commentsCols.id.name, Width: 10},
		{Title: commentsCols.commentType.name, Width: 20},
		{Title: commentsCols.author.name, Width: 20},
		{Title: commentsCols.createdAt.name, Width: 30},
		{Title: commentsCols.updatedAt.name, Width: 30},
		{Title: commentsCols.resolved.name, Width: 10},
		{Title: commentsCols.body.name, Width: 0},
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
		{Title: pipelinesCols.id.name, Width: 10},
		{Title: pipelinesCols.iid.name, Width: 20},
		{Title: pipelinesCols.status.name, Width: 20},
		{Title: pipelinesCols.source.name, Width: 20},
		{Title: pipelinesCols.createdAt.name, Width: 30},
		{Title: pipelinesCols.updatedAt.name, Width: 30},
		{Title: pipelinesCols.url.name, Width: 0},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(msg),
		table.WithFocused(true),
		table.WithHeight(len(msg)),
	)

	return t
}
