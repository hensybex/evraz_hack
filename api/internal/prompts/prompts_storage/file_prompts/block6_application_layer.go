// internal/prompts/prompts_storage/file_prompts/block6_application_layer.go

package file_prompts

import "evraz_api/internal/prompts/types"

type ApplicationLayerCodeData struct {
	FilePath    string
	FileContent string
}

func (d ApplicationLayerCodeData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "File Path",
			Description: "Path of the application layer source code file",
			Content:     d.FilePath,
		},
		{
			Name:        "File Content",
			Description: "Contents of the application layer source code file",
			Content:     d.FileContent,
		},
	}
}

var ApplicationLayerCodePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software architecture, please review the following application layer code.",
	BaseTaskDesc: "Ensure the application layer code adheres to architectural principles.\n\nGuidelines:\n\nContains business logic elements (entities, DTOs, services).\nIs independent of adapters; uses Dependency Injection.\nDefines interfaces for data reception; adapters implement these interfaces.\nUses DTOs instead of simple data structures.\nPerforms data validation within services using Pydantic models.\nManages errors within this layer.\nAvoids excessive coupling between services.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the code meets the application layer requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
