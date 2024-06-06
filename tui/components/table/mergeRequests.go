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
	URL         TableCol
	Desc        TableCol
	ID          TableCol
}

var MergeReqsCols = mergeReqsTable{
	CreatedAd:   TableCol{Idx: 0, Name: icon.Clock},
	Title:       TableCol{Idx: 1, Name: "Title"},
	Author:      TableCol{Idx: 2, Name: "Author"},
	MergeStatus: TableCol{Idx: 3, Name: "Merge Status"},
	Draft:       TableCol{Idx: 4, Name: "Draft"},
	Confilcts:   TableCol{Idx: 5, Name: "Conflicts"},
	URL:         TableCol{Idx: 6, Name: "Url"},
	Desc:        TableCol{Idx: 7, Name: "Description"},
	ID:          TableCol{Idx: 8, Name: "Id"},
}

var MergeReqsIconCols = []int{
	int(MergeReqsCols.MergeStatus.Idx),
	int(MergeReqsCols.Draft.Idx),
	int(MergeReqsCols.Confilcts.Idx),
}

func GetMergeReqsColums(width int) []Column {
	id := int(float32(width) * 0.06)
	title := int(float32(width) * 0.5)
	author := int(float32(width) * 0.2)
	status := int(float32(width) * 0.1)
	i := int(float32(width) * 0.04)
	url := 0

	columns := []Column{
		{Title: MergeReqsCols.CreatedAd.Name, Width: id},
		{Title: MergeReqsCols.Title.Name, Width: title},
		{Title: MergeReqsCols.Author.Name, Width: author},
		{Title: MergeReqsCols.MergeStatus.Name, Width: status},
		{Title: MergeReqsCols.Draft.Name, Width: i},
		{Title: MergeReqsCols.Confilcts.Name, Width: i},
		{Title: MergeReqsCols.URL.Name, Width: url},
		{Title: MergeReqsCols.Desc.Name, Width: 0},
		{Title: MergeReqsCols.ID.Name, Width: 0},
	}

	return columns
}
