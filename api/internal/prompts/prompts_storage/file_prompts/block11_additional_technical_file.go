// internal/prompts/prompts_storage/file_prompts/block11_additional_technical_file.go

package file_prompts

import "evraz_api/internal/prompts/types"

type AdditionalTechnicalFileData struct {
	FilePath    string
	FileContent string
}

func (d AdditionalTechnicalFileData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "File Path",
			Description: "Path of the source code file",
			Content:     d.FilePath,
		},
		{
			Name:        "File Content",
			Description: "Contents of the source code file",
			Content:     d.FileContent,
		},
	}
}

var AdditionalTechnicalFilePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software best practices, please review the following code for technical considerations.",
	BaseTaskDesc: "Ensure adherence to technical requirements at the file level.\n\nGuidelines:\n\nDatabase Transactions: Avoid manual transaction management within code.\nAsynchronous Code: Use only when justified; ensure proper implementation.\nData Science Dependencies: Use pandas, numpy, etc., only within data science modules.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the code meets the additional technical requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
