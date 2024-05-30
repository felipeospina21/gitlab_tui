package table

import (
	"gitlab_tui/internal/icon"
	"gitlab_tui/internal/style"
	"slices"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type TableColIndex uint

type tableCol struct {
	Name string
	Idx  TableColIndex
}

type mergeReqsTable struct {
	Iid         tableCol
	Title       tableCol
	Author      tableCol
	Status      tableCol
	MergeStatus tableCol
	Draft       tableCol
	Confilcts   tableCol
	URL         tableCol
	Desc        tableCol
}

type mergeReqsCommentsTable struct {
	ID          tableCol
	CommentType tableCol
	Author      tableCol
	CreatedAt   tableCol
	UpdatedAt   tableCol
	Resolved    tableCol
	Body        tableCol
}

type mergeReqsPipelinesTable struct {
	ID        tableCol
	Iid       tableCol
	Status    tableCol
	Source    tableCol
	CreatedAt tableCol
	UpdatedAt tableCol
	URL       tableCol
}

var MergeReqsCols = mergeReqsTable{
	Iid:         tableCol{Idx: 0, Name: "Iid"},
	Title:       tableCol{Idx: 1, Name: "Title"},
	Author:      tableCol{Idx: 2, Name: "Author"},
	Status:      tableCol{Idx: 3, Name: "Status"},
	MergeStatus: tableCol{Idx: 4, Name: "Merge Status"},
	Draft:       tableCol{Idx: 5, Name: "Draft"},
	Confilcts:   tableCol{Idx: 6, Name: "Conflicts"},
	URL:         tableCol{Idx: 7, Name: "Url"},
	Desc:        tableCol{Idx: 8, Name: "Description"},
}

var CommentsCols = mergeReqsCommentsTable{
	ID:          tableCol{Idx: 0, Name: "Id"},
	CommentType: tableCol{Idx: 1, Name: "Type"},
	Author:      tableCol{Idx: 2, Name: "Author"},
	CreatedAt:   tableCol{Idx: 3, Name: "Created At"},
	UpdatedAt:   tableCol{Idx: 4, Name: "Updated At"},
	Resolved:    tableCol{Idx: 5, Name: "Resolved"},
	Body:        tableCol{Idx: 6, Name: "Body"},
}

var PipelinesCols = mergeReqsPipelinesTable{
	ID:        tableCol{Idx: 0, Name: "Id"},
	Iid:       tableCol{Idx: 1, Name: "IID"},
	Status:    tableCol{Idx: 2, Name: "Status"},
	Source:    tableCol{Idx: 3, Name: "Source"},
	CreatedAt: tableCol{Idx: 4, Name: "Created At"},
	UpdatedAt: tableCol{Idx: 5, Name: "Updated At"},
	URL:       tableCol{Idx: 6, Name: "URL"},
}

var (
	MergeReqsIconCols = []int{
		int(MergeReqsCols.MergeStatus.Idx),
		int(MergeReqsCols.Draft.Idx),
		int(MergeReqsCols.Confilcts.Idx),
	}

	CommentsIconCols = []int{
		int(CommentsCols.Resolved.Idx),
	}

	PipelinesIconCols = []int{}
)

type Model struct {
	Table  table.Model
	Colums []table.Column
}

func NewModel() Model {
	m := Model{
		Table:  table.Model{},
		Colums: []table.Column{},
	}

	return m
}

type InitModelParams struct {
	Rows      []table.Row
	Colums    []table.Column
	StyleFunc table.StyleFunc
}

func InitModel(params InitModelParams) table.Model {
	s := style.Table()

	t := table.New(
		table.WithColumns(params.Colums),
		table.WithRows(params.Rows),
		table.WithFocused(true),
		table.WithHeight(len(params.Rows)+1),
		table.WithStyles(s),
		table.WithStyleFunc(params.StyleFunc),
	)

	return t
}

func GetMergeReqsColums(width int) []table.Column {
	id := int(float32(width) * 0.06)
	title := int(float32(width) * 0.45)
	author := int(float32(width) * 0.1)
	status := int(float32(width) * 0.13)
	i := int(float32(width) * 0.04)
	url := 0

	if width > 170 {
		id = int(float32(width) * 0.03)
		title = int(float32(width) * 0.35)
		status = int(float32(width) * 0.1)
		total := id + title + author + (status * 2) + (i * 2)
		url = width - total - 10
	}

	columns := []table.Column{
		{Title: MergeReqsCols.Iid.Name, Width: id},
		{Title: MergeReqsCols.Title.Name, Width: title},
		{Title: MergeReqsCols.Author.Name, Width: author},
		{Title: MergeReqsCols.Status.Name, Width: status},
		{Title: MergeReqsCols.MergeStatus.Name, Width: status},
		{Title: MergeReqsCols.Draft.Name, Width: i},
		{Title: MergeReqsCols.Confilcts.Name, Width: i},
		{Title: MergeReqsCols.URL.Name, Width: url},
		{Title: MergeReqsCols.Desc.Name, Width: 0},
	}

	return columns
}

func GetCommentsColums(width int) []table.Column {
	id := int(float32(width) * 0.1)
	t := int(float32(width) * 0.25)
	author := int(float32(width) * 0.3)
	created := int(float32(width) * 0.1)
	updated := int(float32(width) * 0.1)
	resloved := int(float32(width) * 0.1)
	body := 0

	columns := []table.Column{
		{Title: CommentsCols.ID.Name, Width: id},
		{Title: CommentsCols.CommentType.Name, Width: t},
		{Title: CommentsCols.Author.Name, Width: author},
		{Title: CommentsCols.CreatedAt.Name, Width: created},
		{Title: CommentsCols.UpdatedAt.Name, Width: updated},
		{Title: CommentsCols.Resolved.Name, Width: resloved},
		{Title: CommentsCols.Body.Name, Width: body},
	}

	return columns
}

func GetPipelinesColums(width int) []table.Column {
	id := int(float32(width) * 0.1)
	iid := int(float32(width) * 0.1)
	status := int(float32(width) * 0.25)
	source := int(float32(width) * 0.25)
	created := int(float32(width) * 0.1)
	updated := int(float32(width) * 0.1)
	url := 0

	columns := []table.Column{
		{Title: PipelinesCols.ID.Name, Width: id},
		{Title: PipelinesCols.Iid.Name, Width: iid},
		{Title: PipelinesCols.Status.Name, Width: status},
		{Title: PipelinesCols.Source.Name, Width: source},
		{Title: PipelinesCols.CreatedAt.Name, Width: created},
		{Title: PipelinesCols.UpdatedAt.Name, Width: updated},
		{Title: PipelinesCols.URL.Name, Width: url},
	}

	return columns
}

func StyleIconsColumns(s table.Styles, iconColIdx []int) table.StyleFunc {
	return func(row, col int, value string) lipgloss.Style {
		isIconCol := slices.Contains(iconColIdx, col)

		if isIconCol {
			switch value {
			case icon.Check:
				return s.Cell.Foreground(lipgloss.Color(style.Green[300]))

			case icon.Clock:
				return s.Cell.Foreground(lipgloss.Color(style.Yellow[300]))
			}
		}
		return s.Cell
	}
}
