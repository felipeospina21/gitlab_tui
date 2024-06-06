package tui

import (
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components/table"
)

type MergeRequestsModel struct {
	List         table.Model
	Comments     table.Model
	Pipeline     table.Model
	SelectedMr   string
	Error        error
	ListKeys     MergeReqsKeyMap
	CommentsKeys CommentsKeyMap
	PipelineKeys PipelineKeyMap
}

func (m Model) SetMergeRequestsCommentsModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetCommentsColums(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(style.Table()), table.CommentsIconCols),
	})
}

func (m Model) SetMergeRequestPipelinesModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetPipelinesColums(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(style.Table()), table.PipelinesIconCols),
	})
}
