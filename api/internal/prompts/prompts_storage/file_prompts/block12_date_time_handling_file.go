// internal/prompts/prompts_storage/file_prompts/block12_date_time_handling_file.go

package file_prompts

import "evraz_api/internal/prompts/types"

type DateTimeHandlingFileData struct {
	FilePath    string
	FileContent string
}

func (d DateTimeHandlingFileData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "File Path",
			Description: "Path of the source code file handling date and time",
			Content:     d.FilePath,
		},
		{
			Name:        "File Content",
			Description: "Contents of the source code file",
			Content:     d.FileContent,
		},
	}
}

var DateTimeHandlingFilePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in date and time handling in software applications, please review the following code.",
	BaseTaskDesc: "Assess the handling of date and time data within the code.\n\nGuidelines:\n\nUse UTC for all datetime operations.\nProperly manage timezone-aware and naive datetime objects.\nEnsure date conversions are handled correctly.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether date and time handling meets the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
