package utils

import (
	"fmt"
	"os"
)

func GetProjectUrl(projectId string) string {
	environment := os.Getenv("ENVIRONMENT")
	switch environment {
	case "PROD":
		return fmt.Sprintf("https://buildium-frontend-amd.onrender.com/projects/%s", projectId)
	default:
		return fmt.Sprintf("http://localhost:5173/projects/%s", projectId)
	}
}
