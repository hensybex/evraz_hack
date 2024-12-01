// internal/prompts/prompts_storage/project_prompts/master_agent.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type ProjectMasterData struct {
	FileContent string
}

func (t ProjectMasterData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "File content", Description: "Full file content", Content: t.FileContent},
	}
}

var ProjectMasterPrompt = types.Prompt{
	BasePrompt: "As an AI assistant overseeing the project analysis, review the provided file, and determine which review prompts are applicable based on the provided project metadata.",
	BaseTaskDesc: `Your goal is to return a JSON array containing the IDs or names of the applicable prompts for the project-level analysis. Here is a list of possible prompts:

	`,
	JSONStruct: []types.JSONStruct{
		{Key: "applicable_prompts", Description: "Array of prompt IDs that should be applied"},
	},
}
