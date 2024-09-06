package table

import "gitlab_tui/internal/icon"

type mergeReqsDiscussionsTable struct {
	ID        TableCol
	Author    TableCol
	CreatedAt TableCol
	Body      TableCol
	Comments  TableCol
	Type      TableCol
}

var DiscussionsCols = mergeReqsDiscussionsTable{
	CreatedAt: TableCol{Idx: 0, Name: icon.Clock},
	Author:    TableCol{Idx: 1, Name: "Author"},
	Type:      TableCol{Idx: 2, Name: "Type"},
	Comments:  TableCol{Idx: 3, Name: "Comments"},
	Body:      TableCol{Idx: 4, Name: "Body"},
	ID:        TableCol{Idx: 5, Name: "Id"},
}

var DiscussionsIconCols = []int{}

func GetDiscussionsColums(width int) []Column {
	t := int(float32(width) * 0.25)
	author := int(float32(width) * 0.3)
	created := int(float32(width) * 0.1)
	body := int(float32(width) * 0.3)
	id := 0

	columns := []Column{
		{Title: DiscussionsCols.CreatedAt.Name, Width: created},
		{Title: DiscussionsCols.Author.Name, Width: author},
		{Title: DiscussionsCols.Type.Name, Width: created},
		{Title: DiscussionsCols.Comments.Name, Width: t},
		{Title: DiscussionsCols.Body.Name, Width: body},
		{Title: DiscussionsCols.ID.Name, Width: id},
	}

	return columns
}
