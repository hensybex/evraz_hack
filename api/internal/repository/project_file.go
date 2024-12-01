// internal/repository/project_file.go

package repository

import (
	"evraz_api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type ProjectFileRepository interface {
	CreateOne(file *model.ProjectFile) error
	GetOneByID(id uint) (*model.ProjectFile, error)
	GetManyByProjectID(projectID uint) ([]model.ProjectFile, error)
	UpdateOneByID(file *model.ProjectFile) error
	GetFileContentByPath(projectID uint, filePath string) (string, error)
	GetFileContentByName(projectID uint, fileName string) (string, error)
	GetRootFileContentByName(projectID uint, fileName string) (string, error)
	GetFilesByProjectID(projectID uint) ([]model.ProjectFile, error)
}

type GormProjectFileRepository struct {
	db *gorm.DB
}

func NewGormProjectFileRepository(db *gorm.DB) *GormProjectFileRepository {
	return &GormProjectFileRepository{db: db}
}

func (repo *GormProjectFileRepository) CreateOne(file *model.ProjectFile) error {
	return repo.db.Create(file).Error
}

func (repo *GormProjectFileRepository) GetOneByID(id uint) (*model.ProjectFile, error) {
	var file model.ProjectFile
	if err := repo.db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (repo *GormProjectFileRepository) GetManyByProjectID(projectID uint) ([]model.ProjectFile, error) {
	var files []model.ProjectFile
	if err := repo.db.Where(&model.ProjectFile{ProjectID: projectID}).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (repo *GormProjectFileRepository) UpdateOneByID(file *model.ProjectFile) error {
	return repo.db.Save(file).Error
}

func (repo *GormProjectFileRepository) GetFileContentByPath(projectID uint, filePath string) (string, error) {
	var file model.ProjectFile

	// Query the database for the file with the given ProjectID and Path
	if err := repo.db.Where("project_id = ? AND path = ?", projectID, filePath).First(&file).Error; err != nil {
		// If the file is not found, or another error occurs, return an error
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("file with path '%s' not found for project ID %d", filePath, projectID)
		}
		return "", err
	}

	// Return the content of the file
	return file.Content, nil
}

func (repo *GormProjectFileRepository) GetFileContentByName(projectID uint, fileName string) (string, error) {
	var file model.ProjectFile

	// Query the database for the file with the given ProjectID and Path
	if err := repo.db.Where("project_id = ? AND name = ?", projectID, fileName).First(&file).Error; err != nil {
		// If the file is not found, or another error occurs, return an error
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("file with path '%s' not found for project ID %d", fileName, projectID)
		}
		return "", err
	}

	// Return the content of the file
	return file.Content, nil
}

func (repo *GormProjectFileRepository) GetRootFileContentByName(projectID uint, fileName string) (string, error) {
	var file model.ProjectFile

	// Query the database for the file in the root directory
	if err := repo.db.Where("project_id = ? AND (path = ? OR path = ?)", projectID, fileName, "/"+fileName).First(&file).Error; err != nil {
		// If the file is not found, or another error occurs, return an error
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("file '%s' not found in the root directory for project ID %d", fileName, projectID)
		}
		return "", err
	}

	// Return the content of the file
	return file.Content, nil
}

func (repo *GormProjectFileRepository) GetFilesByProjectID(projectID uint) ([]model.ProjectFile, error) {
	var files []model.ProjectFile
	if err := repo.db.Where("project_id = ?", projectID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
