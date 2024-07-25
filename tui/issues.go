package tui

import (
	"gitlab_tui/internal/style"
	"gitlab_tui/tui/components/table"
)

type IssuesModel struct {
	List          table.Model
	SelectedIssue string
	PrevPage      string
	NexPage       string
}

func (m Model) SetIssuesListModel(msg []table.Row) table.Model {
	return table.InitModel(table.InitModelParams{
		Rows:      msg,
		Colums:    table.GetIssuesListColumns(m.Window.Width),
		StyleFunc: table.StyleIconsColumns(table.Styles(style.Table()), table.IssuesListIconCols),
	})
}
