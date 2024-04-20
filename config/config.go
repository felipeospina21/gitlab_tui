package config

import (
	"log"

	"github.com/joho/godotenv"
)

type envVars = map[string]string

type projectsID struct {
	PlanningTool string
}
type config struct {
	BaseURL    string
	APIToken   string
	APIVersion string
	ProjectsID projectsID
}

var Config config

func Load(configObj *config) {
	ev := readEnvVars()

	configObj.APIVersion = "v4"
	configObj.BaseURL = ev["BASE_URL"]
	configObj.APIToken = ev["TOKEN"]
	configObj.ProjectsID = projectsID{
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
