// internal/prompts/prompts_storage/project_prompts/block3_application_architecture.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type ApplicationArchitectureData struct {
	ProjectStructure   string
	ModuleInteractions string
}

func (d ApplicationArchitectureData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Project Structure",
			Description: "Structure of the project directories and files",
			Content:     d.ProjectStructure,
		},
		{
			Name:        "Module Interactions",
			Description: "Description of how modules interact within the project",
			Content:     d.ModuleInteractions,
		},
	}
}

var ApplicationArchitecturePrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software architecture, please review the project's architecture.",
	BaseTaskDesc: "Evaluate whether the project follows the Hexagonal (Ports and Adapters) Architecture.\n\nGuidelines:\n\nApplication Core:\nDomain Layer: Contains business logic, independent of frameworks.\nApplication Layer: Manages use cases and workflows.\nPorts: Interfaces connecting the core to external systems.\nAdapters:\nPrimary Adapters: REST and WebSocket adapters for input.\nSecondary Adapters: Messaging queues, databases, SMS, and email service adapters for output.\nPrinciples:\nCore is independent of external technologies.\nAdapters bridge the core and external systems.\nThe architecture supports scalability and maintainability.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the architecture meets the requirements"},
		{Key: "issues", Description: "(list of str) List of any architectural issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
