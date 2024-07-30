package table

import (
	"gitlab_tui/internal/icon"
)

type mergeReqsTable struct {
	CreatedAd   TableCol
	Title       TableCol
	Author      TableCol
	MergeStatus TableCol
	Draft       TableCol
	Confilcts   TableCol
	Comments    TableCol
	URL         TableCol
	Desc        TableCol
	ID          TableCol
}

var MergeReqsCols = mergeReqsTable{
	CreatedAd:   TableCol{Idx: 0, Name: icon.Clock},
	Draft:       TableCol{Idx: 1, Name: ""},
	Title:       TableCol{Idx: 2, Name: "Title"},
	Author:      TableCol{Idx: 3, Name: "Author"},
	MergeStatus: TableCol{Idx: 4, Name: "Merge Status"},
	Confilcts:   TableCol{Idx: 5, Name: "Conflicts"},
	Comments:    TableCol{Idx: 6, Name: "Comments"},
	URL:         TableCol{Idx: 7, Name: "Url"},
	Desc:        TableCol{Idx: 8, Name: "Description"},
	ID:          TableCol{Idx: 9, Name: "Id"},
}

var MergeReqsIconCols = []int{
	int(MergeReqsCols.MergeStatus.Idx),
	int(MergeReqsCols.Draft.Idx),
	int(MergeReqsCols.Confilcts.Idx),
}

func GetMergeReqsColums(width int) []Column {
	id := int(float32(width) * 0.04)
	title := int(float32(width) * 0.5)
	author := int(float32(width) * 0.2)
	status := int(float32(width) * 0.1)
	i := int(float32(width) * 0.04)
	url := 0

	columns := []Column{
		{Title: MergeReqsCols.CreatedAd.Name, Width: id},
		{Title: MergeReqsCols.Draft.Name, Width: i},
		{Title: MergeReqsCols.Title.Name, Width: title},
		{Title: MergeReqsCols.Author.Name, Width: author},
		{Title: MergeReqsCols.MergeStatus.Name, Width: status},
		{Title: MergeReqsCols.Confilcts.Name, Width: i},
		{Title: MergeReqsCols.Comments.Name, Width: i},
		{Title: MergeReqsCols.URL.Name, Width: url},
		{Title: MergeReqsCols.Desc.Name, Width: 0},
		{Title: MergeReqsCols.ID.Name, Width: 0},
	}

	return columns
}
