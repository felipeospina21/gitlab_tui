package tui

import (
	"gitlab_tui/tui/components/table"
)

type MergeRequestsModel struct {
	List         table.Model
	Discussions  table.Model
	Pipeline     table.Model
	PipelineJobs table.Model
	SelectedMr   string
}

func (m Model) SetMergeRequestsDiscussionsModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetDiscussionsColums(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.DiscussionsIconCols),
	})
}

func (m Model) SetMergeRequestPipelinesModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetPipelinesColums(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.PipelinesIconCols),
	})
}

func (m Model) SetPipelineJobsModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetPipelineJobsColums(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(table.DefaultStyle()), table.PipelineJobsIconCols),
	})
}
