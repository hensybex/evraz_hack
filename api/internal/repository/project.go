// internal/repository/project.go

package repository

import (
	"evraz_api/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateOne(project *model.Project) error
	GetOneByID(id uint) (*model.Project, error)
	UpdateOneByID(project *model.Project) error
	GetAll() ([]model.Project, error)
	GetAllProjects() ([]model.Project, error)
	GetProjectByID(projectID uint) (model.Project, error)
}

type GormProjectRepository struct {
	db *gorm.DB
}

func NewGormProjectRepository(db *gorm.DB) *GormProjectRepository {
	return &GormProjectRepository{db: db}
}

func (repo *GormProjectRepository) CreateOne(project *model.Project) error {
	return repo.db.Create(project).Error
}

func (repo *GormProjectRepository) GetOneByID(id uint) (*model.Project, error) {
	var project model.Project
	if err := repo.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (repo *GormProjectRepository) UpdateOneByID(project *model.Project) error {
	return repo.db.Save(project).Error
}

func (repo *GormProjectRepository) GetAll() ([]model.Project, error) {
	var projects []model.Project
	if err := repo.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (repo *GormProjectRepository) GetAllProjects() ([]model.Project, error) {
	var projects []model.Project
	if err := repo.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (repo *GormProjectRepository) GetProjectByID(projectID uint) (model.Project, error) {
	var project model.Project
	if err := repo.db.First(&project, projectID).Error; err != nil {
		return project, err
	}
	return project, nil
}
