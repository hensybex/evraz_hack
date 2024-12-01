// internal/prompts/prompts_storage/file_prompts/master_agent.go

package file_prompts

import (
	"evraz_api/internal/prompts/types"
)

type FileMasterData struct {
	ProjectTree string
	FilePath    string
	FileContent string
}

func (t FileMasterData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "File content", Description: "Full file content", Content: t.FileContent},
	}
}

var FileMasterPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant overseeing the file analysis, analyse provided python file, and determine whether it is related to application layer or adapters layer of the project",
	BaseTaskDesc: "Return a structured JSON, with a single int value. If the file is related to application layer, return 0, if it is related to adapters layer, return 1, otherwise - return 2",
	JSONStruct: []types.JSONStruct{
		{Key: "value", Description: "An int value based on file type"},
	},
}
