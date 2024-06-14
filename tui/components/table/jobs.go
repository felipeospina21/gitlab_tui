package table

import "gitlab_tui/internal/icon"

type pipelineJobsTable struct {
	ID        TableCol
	Status    TableCol
	Stage     TableCol
	Name      TableCol
	CreatedAt TableCol
	Duration  TableCol
	Coverage  TableCol
	URL       TableCol
}

var PipelineJobsCols = pipelineJobsTable{
	ID:        TableCol{Idx: 0, Name: "Id"},
	CreatedAt: TableCol{Idx: 1, Name: icon.Clock},
	Status:    TableCol{Idx: 2, Name: "Status"},
	Name:      TableCol{Idx: 3, Name: "Name"},
	Stage:     TableCol{Idx: 4, Name: "Stage"},
	Duration:  TableCol{Idx: 5, Name: "Duration"},
	Coverage:  TableCol{Idx: 6, Name: "Coverage"},
	URL:       TableCol{Idx: 7, Name: "URL"},
}

var PipelineJobsIconCols = []int{
	int(PipelineJobsCols.Status.Idx),
}

func GetPipelineJobsColums(width int) []Column {
	id := 0
	created := int(float32(width) * 0.1)
	status := int(float32(width) * 0.20)
	name := int(float32(width) * 0.25)
	duration := int(float32(width) * 0.1)
	url := 0

	columns := []Column{
		{Title: PipelineJobsCols.ID.Name, Width: id},
		{Title: PipelineJobsCols.CreatedAt.Name, Width: created},
		{Title: PipelineJobsCols.Status.Name, Width: status},
		{Title: PipelineJobsCols.Name.Name, Width: name},
		{Title: PipelineJobsCols.Stage.Name, Width: name},
		{Title: PipelineJobsCols.Duration.Name, Width: duration},
		{Title: PipelineJobsCols.Coverage.Name, Width: duration},
		{Title: PipelineJobsCols.URL.Name, Width: url},
	}

	return columns
}
