package main

import (
	"github.com/charmbracelet/bubbles/table"
)

type mergeRequestsComments struct {
	model table.Model
}

// func (m *model) setMergeRequestsCommentsModel() table.Model {
// 	r := api.GetMRComments(m.mergeRequests.selectedMr)
//
// 	columns := []table.Column{
// 		{Title: "Id", Width: 4},
// 		{Title: "Type", Width: 10},
// 		{Title: "Author", Width: 20},
// 		{Title: "Created At", Width: 20},
// 		{Title: "Updated At", Width: 20},
// 		{Title: "Resolved", Width: 10},
// 		{Title: "Body", Width: 0},
// 	}
//
// 	t := table.New(
// 		table.WithColumns(columns),
// 		table.WithRows(r),
// 		table.WithFocused(true),
// 		table.WithHeight(len(r)),
// 	)
//
// 	return t
// }
