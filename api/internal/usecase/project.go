// internal/usecase/project.go

package usecase

import (
	"errors"
	"evraz_api/internal/dto"
	"evraz_api/internal/model"
	"evraz_api/internal/repository"
	"evraz_api/internal/service"
	"fmt"
	"os"
	"path/filepath"
)

type ProjectUsecase struct {
	ProjectRepo               repository.ProjectRepository
	ProjectFileRepo           repository.ProjectFileRepository
	ProjectAnalysisResultRepo repository.ProjectAnalysisRepository
	FileManager               service.FileManager
}

func NewProjectUsecase(projectRepo repository.ProjectRepository, projectFileRepo repository.ProjectFileRepository, projectAnalysisResultRepo repository.ProjectAnalysisRepository, fileManager service.FileManager) *ProjectUsecase {
	return &ProjectUsecase{
		ProjectRepo:               projectRepo,
		ProjectFileRepo:           projectFileRepo,
		ProjectAnalysisResultRepo: projectAnalysisResultRepo,
		FileManager:               fileManager,
	}
}

func (uc *ProjectUsecase) UploadProject(req dto.UploadProjectRequest) (dto.ProjectDTO, error) {
	req.UserID = "1"
	if req.UserID == "" {
		return dto.ProjectDTO{}, errors.New("User ID is required")
	}

	// Create the project directory
	if err := uc.FileManager.CreateProject(req.UserID, req.Name); err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to create project")
	}

	// Define the path where the file will be saved
	tempFilePath := filepath.Join(uc.FileManager.FormulatePath(req.UserID, req.Name, false), req.File.Filename)

	// Save the uploaded file
	if err := uc.FileManager.SaveUploadedFile(req.File, tempFilePath); err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to save file")
	}

	// Determine file type and extract accordingly
	extension := filepath.Ext(req.File.Filename)
	extractPath := uc.FileManager.FormulatePath(req.UserID, req.Name, false)
	var err error
	switch extension {
	case ".zip":
		err = uc.FileManager.ExtractZip(tempFilePath, extractPath)
	case ".tar.gz", ".tgz":
		err = uc.FileManager.ExtractTarGz(tempFilePath, extractPath)
	case ".7z":
		err = uc.FileManager.Extract7z(tempFilePath, extractPath)
	default:
		return dto.ProjectDTO{}, errors.New("unsupported file type")
	}
	if err != nil {
		return dto.ProjectDTO{}, fmt.Errorf("failed to extract archive: %w", err)
	}

	// Delete the temporary archive file
	if err := uc.FileManager.RemovePath(tempFilePath); err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to remove temp file")
	}

	// Adjust the extractedPath if necessary
	extractedPath := extractPath

	entries, err := os.ReadDir(extractedPath)
	if err != nil {
		return dto.ProjectDTO{}, fmt.Errorf("failed to read extracted directory: %w", err)
	}

	// Filter out __MACOSX entries
	filteredEntries := []os.DirEntry{}
	for _, entry := range entries {
		if entry.Name() != "__MACOSX" {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	if len(filteredEntries) == 1 && filteredEntries[0].IsDir() {
		// Only one directory (excluding __MACOSX), adjust the extractedPath
		extractedPath = filepath.Join(extractedPath, filteredEntries[0].Name())
		fmt.Printf("Only one subdirectory detected, adjusting extractedPath to %s\n", extractedPath)
	}

	// Execute tree command to get the directory structure
	treeOutput, err := uc.FileManager.RunTreeAtPath(extractedPath)
	if err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to execute tree command")
	}

	// Create Project object and save to DB
	project := model.Project{
		ProgrammingLanguageID: 1,
		Name:                  req.Name,
		Path:                  extractedPath,
		Tree:                  treeOutput,
		WasAnalyzed:           false,
	}
	if err := uc.ProjectRepo.CreateOne(&project); err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to create project")
	}

	// Process files in directory
	err = uc.FileManager.ProcessFilesInDirectory(extractedPath, func(relPath string, content []byte) error {
		fileName := filepath.Base(relPath)
		projectFile := model.ProjectFile{
			ProjectID:   project.ID,
			Path:        relPath,
			Content:     string(content),
			WasAnalyzed: false,
			Name:        fileName,
		}
		return uc.ProjectFileRepo.CreateOne(&projectFile)
	})
	if err != nil {
		return dto.ProjectDTO{}, errors.New("Failed to process project files")
	}

	return dto.ProjectDTO{
		ID:                    project.ID,
		ProgrammingLanguageID: project.ProgrammingLanguageID,
		Name:                  project.Name,
		Path:                  project.Path,
		Tree:                  project.Tree,
	}, nil
}

func (uc *ProjectUsecase) GetAllProjects() ([]model.Project, error) {
	return uc.ProjectRepo.GetAllProjects()
}

func (uc *ProjectUsecase) GetProjectOverview(projectID uint) (model.Project, []model.ProjectFile, []model.ProjectAnalysisResult, error) {
	project, err := uc.ProjectRepo.GetProjectByID(projectID)
	if err != nil {
		return model.Project{}, nil, nil, err
	}

	files, err := uc.ProjectFileRepo.GetFilesByProjectID(projectID)
	if err != nil {
		return project, nil, nil, err
	}

	analysisResults, err := uc.ProjectAnalysisResultRepo.GetResultsByProjectID(projectID)
	if err != nil {
		return project, files, nil, err
	}

	return project, files, analysisResults, nil
}
