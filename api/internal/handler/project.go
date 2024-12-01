// internal/handler/project.go

package handler

import (
	"bytes"
	"evraz_api/internal/dto"
	"evraz_api/internal/repository"
	"evraz_api/internal/usecase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
)

type ProjectHandlers struct {
	ProjectUsecase         *usecase.ProjectUsecase
	ProjectFileUsecase     *usecase.ProjectFileUsecase
	ProjectAnalysisUsecase *usecase.ProjectAnalysisUsecase
	FileAnalysisRepo       repository.FileAnalysisRepository
}

func NewProjectHandlers(
	projectUsecase *usecase.ProjectUsecase,
	projectFileUsecase *usecase.ProjectFileUsecase,
	projectAnalysisUsecase *usecase.ProjectAnalysisUsecase,
	fileAnalysisRepo repository.FileAnalysisRepository,
) *ProjectHandlers {
	return &ProjectHandlers{
		ProjectUsecase:         projectUsecase,
		ProjectFileUsecase:     projectFileUsecase,
		ProjectAnalysisUsecase: projectAnalysisUsecase,
		FileAnalysisRepo:       fileAnalysisRepo,
	}
}

func (h *ProjectHandlers) UploadProject(c *gin.Context) {
	var req dto.UploadProjectRequest

	// Bind request JSON to DTO
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Get the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	req.File = file

	// Call the use case with the DTO
	projectDTO, err := h.ProjectUsecase.UploadProject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response using DTO
	resp := dto.UploadProjectResponse{
		Message: "Project uploaded successfully",
		Project: projectDTO,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *ProjectHandlers) AnalyzeProject(c *gin.Context) {
	var req dto.AnalyzeProjectRequest

	// Bind JSON input to DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Call the use case with the project ID
	if err := h.ProjectAnalysisUsecase.AnalyzeProject(req.ProjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response DTO
	resp := dto.AnalyzeProjectResponse{
		Message: "Project analysis completed",
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProjectHandlers) AnalyzeFile(c *gin.Context) {
	var req dto.AnalyzeFileRequest

	// Bind JSON input to DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Call the use case with the file ID
	if err := h.ProjectAnalysisUsecase.AnalyzeFile(req.FileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response DTO
	resp := dto.AnalyzeFileResponse{
		Message: "File analysis completed",
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for the "all projects" endpoint
func (h *ProjectHandlers) GetAllProjects(c *gin.Context) {
	projects, err := h.ProjectUsecase.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projectDTOs := make([]dto.ProjectDTO, len(projects))
	for i, project := range projects {
		projectDTOs[i] = dto.ProjectDTO{
			ID:                    project.ID,
			ProgrammingLanguageID: project.ProgrammingLanguageID,
			Name:                  project.Name,
			Path:                  project.Path,
			Tree:                  project.Tree,
		}
	}

	resp := dto.GetAllProjectsResponse{
		Projects: projectDTOs,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for the "project overview" endpoint
func (h *ProjectHandlers) GetProjectOverview(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}

	project, files, analysisResults, err := h.ProjectUsecase.GetProjectOverview(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projectDTO := dto.ProjectDTO{
		ID:                    project.ID,
		ProgrammingLanguageID: project.ProgrammingLanguageID,
		Name:                  project.Name,
		Path:                  project.Path,
		Tree:                  project.Tree,
	}

	fileDTOs := make([]dto.ProjectFileDTO, len(files))
	for i, file := range files {
		fileDTOs[i] = dto.ProjectFileDTO{
			ID:          file.ID,
			Name:        file.Name,
			WasAnalyzed: file.WasAnalyzed,
		}
	}

	analysisResultDTOs := make([]dto.ProjectAnalysisResultDTO, len(analysisResults))
	for i, result := range analysisResults {
		analysisResultDTOs[i] = dto.ProjectAnalysisResultDTO{
			ID:              result.ID,
			PromptName:      result.PromptName,
			Compliance:      result.Compliance,
			Issues:          result.Issues,
			Recommendations: result.Recommendations,
		}
	}

	resp := dto.GetProjectOverviewResponse{
		Project:                projectDTO,
		Files:                  fileDTOs,
		ProjectAnalysisResults: analysisResultDTOs,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProjectHandlers) GetFileAnalysisResults(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file_id"})
		return
	}

	analysisResults, err := h.FileAnalysisRepo.GetManyByFileID(uint(fileID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	analysisResultDTOs := make([]dto.FileAnalysisResultDTO, len(analysisResults))
	for i, result := range analysisResults {
		analysisResultDTOs[i] = dto.FileAnalysisResultDTO{
			ID:              result.ID,
			PromptName:      result.PromptName,
			Compliance:      result.Compliance,
			Issues:          result.Issues,
			Recommendations: result.Recommendations,
		}
	}

	resp := dto.GetFileAnalysisResultsResponse{
		FileAnalysisResults: analysisResultDTOs,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProjectHandlers) GenerateProjectPDF(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid project_id '%s': %v\n", projectIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}
	fmt.Printf("Generating PDF for project ID: %d\n", projectID)

	project, _, analysisResults, err := h.ProjectUsecase.GetProjectOverview(uint(projectID))
	if err != nil {
		fmt.Printf("Error getting project overview: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Retrieved project: %+v\n", project)

	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(project.Name, false)
	pdf.AddPage()
	fmt.Println("Initialized PDF document.")

	// Register the font that supports Russian characters
	fontDir := "fonts/" // Ensure this path is correct
	fontName := "DejaVuSans"
	fontFile := fontDir + "DejaVuSans.ttf"
	fmt.Printf("Registering font '%s' from file '%s'\n", fontName, fontFile)
	// Register regular font
	pdf.AddUTF8Font(fontName, "", fontDir+"DejaVuSans.ttf")
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error adding regular font: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Regular font registered successfully.")

	// Register bold font
	pdf.AddUTF8Font(fontName, "B", fontDir+"DejaVuSans-Bold.ttf")
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error adding bold font: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Font registered successfully.")

	// Set the font
	pdf.SetFont(fontName, "", 16)
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error setting font: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Font set successfully.")

	// Add Project Title
	title := fmt.Sprintf("Project Report: %s", project.Name)
	pdf.Cell(40, 10, title)
	pdf.Ln(20)
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error after adding project title: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Added project title to PDF.")

	// Set font size for content
	pdf.SetFont(fontName, "", 12)
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error setting font size for content: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Font size set for content.")

	// Add Project Details
	pdf.Cell(40, 10, fmt.Sprintf("Project ID: %d", project.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Programming Language ID: %d", project.ProgrammingLanguageID))
	pdf.Ln(10)
	// Add other project details as needed
	if pdf.Err() {
		errMsg := fmt.Sprintf("Error after adding project details: %v", pdf.Error())
		fmt.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Println("Added project details to PDF.")

	// Add Analysis Results
	for _, result := range analysisResults {
		pdf.SetFont(fontName, "B", 14)
		if pdf.Err() {
			errMsg := fmt.Sprintf("Error setting bold font: %v", pdf.Error())
			fmt.Println(errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}

		pdf.Cell(40, 10, result.PromptName)
		pdf.Ln(10)

		pdf.SetFont(fontName, "", 12)
		if pdf.Err() {
			errMsg := fmt.Sprintf("Error setting normal font: %v", pdf.Error())
			fmt.Println(errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}

		pdf.MultiCell(0, 10, fmt.Sprintf("Compliance: %s", result.Compliance), "", "", false)
		pdf.MultiCell(0, 10, fmt.Sprintf("Issues: %s", result.Issues), "", "", false)
		pdf.MultiCell(0, 10, fmt.Sprintf("Recommendations: %s", result.Recommendations), "", "", false)
		pdf.Ln(10)

		if pdf.Err() {
			errMsg := fmt.Sprintf("Error after adding analysis result for prompt '%s': %v", result.PromptName, pdf.Error())
			fmt.Println(errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}
		fmt.Printf("Added analysis result for prompt '%s' to PDF.\n", result.PromptName)
	}

	// Output PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		fmt.Printf("Error generating PDF: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}
	fmt.Println("PDF generated successfully.")

	// Set headers and send PDF
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", project.Name))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
	fmt.Println("PDF sent to client.")
}
