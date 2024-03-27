package main

import (
	"gitlab_tui/api"

	"github.com/charmbracelet/bubbles/table"
)

type mergeRequests struct {
	model      table.Model
	selectedMr string
}

func setMergeRequestsModel() table.Model {
	r := api.GetMergeRequests()

	columns := []table.Column{
		{Title: "Iid", Width: 4},
		{Title: "Title", Width: 80},
		{Title: "Author", Width: 20},
		{Title: "Status", Width: 30},
		{Title: "Draft", Width: 10},
		{Title: "Conflicts", Width: 10},
		{Title: "Url", Width: 0},
		{Title: "Description", Width: 0},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(r),
		table.WithFocused(true),
		table.WithHeight(len(r)),
	)

	return t
}
