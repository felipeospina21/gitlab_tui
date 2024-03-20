package config

import (
	"log"

	"github.com/joho/godotenv"
)

type envVars = map[string]string

type projectsId struct {
	PlanningTool string
}
type config struct {
	BaseUrl    string
	ApiToken   string
	ApiVersion string
	ProjectsId projectsId
}

var Config config

func Load(configObj *config) {
	ev := readEnvVars()

	configObj.ApiVersion = "v4"
	configObj.BaseUrl = ev["BASE_URL"]
	configObj.ApiToken = ev["TOKEN"]
	configObj.ProjectsId = projectsId{
		PlanningTool: ev["PLANNING_TOOL_ID"],
	}
}

func readEnvVars() envVars {
	var m envVars
	m, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading env")
	}
	return m
}
