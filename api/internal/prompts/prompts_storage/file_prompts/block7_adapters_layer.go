// internal/prompts/prompts_storage/file_prompts/block7_adapters_layer.go

package file_prompts

import "evraz_api/internal/prompts/types"

type AdaptersLayerCodeData struct {
	FilePath    string
	FileContent string
}

func (d AdaptersLayerCodeData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "File Path",
			Description: "Path of the adapters layer source code file",
			Content:     d.FilePath,
		},
		{
			Name:        "File Content",
			Description: "Contents of the adapters layer source code file",
			Content:     d.FileContent,
		},
	}
}

var AdaptersLayerCodePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software architecture, please review the following adapters layer code.",
	BaseTaskDesc: "Review the adapters layer code for compliance with guidelines.\n\nGuidelines:\n\nManages integrations with external systems.\nContains web frameworks, CLI tools, and API clients.\nHandles database interactions using SQLAlchemy.\nAvoids embedding business logic in query code.\nControllers inject services from the application layer.\nPrepares data for serialization; manages asynchronous tasks.\nFollows serialization rules for specific data types.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the code meets the adapters layer requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
