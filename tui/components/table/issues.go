package table

import "gitlab_tui/internal/icon"

type issuesListTable struct {
	CreatedAt  TableCol
	Title      TableCol
	Author     TableCol
	Assignees  TableCol
	ClosedBy   TableCol
	Labels     TableCol
	NotesCount TableCol
	URL        TableCol
	Desc       TableCol
	ID         TableCol
}

var IssuesListCols = issuesListTable{
	CreatedAt:  TableCol{Idx: 0, Name: icon.Clock},
	Title:      TableCol{Idx: 1, Name: "Title"},
	Assignees:  TableCol{Idx: 2, Name: "Assignees"},
	Labels:     TableCol{Idx: 3, Name: "Labels"},
	NotesCount: TableCol{Idx: 4, Name: "Notes"},
	Author:     TableCol{Idx: 5, Name: "Author"},
	ClosedBy:   TableCol{Idx: 6, Name: "Closed By"},
	URL:        TableCol{Idx: 7, Name: "URL"},
	Desc:       TableCol{Idx: 8, Name: "Desc"},
	ID:         TableCol{Idx: 9, Name: "ID"},
}

var IssuesListIconCols = []int{}

func GetIssuesListColumns(width int) []Column {
	i := int(float32(width) * 0.03)
	title := int(float32(width) * 0.4)
	author := int(float32(width) * 0.08)
	array := int(float32(width) * 0.15)

	columns := []Column{
		{Title: IssuesListCols.CreatedAt.Name, Width: i},
		{Title: IssuesListCols.Title.Name, Width: title},
		{Title: IssuesListCols.Assignees.Name, Width: array},
		{Title: IssuesListCols.Labels.Name, Width: array},
		{Title: IssuesListCols.NotesCount.Name, Width: i},
		{Title: IssuesListCols.Author.Name, Width: author},
		{Title: IssuesListCols.ClosedBy.Name, Width: author},
		{Title: IssuesListCols.URL.Name, Width: 0},
		{Title: IssuesListCols.Desc.Name, Width: 0},
		{Title: IssuesListCols.ID.Name, Width: 0},
	}

	return columns
}
