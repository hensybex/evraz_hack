// internal/dto/llm_responses/file_analysis_response.go

package llm_responses

type FileAnalysisResponse struct {
	Compliance      bool     `json:"compliance"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
}
