package table

import "gitlab_tui/internal/icon"

type mergeReqsCommentsTable struct {
	ID          TableCol
	CommentType TableCol
	Author      TableCol
	CreatedAt   TableCol
	Resolved    TableCol
	Body        TableCol
}

var CommentsCols = mergeReqsCommentsTable{
	CreatedAt:   TableCol{Idx: 0, Name: icon.Clock},
	CommentType: TableCol{Idx: 1, Name: "Type"},
	Author:      TableCol{Idx: 2, Name: "Author"},
	Resolved:    TableCol{Idx: 3, Name: "Resolved"},
	Body:        TableCol{Idx: 4, Name: "Body"},
	ID:          TableCol{Idx: 5, Name: "Id"},
}

var CommentsIconCols = []int{
	int(CommentsCols.Resolved.Idx),
}

func GetCommentsColums(width int) []Column {
	t := int(float32(width) * 0.25)
	author := int(float32(width) * 0.3)
	created := int(float32(width) * 0.1)
	resloved := int(float32(width) * 0.1)
	body := 0
	id := 0

	columns := []Column{
		{Title: CommentsCols.CreatedAt.Name, Width: created},
		{Title: CommentsCols.CommentType.Name, Width: t},
		{Title: CommentsCols.Author.Name, Width: author},
		{Title: CommentsCols.Resolved.Name, Width: resloved},
		{Title: CommentsCols.Body.Name, Width: body},
		{Title: CommentsCols.ID.Name, Width: id},
	}

	return columns
}
