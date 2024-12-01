// internal/prompts/prompts_storage/project_prompts/block1_project_structure.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type ProjectStructureData struct {
	ProjectTree       string
	RootFiles         []string
	ComponentsPresent bool
	DocsPresent       bool
	DeploymentPresent bool
}

func (d ProjectStructureData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Project Tree",
			Description: "Structure of the project directories and files",
			Content:     d.ProjectTree,
		},
	}
}

var ProjectStructurePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in code analysis, please analyze the following project structure.",
	BaseTaskDesc: "Verify that the project adheres to the required structural and organizational standards.\n\nGuidelines:\n\nMonorepository Structure: Confirm the project uses a monorepository layout similar to the demo project.\nRoot Files: Check for the presence of .gitignore, .editorconfig, and .gitattributes in the root directory.\nDirectories:\ndeployment: Contains CI/CD files (coordinate with DevOps if needed).\ndocs: Stores technical documentation, including PlantUML diagrams.\ncomponents: Separates frontend and backend code.\nWithin components, demo_project_backend should serve as the backend root.\nBackend Module:\nShould be recognized as the root for Python modules in IDEs (sources_root) and via PYTHONPATH.\nSwagger Documentation: Generated on the backend when the corresponding endpoint is called.\nBusiness Process Documentation: Maintained in the docs directory or a separate wiki.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the project meets the structural requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
