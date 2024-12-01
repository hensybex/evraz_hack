// internal/di/di.go

package di

import (
	"evraz_api/internal/config"
	"evraz_api/internal/handler"
	"evraz_api/internal/repository"
	"evraz_api/internal/service"
	"evraz_api/internal/usecase"
	"gorm.io/gorm"
)

type DIContainer struct {
	Config             *config.Config
	DB                 *gorm.DB
	ProjectFileRepo    repository.ProjectFileRepository
	ProjectRepo        repository.ProjectRepository
	FileManager        service.FileManager
	ProjectFileUsecase usecase.ProjectFileUsecase
	ProjectUsecase     usecase.ProjectUsecase
	ProjectHandlers    *handler.ProjectHandlers
}

func NewDIContainer(cfg *config.Config, db *gorm.DB) *DIContainer {
	// Initialize repositories
	projectFileRepo := repository.NewGormProjectFileRepository(db)
	projectRepo := repository.NewGormProjectRepository(db)
	fileManager := service.NewOSFileManager()
	projectAnalysisRepo := repository.NewGormProjectAnalysisRepository(db)
	fileAnalysisRepo := repository.NewGormFileAnalysisRepository(db)
	gptCallRepo := repository.NewGormGPTCallRepository(db)

	projectFileUsecase := usecase.NewProjectFileUsecase(projectFileRepo)
	projectUsecase := usecase.NewProjectUsecase(projectRepo, projectFileRepo, projectAnalysisRepo, fileManager)

	// Initialize services
	mistralService := service.NewMistralService(gptCallRepo)

	// Initialize use cases
	projectAnalysisUsecase := usecase.NewProjectAnalysisUsecase(
		projectRepo,
		projectFileRepo,
		projectAnalysisRepo,
		fileAnalysisRepo,
		*mistralService,
	)

	// Initialize handlers
	projectHandlers := handler.NewProjectHandlers(
		projectUsecase,
		projectFileUsecase,
		projectAnalysisUsecase,
		fileAnalysisRepo,
	)

	return &DIContainer{
		Config:             cfg,
		DB:                 db,
		ProjectFileRepo:    projectFileRepo,
		ProjectRepo:        projectRepo,
		FileManager:        fileManager,
		ProjectFileUsecase: *projectFileUsecase,
		ProjectUsecase:     *projectUsecase,
		ProjectHandlers:    projectHandlers,
	}
}
