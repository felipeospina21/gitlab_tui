package config

import "fmt"

type routes struct {
	ProjectMergeReqs string
}

var Routes = routes{
	ProjectMergeReqs: fmt.Sprintf("%s/%s/projects/%s/merge_requests", Config.BaseURL, Config.APIVersion, Config.ProjectsID.PlanningTool),
}
