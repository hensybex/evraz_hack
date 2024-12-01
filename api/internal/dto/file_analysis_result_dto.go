package dto

type FileAnalysisResultDTO struct {
	ID              uint   `json:"id"`
	PromptName      string `json:"prompt_name"`
	Compliance      string `json:"compliance"`
	Issues          string `json:"issues"`
	Recommendations string `json:"recommendations"`
}

type GetFileAnalysisResultsResponse struct {
	FileAnalysisResults []FileAnalysisResultDTO `json:"analysis_results"`
}
