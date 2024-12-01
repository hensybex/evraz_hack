// internal/dto/llm_responses/coding_standards_response.go

package llm_responses

type CodingStandardsResponse struct {
	Compliance      string   `json:"compliance"`
	Issues          []string `json:"issues"`
	Recommendations []string `json:"recommendations"`
}
