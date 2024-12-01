// internal/prompts/prompts_storage/project_prompts/block12_date_time_handling.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
)

type DateTimeHandlingData struct {
	DateTimeCodeSamples    string
	ETLProcesses           string
	TimezoneConfigurations string
}

func (d DateTimeHandlingData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Date and Time Code Samples",
			Description: "Code snippets where date and time operations are performed",
			Content:     d.DateTimeCodeSamples,
		},
		{
			Name:        "ETL Processes",
			Description: "Descriptions and code of ETL processes handling dates",
			Content:     d.ETLProcesses,
		},
		{
			Name:        "Timezone Configurations",
			Description: "Settings and configurations related to timezones",
			Content:     d.TimezoneConfigurations,
		},
	}
}

var DateTimeHandlingPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in date and time handling in software applications, please review the project's approach to managing dates and times.",
	BaseTaskDesc: "Assess the project's handling of date and time data.\n\nGuidelines:\n\nStore all times in the database as UTC.\nBackend calculations use UTC.\nManage aware and naive datetime objects appropriately.\nConvert timezones correctly in ETL tasks.\nDefine timezone settings in project configurations.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether date and time handling meets the requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
