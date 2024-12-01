// internal/prompts/prompts.go

package prompts

import (
	"evraz_api/internal/prompts/prompts_storage/file_prompts"
	helper_prompts "evraz_api/internal/prompts/prompts_storage/helpers"
	"evraz_api/internal/prompts/prompts_storage/project_prompts"
	"evraz_api/internal/prompts/types"
)

type Prompts struct {
	// Master Prompts
	ProjectMasterPrompt types.Prompt
	FileMasterPrompt    types.Prompt

	// Project-Level Prompts
	ProjectStructure        types.Prompt
	KeyFiles                types.Prompt
	ApplicationArchitecture types.Prompt
	DependencyManagement    types.Prompt
	ProjectSettings         types.Prompt
	TestingStrategy         types.Prompt
	AdditionalTechnical     types.Prompt
	DateTimeHandling        types.Prompt

	// File-Level Prompts
	CodingStandards         types.Prompt
	ApplicationLayerCode    types.Prompt
	AdaptersLayerCode       types.Prompt
	ErrorHandlingAndLogging types.Prompt
	AdditionalTechnicalFile types.Prompt
	DateTimeHandlingFile    types.Prompt

	// Helper Prompts
	ExtractTests types.Prompt
}

func NewPrompts() *Prompts {
	return &Prompts{
		// Master Prompts
		ProjectMasterPrompt: project_prompts.ProjectMasterPrompt,
		FileMasterPrompt:    file_prompts.FileMasterPrompt,

		// Project-Level Prompts
		ProjectStructure:        project_prompts.ProjectStructurePrompt,
		KeyFiles:                project_prompts.KeyFilesPrompt,
		ApplicationArchitecture: project_prompts.ApplicationArchitecturePrompt,
		DependencyManagement:    project_prompts.DependencyManagementPrompt,
		ProjectSettings:         project_prompts.ProjectSettingsPrompt,
		TestingStrategy:         project_prompts.TestingStrategyPrompt,
		AdditionalTechnical:     project_prompts.AdditionalTechnicalPrompt,
		DateTimeHandling:        project_prompts.DateTimeHandlingPrompt,

		// File-Level Prompts
		CodingStandards:         file_prompts.CodingStandardsPrompt,
		ApplicationLayerCode:    file_prompts.ApplicationLayerCodePrompt,
		AdaptersLayerCode:       file_prompts.AdaptersLayerCodePrompt,
		ErrorHandlingAndLogging: file_prompts.ErrorHandlingAndLoggingPrompt,
		AdditionalTechnicalFile: file_prompts.AdditionalTechnicalFilePrompt,
		DateTimeHandlingFile:    file_prompts.DateTimeHandlingFilePrompt,

		// Helper Prompts
		ExtractTests: helper_prompts.ExtractTestsPrompt,
	}
}
