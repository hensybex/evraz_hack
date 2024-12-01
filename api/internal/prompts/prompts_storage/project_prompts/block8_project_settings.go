// internal/prompts/prompts_storage/project_prompts/block8_project_settings.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
	"fmt"
)

type ProjectSettingsData struct {
	SettingsFilesContent map[string]string // Key: filename, Value: content
	EnvironmentVariables map[string]string // Key: variable name, Value: value
	CompositeModules     []string          // List of composite module names
}

func (d ProjectSettingsData) ToPassedData() []types.PassedData {
	var passedData []types.PassedData
	for filename, content := range d.SettingsFilesContent {
		passedData = append(passedData, types.PassedData{
			Name:        fmt.Sprintf("Content of %s", filename),
			Description: fmt.Sprintf("Contents of the settings file %s", filename),
			Content:     content,
		})
	}
	passedData = append(passedData, types.PassedData{
		Name:        "Environment Variables",
		Description: "List of environment variables and their values",
		Content:     fmt.Sprintf("%v", d.EnvironmentVariables),
	})
	passedData = append(passedData, types.PassedData{
		Name:        "Composite Modules",
		Description: "List of composite module names",
		Content:     fmt.Sprintf("%v", d.CompositeModules),
	})
	return passedData
}

var ProjectSettingsPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in configuration management, please review the project's settings and configurations.",
	BaseTaskDesc: "Ensure that project settings and configurations adhere to standards.\n\nGuidelines:\n\nSettings are passed via environment variables.\nEach component has its own settings.py using Pydantic's BaseSettings.\nComposite modules in composites assemble components and manage dependency injection.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the settings meet the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
