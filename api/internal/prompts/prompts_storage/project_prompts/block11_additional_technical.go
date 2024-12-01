// internal/prompts/prompts_storage/project_prompts/block11_additional_technical.go

package project_prompts

import (
	"evraz_api/internal/prompts/types"
	"fmt"
)

type AdditionalTechnicalData struct {
	TransactionManagementCode string
	AsynchronousCodeUsage     string
	DataScienceDependencies   []string
	MonitoringPlan            string
}

func (d AdditionalTechnicalData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{
			Name:        "Transaction Management Code",
			Description: "Code snippets related to database transaction management",
			Content:     d.TransactionManagementCode,
		},
		{
			Name:        "Asynchronous Code Usage",
			Description: "Description and code examples of asynchronous code used",
			Content:     d.AsynchronousCodeUsage,
		},
		{
			Name:        "Data Science Dependencies",
			Description: "List of data science libraries used",
			Content:     fmt.Sprintf("%v", d.DataScienceDependencies),
		},
		{
			Name:        "Monitoring Plan",
			Description: "Description of the plan for implementing monitoring",
			Content:     d.MonitoringPlan,
		},
	}
}

var AdditionalTechnicalPrompt = types.Prompt{
	BasePrompt:   "As an AI assistant specialized in software engineering best practices, please review the project's technical considerations.",
	BaseTaskDesc: "Ensure adherence to technical requirements at the project level.\n\nGuidelines:\n\nDatabase Transactions: Implement \"Unit of Work\" pattern; avoid nested transactions.\nAsynchronous Code: Justify usage; ensure proper implementation with gevent.\nData Science Dependencies: Use libraries like pandas and numpy only within data science modules.\nMonitoring: Acknowledge current limitations; plan for future implementation.",
	JSONStruct: []types.JSONStruct{
		{Key: "compliance", Description: "(bool) Whether the project meets the additional technical requirements"},
		{Key: "issues", Description: "(list of str) List of any issues found"},
		{Key: "recommendations", Description: "(list of str) Suggestions for improvement"},
	},
}
