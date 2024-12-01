// internal/prompts/prompts_storage/project_prompts/block5_dependency_management.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type DependencyManagementData struct {
	DependenciesContent string // Contents of dependency files
}

func (d DependencyManagementData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Dependencies File",
			Description: "File with all dependencies with their versions",
			Content:     d.DependenciesContent,
		},
	}
}

var DependencyManagementPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in dependency management, please review the project's dependencies.",
	BaseTaskDesc: "Verify that the project uses the correct dependencies as per the specified stack.\n\nGuidelines:\n\nEnsure the latest versions of evraz-classic packages are used.\nCheck that development packages match the specified versions.\nConfirm no unauthorized packages are included without approval.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the dependencies meet the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues with dependencies"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
