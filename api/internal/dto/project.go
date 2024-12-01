// internal/dto/project.go

package dto

import "mime/multipart"

type UploadProjectRequest struct {
	UserID string                `json:"user_id"`
	Name   string                `json:"name"`
	File   *multipart.FileHeader `json:"file"`
}

type UploadProjectResponse struct {
	Message string     `json:"message"`
	Project ProjectDTO `json:"project"`
}

type ProjectDTO struct {
	ID                    uint   `json:"id"`
	ProgrammingLanguageID uint   `json:"programming_language_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Path                  string `json:"path"`
	Tree                  string `json:"tree"`
}

type AnalyzeProjectRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
}

type AnalyzeProjectResponse struct {
	Message string `json:"message"`
}

type AnalyzeFileRequest struct {
	FileID uint `json:"file_id" binding:"required"`
}

type AnalyzeFileResponse struct {
	Message string `json:"message"`
}

// New DTOs for the first endpoint
type GetAllProjectsResponse struct {
	Projects []ProjectDTO `json:"projects"`
}

// DTO for ProjectFile
type ProjectFileDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	WasAnalyzed bool   `json:"was_analyzed"`
}

// DTO for ProjectAnalysisResult
type ProjectAnalysisResultDTO struct {
	ID              uint   `json:"id"`
	PromptName      string `json:"prompt_name"`
	Compliance      string `json:"compliance"`
	Issues          string `json:"issues"`
	Recommendations string `json:"recommendations"`
}

// DTO for the second endpoint
type GetProjectOverviewResponse struct {
	Project                ProjectDTO                 `json:"project"`
	Files                  []ProjectFileDTO           `json:"files"`
	ProjectAnalysisResults []ProjectAnalysisResultDTO `json:"analysis_results"`
}
