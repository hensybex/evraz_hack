// internal/dto/llm_responses/project_analysis_response.go

package llm_responses

type ProjectAnalysisResponse struct {
	Description   string         `json:"description"`
	FilesOverview []FileOverview `json:"files_overview"`
}

type FileOverview struct {
	ID      uint   `json:"id"`
	Purpose string `json:"purpose"`
	Context string `json:"context"`
}
