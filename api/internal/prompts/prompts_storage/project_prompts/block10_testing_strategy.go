// internal/prompts/prompts_storage/project_prompts/block10_testing_strategy.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type TestingStrategyData struct {
	ProjectTree       string
	TestsFilesContent string
}

func (t TestingStrategyData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Project Tree", Description: "Structure of the project", Content: t.ProjectTree},
		{Name: "Tests Files Content", Description: "Files related to testing", Content: t.TestsFilesContent},
	}
}

var TestingStrategyPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software testing, please review the project's testing strategy and structure.",
	BaseTaskDesc: "Evaluate the project's testing strategy and structure.\n\nGuidelines:\n\nUnit Tests: Prioritized, with adapters mocked.\nIntegration Tests: Use SQLite in-memory databases.\nTest Structure: Mirrors project structure; test files correspond to modules/classes.\nTesting Practices: Tests cover various scenarios, including edge cases.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the testing strategy meets the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
