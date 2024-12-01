// internal/prompts/prompts_storage/file_prompts/block4_coding_standards.go

package helper_prompts

import (
	"evraz_api/internal/prompts/types"
)

type ExtractTestsData struct {
	ProjectTree string
}

func (d ExtractTestsData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Project Tree",
			Description: "The file structure of the project",
			Content:     d.ProjectTree,
		},
	}
}

var ExtractTestsPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in code analysis, analyze the following project structure to identify all files related to project testing. Consider any files that may have testing-related content (e.g., unit tests, integration tests, etc.). The files might be in specific directories or named with common testing patterns.",
	BaseTaskDesc: "Analyze the project directory structure and identify all file paths related to testing. Common testing-related file patterns include:\n\n- Files in directories like 'test', 'tests', 'spec', or 'integration'.\n- Files named with patterns such as 'test_', 'spec_', 'unit_', etc.\n- Files that may contain code related to tests (even if not explicitly named).\n\nReturn a structured JSON containing a list of full file paths for all files related to testing.",
	JSONStruct: []types.JSONStruct{
		{Key: "test_files_routes", Description: "List of full paths to testing-related files"},
	},
}
