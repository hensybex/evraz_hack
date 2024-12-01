// internal/prompts/prompts_storage/project_prompts/block2_key_files.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
	"fmt"
)

type KeyFilesData struct {
	SetupFilesContent map[string]string // Key: filename, Value: content
	ReadmeContent     string
	SourceDirName     string
}

func (d KeyFilesData) ToPassedData() []types.PassedData {
	var passedData []types.PassedData
	for filename, content := range d.SetupFilesContent {
		passedData = append(passedData, types.PassedData{
			Name:        fmt.Sprintf("Content of %s", filename),
			Description: fmt.Sprintf("Contents of the file %s", filename),
			Content:     content,
		})
	}
	passedData = append(passedData, types.PassedData{
		Name:        "README.md",
		Description: "Contents of README.md",
		Content:     d.ReadmeContent,
	})
	passedData = append(passedData, types.PassedData{
		Name:        "Source Directory Name",
		Description: "Name of the source code directory",
		Content:     d.SourceDirName,
	})
	return passedData
}

var KeyFilesPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in code analysis, please review the key configuration files of the project.",
	BaseTaskDesc: "Ensure that the backend directory contains the essential files with correct configurations.\n\nGuidelines:\n\nsetup.py or setup.cfg: Contains package metadata and dependencies.\npyproject.toml: Includes configurations for builders and autoformatters.\nREADME.md: Provides a project overview, deployment instructions, testing procedures, and permission/group schemes.\nSource Code Directory: Acts as the root for imports and has a concise, meaningful name.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the key files meet the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
