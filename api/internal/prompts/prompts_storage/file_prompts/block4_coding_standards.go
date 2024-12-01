// internal/prompts/prompts_storage/file_prompts/block4_coding_standards.go

package file_prompts

import "evraz_api/internal/prompts/types"

type CodingStandardsData struct {
	FilePath    string
	FileContent string
}

func (d CodingStandardsData) ToPassedData() []types.PassedData {
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

var CodingStandardsPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in code analysis, please review the following Python source code for adherence to coding standards.",
	BaseTaskDesc: "Check the code for adherence to coding standards.\n\nGuidelines:\n\nCode follows PEP8 guidelines.\nDocstrings comply with PEP256 and PEP257.\nUse yapf and isort for formatting.\nLine length does not exceed 80 characters (exceptions up to 100 characters with proper management).\nCode is decomposed and refactored for readability and maintainability.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the code meets coding standards"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
