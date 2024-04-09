package tui

import "github.com/charmbracelet/bubbles/table"

// Model to list pipelines
type PipelinesModel struct {
	List  table.Model
	Jobs  table.Model
	Error error
}

// pipelines https://docs.gitlab.com/ee/api/pipelines.html

// jobs https://docs.gitlab.com/ee/api/jobs.html#list-pipeline-jobs
