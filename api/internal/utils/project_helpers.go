// internal/utils/project_helpers.go

package utils

import (
	"evraz_api/internal/model"
	"strings"
)

func CheckIfProjectHasTests(project *model.Project) bool {
	// Implement logic to check if the project has tests
	return strings.Contains(project.Tree, "tests")
}

func CheckIfProjectUsesAsync(project *model.Project) bool {
	// Implement logic to check if the project uses asynchronous code
	// For example, search through project files for "async" or "await"
	return strings.Contains(project.Tree, "async") || strings.Contains(project.Tree, "await")
}

func CheckIfProjectContainsDatetime(project *model.Project) bool {
	// Implement logic to check if the project handles date and time operations
	return strings.Contains(project.Tree, "datetime") || strings.Contains(project.Tree, "time")
}
