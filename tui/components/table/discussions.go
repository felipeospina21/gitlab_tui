package table

import "gitlab_tui/internal/icon"

type mergeReqsDiscussionsTable struct {
	ID          TableCol
	Author      TableCol
	CreatedAt   TableCol
	Body        TableCol
	Discussions TableCol
	Type        TableCol
	IsResolved  TableCol
}

var DiscussionsCols = mergeReqsDiscussionsTable{
	CreatedAt:   TableCol{Idx: 0, Name: icon.Clock},
	Author:      TableCol{Idx: 1, Name: "Author"},
	Discussions: TableCol{Idx: 2, Name: "Discussions"},
	IsResolved:  TableCol{Idx: 3, Name: "Resolved"},
	Type:        TableCol{Idx: 4, Name: "Type"},
	Body:        TableCol{Idx: 5, Name: "Body"},
	ID:          TableCol{Idx: 6, Name: "Id"},
}

var DiscussionsIconCols = []int{}

func GetDiscussionsColums(width int) []Column {
	cols := int(float32(width) * 0.1)
	body := int(float32(width) * 0.5)
	id := 0

	columns := []Column{
		{Title: DiscussionsCols.CreatedAt.Name, Width: cols},
		{Title: DiscussionsCols.Author.Name, Width: cols},
		{Title: DiscussionsCols.Discussions.Name, Width: cols / 2},
		{Title: DiscussionsCols.IsResolved.Name, Width: cols / 2},
		{Title: DiscussionsCols.Type.Name, Width: cols / 2},
		{Title: DiscussionsCols.Body.Name, Width: body},
		{Title: DiscussionsCols.ID.Name, Width: id},
	}

	return columns
}
