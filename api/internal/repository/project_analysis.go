// internal/repository/project_analysis.go

package repository

import (
	"evraz_api/internal/model"
	"log"

	"gorm.io/gorm"
)

type ProjectAnalysisRepository interface {
	CreateOne(analysis *model.ProjectAnalysisResult) error
	//GetFilesByProjectID(projectID uint) ([]model.ProjectFile, error)
	GetResultsByProjectID(projectID uint) ([]model.ProjectAnalysisResult, error)
	// Add more methods as needed
}

type GormProjectAnalysisRepository struct {
	db *gorm.DB
}

func NewGormProjectAnalysisRepository(db *gorm.DB) *GormProjectAnalysisRepository {
	return &GormProjectAnalysisRepository{db: db}
}

func (repo *GormProjectAnalysisRepository) CreateOne(analysis *model.ProjectAnalysisResult) error {
	log.Printf("Attempting to create FileAnalysis: %+v", analysis)
	err := repo.db.Create(analysis).Error
	if err != nil {
		log.Printf("Error during FileAnalysis creation: %v", err)
	} else {
		log.Println("FileAnalysis creation succeeded.")
	}
	return err
}

func (repo *GormProjectAnalysisRepository) GetResultsByProjectID(projectID uint) ([]model.ProjectAnalysisResult, error) {
	var results []model.ProjectAnalysisResult
	if err := repo.db.Where("project_id = ?", projectID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
