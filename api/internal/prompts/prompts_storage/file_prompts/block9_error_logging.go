// internal/prompts/prompts_storage/file_prompts/block9_error_logging.go

package file_prompts

import "evraz_api/internal/prompts/types"

type ErrorHandlingAndLoggingData struct {
	FilePath    string
	FileContent string
}

func (d ErrorHandlingAndLoggingData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "File Path",
			Description: "Path of the source code file implementing error handling and logging",
			Content:     d.FilePath,
		},
		{
			Name:        "File Content",
			Description: "Contents of the source code file",
			Content:     d.FileContent,
		},
	}
}

var ErrorHandlingAndLoggingPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in error handling and logging practices, please review the following code.",
	BaseTaskDesc: "Verify that error handling and logging adhere to standards.\n\nGuidelines:\n\nErrors are defined in the business logic layer.\nServices perform validation and raise custom errors.\nAdapters catch errors and format responses appropriately.\nUses the standard logging module.\nConfigured in settings.py; logs in JSON format using python-json-logger.\nLoggers are module-level; avoid global loggers.\nUse %s placeholders in logging statements instead of f-strings.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether error handling and logging meet the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
