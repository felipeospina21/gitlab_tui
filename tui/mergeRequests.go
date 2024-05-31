package tui

import (
	"gitlab_tui/internal/style"
	tbl "gitlab_tui/tui/components/table"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
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

func (m Model) InitMergeRequestsModel(listModel table.Model, commentsModel table.Model, pipelinesModel table.Model, projectsModel list.Model) Model {
	newM := Model{
		MergeRequests: MergeRequestsModel{
			List:         listModel,
			Comments:     commentsModel,
			Pipeline:     pipelinesModel,
			ListKeys:     MergeReqsKeys,
			CommentsKeys: CommentsKeys,
			PipelineKeys: PipelinKeys,
		},
		Projects: ProjectsModel{
			List:      projectsModel,
			ProjectID: m.Projects.ProjectID,
		},
		CurrView: m.CurrView,
		Md:       m.Md,
		Help:     m.Help,
	}

	return newM
}

func (m Model) SetMergeRequestsCommentsModel(msg []table.Row) table.Model {
	return tbl.InitModel(tbl.InitModelParams{
		Rows:      msg,
		Colums:    tbl.GetCommentsColums(m.Window.Width),
		StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.CommentsIconCols),
	})
}

func (m Model) SetMergeRequestPipelinesModel(msg []table.Row) table.Model {
	return tbl.InitModel(tbl.InitModelParams{
		Rows:      msg,
		Colums:    tbl.GetPipelinesColums(m.Window.Width),
		StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.PipelinesIconCols),
	})
}
