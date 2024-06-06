package table

import "gitlab_tui/internal/icon"

type mergeReqsPipelinesTable struct {
	ID        TableCol
	Iid       TableCol
	Status    TableCol
	Source    TableCol
	CreatedAt TableCol
	UpdatedAt TableCol
	URL       TableCol
}

var PipelinesCols = mergeReqsPipelinesTable{
	ID:        TableCol{Idx: 0, Name: "Id"},
	Iid:       TableCol{Idx: 1, Name: "IID"},
	CreatedAt: TableCol{Idx: 2, Name: icon.Clock},
	Status:    TableCol{Idx: 3, Name: "Status"},
	Source:    TableCol{Idx: 4, Name: "Source"},
	URL:       TableCol{Idx: 5, Name: "URL"},
}

var PipelinesIconCols = []int{}

func GetPipelinesColums(width int) []Column {
	id := 0
	iid := 0
	status := int(float32(width) * 0.25)
	source := int(float32(width) * 0.25)
	created := int(float32(width) * 0.1)
	url := 0

	columns := []Column{
		{Title: PipelinesCols.ID.Name, Width: id},
		{Title: PipelinesCols.Iid.Name, Width: iid},
		{Title: PipelinesCols.CreatedAt.Name, Width: created},
		{Title: PipelinesCols.Status.Name, Width: status},
		{Title: PipelinesCols.Source.Name, Width: source},
		{Title: PipelinesCols.URL.Name, Width: url},
	}

	return columns
}
