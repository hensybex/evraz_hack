// internal/usecase/project_file.go

package usecase

import (
	"evraz_api/internal/model"
	"evraz_api/internal/repository"
)

type ProjectFileUsecase struct {
	Repo repository.ProjectFileRepository
}

func NewProjectFileUsecase(repo repository.ProjectFileRepository) *ProjectFileUsecase {
	return &ProjectFileUsecase{
		Repo: repo,
	}
}

func (uc *ProjectFileUsecase) CreateOne(file *model.ProjectFile) error {
	return uc.Repo.CreateOne(file)
}

func (uc *ProjectFileUsecase) GetOneByID(id uint) (*model.ProjectFile, error) {
	return uc.Repo.GetOneByID(id)
}

func (uc *ProjectFileUsecase) GetManyByProjectID(id uint) ([]model.ProjectFile, error) {
	return uc.Repo.GetManyByProjectID(id)
}

func (uc *ProjectFileUsecase) GetManyByProjectFileID(id uint) ([]model.ProjectFile, error) {
	return uc.Repo.GetManyByProjectID(id)
}

func (uc *ProjectFileUsecase) UpdateOneByID(file *model.ProjectFile) error {
	return uc.Repo.UpdateOneByID(file)
}

func (uc *ProjectFileUsecase) GetProjectFilesWithAnalysis(projectID uint) ([]model.ProjectFile, error) {
	return uc.Repo.GetFilesWithAnalysisByProjectID(projectID)
}
